package err

import "errors"

var (
	NeedImplement = errors.New("This function need to be implemented.")

	// User
	NilUser           = errors.New("A nil user/*user is to be used.")
	ExistedUser       = errors.New("The user has been existed.")
	UserNotFound      = errors.New("Cannot find the user.")
	UserNotRegistered = errors.New("Cannot find the user for registered one.")

	EmptyUsername = errors.New("Given username cannot be empty.")

	SponsorAuthority      = errors.New("Only the sponsor can modify the meeting.")
	SponsorResponsibility = errors.New("The sponsor can only cancel but not quit the meeting.")

	// Meeting
	NilMeeting          = errors.New("A nil meeting/*meeting is to be used.")
	ExistedMeeting      = errors.New("The meeting has been existed.")
	ExistedMeetingTitle = errors.New("The meeting title has been existed.")
	MeetingNotFound     = errors.New("Cannot find the meeting.")

	EmptyMeetingTitle = errors.New("Given meeting title cannot be empty.")

	// Time
	// InvalidTime         = errors.New("StartTime/EndTime is not valid.")
	InvalidTimeInterval    = errors.New("EndTime must be after StartTime.")
	ConflictedTimeInterval = errors.New("Given time interval conflicts with existed interval.")

	// Information
	GivenConflictedInfo = errors.New("Given a not reasonable information")
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
