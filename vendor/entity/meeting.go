package entity

import (
	"fmt"
	"time"
	"util"
)

type MeetingTitle = util.Identifier

type MeetingInfo struct {
	// ID            int
	Title MeetingTitle
	// Time          time.Time
	StartTime time.Time
	EndTime   time.Time
}
type Meeting struct {
	MeetingInfo
	Sponsor       *User
	Participators UserList
}

// TODO: abstract Meeting(List) and User(List)

func NewMeeting(info MeetingInfo) *Meeting {
	if info.Title.Empty() {
		// FIXME: more elegant ?
		return nil
	}
	m := new(Meeting)
	m.MeetingInfo = info
	return m
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
	meetings := make(MeetingListRaw, ml.Size()) // CHECK: diff between `len` and `cap` in golang
	for _, m := range ml.Meetings {
		meetings = append(meetings, m) // CHECK: maybe better to use index in golang ?
	}
	return meetings
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
