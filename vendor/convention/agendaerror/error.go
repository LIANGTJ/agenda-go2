package agendaerror

import "errors"

var (
	ErrNeedImplement = errors.New("this function need to be implemented")

	// User
	NilUser           = errors.New("a nil user/*user is to be used")
	ExistedUser       = errors.New("the user has been existed")
	UserNotFound      = errors.New("cannot find the user")
	UserNotRegistered = errors.New("cannot find the user for registered one")
	UserNotLogined    = errors.New("no user logined")

	EmptyUsername = errors.New("given username cannot be empty")

	UserAuthority = errors.New("only the User self can modify his/her account")

	SponsorAuthority      = errors.New("only the sponsor can modify the meeting")
	SponsorResponsibility = errors.New("the sponsor can only cancel but not quit the meeting")

	// Meeting
	NilMeeting          = errors.New("a nil meeting/*meeting is to be used")
	ExistedMeeting      = errors.New("the meeting has been existed")
	ExistedMeetingTitle = errors.New("the meeting title has been existed")
	MeetingNotFound     = errors.New("cannot find the meeting")

	EmptyMeetingTitle = errors.New("given meeting title cannot be empty")

	// Time
	// InvalidTime         = errors.New("startTime/EndTime is not valid")
	InvalidTimeInterval    = errors.New("endTime must be after StartTime")
	ConflictedTimeInterval = errors.New("given time interval conflicts with existed interval")

	// Information
	GivenConflictedInfo = errors.New("given a not reasonable information")
)

type AgendaError struct {
	msg string
}

func NewAgendaError(msg string) *AgendaError {
	return &AgendaError{
		msg: msg,
	}
}

func (e *AgendaError) Error() string {
	return e.msg
}
