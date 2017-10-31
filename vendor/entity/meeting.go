package entity

import (
	"convention/agendaerror"
	"convention/codec"
	"log"
	"time"
)

// Identifier
type MeetingTitle string

func (t *MeetingTitle) Empty() bool {
	return *t == ""
}
func (t *MeetingTitle) Valid() bool {
	// FIXME: not only !empty
	return !t.Empty()
}

// TODO: Not sure where to place ...
var allMeetings = *(NewMeetingList())

func (title MeetingTitle) RefInAllMeetings() *Meeting {
	return allMeetings.Ref(title)
}
func GetAllMeetings() *MeetingList {
	return &allMeetings
}

var dissolvedMeetings = *(NewMeetingList())

func (title MeetingTitle) RefInDissolvedMeetings() *Meeting {
	return dissolvedMeetings.Ref(title)
}
func GetDissolvedMeetings() *MeetingList {
	return &dissolvedMeetings
}

type MeetingInfo struct {
	// ID            int
	Title         MeetingTitle
	Sponsor       *User
	Participators UserList

	StartTime time.Time
	EndTime   time.Time
}
type MeetingInfoSerializable struct {
	Title   MeetingTitle
	Sponsor Username
	// Participators []UserInfoSerializable
	Participators []Username
	StartTime     string // TODO:
	EndTime       string
}

type MeetingInfoListPrintable = MeetingInfoListSerializable

type Meeting struct {
	MeetingInfo
}

func (info *MeetingInfo) Serialize() *MeetingInfoSerializable {
	serialInfo := new(MeetingInfoSerializable)

	serialInfo.Title = info.Title

	if sponsor := info.Sponsor; sponsor != nil {
		serialInfo.Sponsor = sponsor.Name
	}

	serialInfo.Participators = info.Participators.identifiers()

	serialInfo.StartTime = info.StartTime.Format(TimeLayout)
	serialInfo.EndTime = info.EndTime.Format(TimeLayout)

	return serialInfo
}
func (infoSerial *MeetingInfoSerializable) Deserialize() *MeetingInfo {
	info := new(MeetingInfo)

	info.Title = infoSerial.Title

	// CHECK: Need ensure Sponsor not nil ?
	info.Sponsor = infoSerial.Sponsor.RefInAllUsers()

	// TODO: TODEL:
	for _, name := range infoSerial.Participators {
		u := name.RefInAllUsers() // CHECK: ditto
		info.Participators.Add(u)
	}

	// FIXME: for no error returned
	var err1, err2 error
	info.StartTime, err1 = time.Parse(TimeLayout, infoSerial.StartTime)
	info.EndTime, err2 = time.Parse(TimeLayout, infoSerial.EndTime)
	if err1 != nil || err2 != nil {
		log.Fatalf("time.Parse fail when parsing %v / %v", infoSerial.StartTime, infoSerial.EndTime)
	}

	return info
}

type MeetingInfoListSerializable []MeetingInfoSerializable

func (ml MeetingInfoListSerializable) Size() int {
	return len(ml)
}

func (mlSerial MeetingInfoListSerializable) Deserialize() *MeetingList {
	ret := NewMeetingList()

	for _, mInfoSerial := range mlSerial {

		// FIXME: these are introduced since up to now, it is possible that UserList contains nil User
		// FIXME: Not use `== nil` because `mInfoSerial` is a struct
		if mInfoSerial.Title.Empty() {
			log.Printf("A No-Title MeetingInfo is to be used. Just SKIP OVER it.")
			continue
		}

		m := NewMeeting(*(mInfoSerial.Deserialize()))
		if err := ret.Add(m); err != nil {
			log.Printf(err.Error()) // CHECK:
		}
	}
	return ret
}

const TimeLayout = time.RFC3339

// TODO: abstract Meeting(List) and User(List)

// CHECK: if need, use pointer instead of value
func NewMeeting(info MeetingInfo) *Meeting {
	if info.Title.Empty() {
		// FIXME: more elegant ?
		log.Printf("An empty MeetingInfo is passed to new a Meeting. Just return `nil`.")
		return nil
	}
	m := new(Meeting)
	m.MeetingInfo = info
	return m
}

func LoadMeeting(decoder codec.Decoder, m *Meeting) {
	mInfoSerial := new(MeetingInfoSerializable)

	err := decoder.Decode(mInfoSerial)
	if err != nil {
		log.Fatal(err) // FIXME:
	}
	m.MeetingInfo = *(mInfoSerial.Deserialize())
}
func LoadedMeeting(decoder codec.Decoder) *Meeting {
	m := new(Meeting)
	LoadMeeting(decoder, m)
	return m
}

func (m *Meeting) Save(encoder codec.Encoder) error {
	return encoder.Encode(*m.MeetingInfo.Serialize())
}

// SponsoredBy checks if Meeting sponsored by User
func (m *Meeting) SponsoredBy(name Username) bool {
	return m.Sponsor.Name == name
}

// ContainsParticipator checks if Meeting's participators contains the User
func (m *Meeting) ContainsParticipator(name Username) bool {
	return m.Participators.Contains(name)
}

// Dissolve deletes the Meeting (, not by a User)
func (m *Meeting) Dissolve() error {
	if err := GetAllMeetings().Remove(m); err != nil {
		log.Printf(err.Error())
		return err
	}
	if err := GetDissolvedMeetings().Add(m); err != nil {
		log.Printf(err.Error())
		return err
	}
	return nil
}

