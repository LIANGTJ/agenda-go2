package entity

import "errors"

var (
	ErrNeedImplement = errors.New("This function need to be implemented.")

	// User
	ErrNilUser           = errors.New("A nil user/*user is to be used.")
	ErrExistedUser       = errors.New("The user has been existed.")
	ErrUserNotFound      = errors.New("Cannot find the user.")
	ErrUserNotRegistered = errors.New("Cannot find the user for registered one.")

	ErrEmptyUsername = errors.New("Given username cannot be empty.")

	ErrSponsorAuthority      = errors.New("Only the sponsor can modify the meeting.")
	ErrSponsorResponsibility = errors.New("The sponsor can only cancel but not quit the meeting.")

	// Meeting
	ErrNilMeeting          = errors.New("A nil meeting/*meeting is to be used.")
	ErrExistedMeeting      = errors.New("The meeting has been existed.")
	ErrExistedMeetingTitle = errors.New("The meeting title has been existed.")
	ErrMeetingNotFound     = errors.New("Cannot find the meeting.")

	ErrEmptyMeetingTitle = errors.New("Given meeting title cannot be empty.")

	// Time
	// ErrInvalidTime         = errors.New("StartTime/EndTime is not valid.")
	ErrInvalidTimeInterval    = errors.New("EndTime must be after StartTime.")
	ErrConflictedTimeInterval = errors.New("Given time interval conflicts with existed interval.")

	// Information
	ErrGivenConflictedInfo = errors.New("Given a not reasonable information")
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
