package entity

import (
	"fmt"
	"log"
	"time"
	"util"
)

type MeetingTitle = util.Identifier

type MeetingInfo struct {
	// ID            int
	Title         MeetingTitle
	Sponsor       *User
	Participators UserList

	StartTime time.Time
	EndTime   time.Time
}
type MeetingInfoSerializable struct {
	Title         MeetingTitle
	Sponsor       Username
	Participators []UserInfoSerializable
	StartTime     string // TODO:
	EndTime       string
}

type Meeting struct {
	MeetingInfo
}

func (info *MeetingInfo) Serialize() *MeetingInfoSerializable {
	serialInfo := new(MeetingInfoSerializable)

	serialInfo.Title = info.Title

	sponsor, name := info.Sponsor, util.EmptyIdentifier
	if sponsor != nil {
		name = sponsor.Name
	}
	serialInfo.Sponsor = name

	serialInfo.Participators = info.Participators.Serialize()

	serialInfo.StartTime = info.StartTime.Format(TimeLayout)
	serialInfo.EndTime = info.EndTime.Format(TimeLayout)

	return serialInfo
}
func (infoSerial *MeetingInfoSerializable) Deserialize() *MeetingInfo {
	info := new(MeetingInfo)

	info.Title = infoSerial.Title

	USERS := GetAllUsersRegistered()

	// CHECK: Need ensure Sponsor not nil ?
	info.Sponsor = USERS.Ref(infoSerial.Sponsor)

	for _, infoSerial := range infoSerial.Participators {
		u := USERS.Ref(infoSerial.Name)
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

type MeetingInfoSerializableList []MeetingInfoSerializable

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

func LoadMeeting(decoder Decoder, m *Meeting) {
	mInfoSerial := new(MeetingInfoSerializable)

	err := decoder.Decode(mInfoSerial)
	if err != nil {
		log.Fatal(err) // FIXME:
	}
	m.MeetingInfo = *(mInfoSerial.Deserialize())
}
func LoadedMeeting(decoder Decoder) *Meeting {
	m := new(Meeting)
	LoadMeeting(decoder, m)
	return m
}

func (m *Meeting) Save(encoder Encoder) error {
	return encoder.Encode(*m.MeetingInfo.Serialize())
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

func LoadMeetingList(decoder Decoder, ml *MeetingList) {
	// for decoder.More() {
	// 	meeting := LoadedMeeting(decoder)
	// 	if err := ml.Add(meeting); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	mlSerial := new(MeetingInfoSerializableList)
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
func (ml *MeetingList) LoadFrom(decoder Decoder) {
	LoadMeetingList(decoder, ml)
}

func LoadedMeetingList(decoder Decoder) *MeetingList {
	ml := NewMeetingList()
	LoadMeetingList(decoder, ml)
	return ml
}

func (ml *MeetingList) Serialize() MeetingInfoSerializableList {
	meetings := ml.Slice()
	ret := make(MeetingInfoSerializableList, 0, ml.Size())

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

func (ml MeetingInfoSerializableList) Size() int {
	return len(ml)
}

func (mlSerial MeetingInfoSerializableList) Deserialize() *MeetingList {
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

func (ml *MeetingList) Save(encoder Encoder) error {
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
		return ErrNilMeeting
	}
	title := meeting.Title
	if ml.Contains(title) {
		return ErrExistedMeeting
	}
	ml.Meetings[title] = meeting
	return nil
}
func (ml *MeetingList) Remove(meeting *Meeting) error {
	if meeting == nil {
		return ErrNilMeeting
	}
	title := meeting.Title
	if ml.Contains(title) {
		delete(ml.Meetings, title) // NOTE: never error, according to 'go-maps-in-action'
		return nil
	}
	return ErrMeetingNotFound
}
func (ml *MeetingList) PickOut(title MeetingTitle) (*Meeting, error) {
	if title.Empty() {
		return nil, ErrEmptyMeetingTitle
	}
	m := ml.Ref(title)
	if m == nil {
		return nil, ErrMeetingNotFound
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

// CHECK: should limit func type ?
func (ml *MeetingList) ForEach(fn interface{}) error {
	switch f := fn.(type) {
	case func(MeetingTitle) error:
		for k := range ml.Meetings {
			if err := f(k); err != nil {
				return err
			}
		}
	case func(*Meeting) error:
		for _, v := range ml.Meetings {
			if err := f(v); err != nil {
				return err
			}
		}
	case func(MeetingTitle, *Meeting) error:
		for k, v := range ml.Meetings {
			if err := f(k, v); err != nil {
				return err
			}
		}
	default:
		return NewAgendaError(fmt.Sprintf("Given function has unmatched signature: %T", f))
	}
	return nil
}