// Exclude removes User from Meeting's participators list
func (m *Meeting) Exclude(u *User) error {
	if err := m.Participators.Remove(u); err != nil {
		log.Printf(err.Error())
		return err
	}
	if m.Participators.Size() <= 0 {
		return m.Dissolve()
	}
	return nil
}

// Involve adds User to Meeting's participators list
func (m *Meeting) Involve(u *User) error {
	if err := m.Participators.Add(u); err != nil {
		log.Printf(err.Error())
		return err
	}
	return nil
}

// ................................................................

type MeetingList struct {
	Meetings map[MeetingTitle](*Meeting)
}

type MeetingListRaw = []*Meeting

func NewMeetingList() *MeetingList {
	ml := new(MeetingList)
	ml.Meetings = make(map[MeetingTitle](*Meeting))
	return ml
}

func LoadMeetingList(decoder codec.Decoder, ml *MeetingList) {
	// for decoder.More() {
	// 	meeting := LoadedMeeting(decoder)
	// 	if err := ml.Add(meeting); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	mlSerial := new(MeetingInfoListSerializable)
	if err := decoder.Decode(mlSerial); err != nil {
		log.Fatal(err)
	}
	for _, mInfoSerial := range *mlSerial {
		m := NewMeeting(*(mInfoSerial.Deserialize()))
		if err := ml.Add(m); err != nil {
			log.Printf(err.Error())
		}
	}
}
func (ml *MeetingList) LoadFrom(decoder codec.Decoder) {
	LoadMeetingList(decoder, ml)
}

func LoadedMeetingList(decoder codec.Decoder) *MeetingList {
	ml := NewMeetingList()
	LoadMeetingList(decoder, ml)
	return ml
}

func (ml *MeetingList) Serialize() MeetingInfoListSerializable {
	meetings := ml.Slice()
	ret := make(MeetingInfoListSerializable, 0, ml.Size())

	// logln("ml.Size(): ", ml.Size())
	// logf("Serialize: %+v \n", meetings)
	for _, m := range meetings {

		// FIXME: these are introduced since up to now, it is possible that MeetingList contains nil Meeting
		if m == nil {
			log.Printf("A nil Meeting is to be used. Just SKIP OVER it.")
			continue
		}

		ret = append(ret, *(m.MeetingInfo.Serialize()))
		// logf("%+v\n", m.MeetingInfo)
		// smi := *(m.MeetingInfo.Serialize())
		// logf("%+v\n", smi)
		// ret = append(ret, smi)
	}
	return ret
}

func (ml *MeetingList) Textualize() MeetingInfoListPrintable {
	return ml.Serialize()
}

func (ml *MeetingList) Save(encoder codec.Encoder) error {
	sl := ml.Serialize()
	// logf("sl: %+v\n", sl)
	return encoder.Encode(sl)
}

func (ml *MeetingList) Size() int {
	return len(ml.Meetings)
}

func (ml *MeetingList) Ref(title MeetingTitle) *Meeting {
	return ml.Meetings[title] // NOTE: if directly return accessed result from a map like this, would not get the (automatical) `ok`
}
func (ml *MeetingList) Contains(title MeetingTitle) bool {
	m := ml.Ref(title)
	return m != nil
}

func (ml *MeetingList) Add(meeting *Meeting) error {
	if meeting == nil {
		return agendaerror.ErrNilMeeting
	}
	title := meeting.Title
	if ml.Contains(title) {
		return agendaerror.ErrExistedMeeting
	}
	ml.Meetings[title] = meeting
	return nil
}
func (ml *MeetingList) Remove(meeting *Meeting) error {
	if meeting == nil {
		return agendaerror.ErrNilMeeting
	}
	title := meeting.Title
	if ml.Contains(title) {
		delete(ml.Meetings, title) // NOTE: never error, according to 'go-maps-in-action'
		return nil
	}
	return agendaerror.ErrMeetingNotFound
}
func (ml *MeetingList) PickOut(title MeetingTitle) (*Meeting, error) {
	if title.Empty() {
		return nil, agendaerror.ErrEmptyMeetingTitle
	}
	m := ml.Ref(title)
	if m == nil {
		return nil, agendaerror.ErrMeetingNotFound
	}
	defer ml.Remove(m)
	return m, nil
}

func (ml *MeetingList) Slice() MeetingListRaw {
	meetings := make(MeetingListRaw, 0, ml.Size())
	for _, m := range ml.Meetings {
		meetings = append(meetings, m) // CHECK: maybe better to use index in golang ?
	}
	return meetings
}

// ForEach used to extension/concrete logic for whole MeetingList
func (ml *MeetingList) ForEach(fn func(*Meeting) error) error {
	for _, v := range ml.Meetings {
		if err := fn(v); err != nil {
			// CHECK: Or, lazy error ?
			return err
		}
	}
	return nil
}

// Filter used for all extension/concrete select for whole MeetingList
func (ml *MeetingList) Filter(pred func(Meeting) bool) *MeetingList {
	ret := NewMeetingList()
	for _, m := range ml.Meetings {
		if pred(*m) {
			ret.Add(m)
		}
	}
	return ret
}
