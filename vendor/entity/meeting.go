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
	// Time          time.Time
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

const TimeLayout = time.RFC3339

func (info MeetingInfo) toSerializable() *MeetingInfoSerializable {
	serialInfo := new(MeetingInfoSerializable)

	serialInfo.Title = info.Title
	sponsor, name := info.Sponsor, util.EmptyIdentifier
	if sponsor != nil {
		name = sponsor.Name
	}
	serialInfo.Sponsor = name
	// LOGF("mInfo.toSerializable: %+v\n", serialInfo)
	serialInfo.Participators = info.Participators.toSerializable()
	// us := info.Participators
	// LOGF(".. %+v \n", us)
	// sus := us.toSerializable()
	// LOGF(".. %+v \n", sus)
	// serialInfo.Participators = sus
	// LOGF("mInfo.toSerializable: %+v\n", serialInfo)
	serialInfo.StartTime = info.StartTime.Format(TimeLayout)
	serialInfo.EndTime = info.EndTime.Format(TimeLayout)
	return serialInfo
}
func (infoSerial MeetingInfoSerializable) toRegular() *MeetingInfo {
	info := new(MeetingInfo)

	info.Title = infoSerial.Title
	USERS := GetUserTableTotal()
	info.Sponsor = USERS.Get(infoSerial.Sponsor)
	for _, infoSerial := range infoSerial.Participators {
		u := USERS.Get(infoSerial.Name)
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

type Meeting struct {
	MeetingInfo
}

// TODO: abstract Meeting(List) and User(List)

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

func DeserializeMeeting(decoder Decoder) (*Meeting, error) {
	mInfoSerial := new(MeetingInfoSerializable)

	err := decoder.Decode(mInfoSerial)
	if err != nil {
		log.Fatal(err) // FIXME:
		return nil, err
	}
	meeting := NewMeeting(*(mInfoSerial.toRegular()))
	return meeting, nil
}

func (m Meeting) Serialize(encoder Encoder) error {
	return encoder.Encode(*m.MeetingInfo.toSerializable())
}

// ................................................................

type MeetingList struct {
	Meetings map[MeetingTitle](*Meeting)
}

type MeetingListRaw = []*Meeting

const (
	defaultMeetingListLength = 5
)

func NewMeetingList() *MeetingList {
	ml := new(MeetingList)
	ml.Meetings = make(map[MeetingTitle](*Meeting))
	return ml
}

func DeserializeMeetingList(decoder Decoder) (*MeetingList, error) {
	ml := NewMeetingList()
	for decoder.More() {
		// TODO: MeetingInfo <---> MeetingInfoSerializable integration
		mInfoSerial := new(MeetingInfoSerializable)
		err := decoder.Decode(mInfoSerial)
		if err != nil {
			log.Fatal(err) // FIXME:
			return nil, err
		}
		meeting := NewMeeting(*(mInfoSerial.toRegular()))
		if err := ml.Add(meeting); err != nil {
			log.Fatal(err)
			return ml, err // FIXME:
		}
	}
	return ml, nil
}

func (ml *MeetingList) Serialize(encoder Encoder) error {
	sl := ml.toSerializable()
	// LOGF("sl: %+v\n", sl)
	return encoder.Encode(sl)
}

func (ml MeetingList) Size() int {
	return len(ml.Meetings)
}

func (ml *MeetingList) Get(title MeetingTitle) *Meeting {
	return ml.Meetings[title] // NOTE: if directly return accessed result from a map like this, would not get the (automatical) `ok`
}
func (ml MeetingList) Contains(title MeetingTitle) bool {
	m := ml.Get(title)
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
	m := ml.Get(title)
	if m == nil {
		return m, ErrMeetingNotFound
	}
	defer ml.Remove(m)
	return m, nil
}

func (ml MeetingList) Slice() MeetingListRaw {
	meetings := make(MeetingListRaw, 0, ml.Size())
	for _, m := range ml.Meetings {
		meetings = append(meetings, m) // CHECK: maybe better to use index in golang ?
	}
	return meetings
}

func (ml MeetingList) toSerializable() []MeetingInfoSerializable {
	meetings := ml.Slice()
	ret := make([]MeetingInfoSerializable, 0, ml.Size())

	// LOG("ml.Size(): ", ml.Size())
	// LOGF("toSerializable: %+v \n", meetings)
	for _, m := range meetings {

		// FIXME: these are introduced since up to now, it is possible that MeetingList contains nil Meeting
		if m == nil {
			log.Printf("A nil Meeting is to be used. Just SKIP OVER it.")
			continue
		}

		ret = append(ret, *(m.MeetingInfo.toSerializable()))
		// LOGF("%+v\n", m.MeetingInfo)
		// smi := *(m.MeetingInfo.toSerializable())
		// LOGF("%+v\n", smi)
		// ret = append(ret, smi)
	}
	return ret
}

// func (ml *MeetingList) ForEach(f func(*Meeting)) error {
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
