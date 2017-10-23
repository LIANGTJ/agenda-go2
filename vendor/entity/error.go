package entity

import "errors"

var (
	ErrNeedImplement = errors.New("This function need to be implemented.")

	// User
	ErrNilUser      = errors.New("A nil user/*user is to be used.")
	ErrExistedUser  = errors.New("The user has been existed.")
	ErrUserNotFound = errors.New("Cannot find the user.")

	ErrEmptyUsername = errors.New("Given username cannot be empty.")

	// Meeting
	ErrNilMeeting      = errors.New("A nil meeting/*meeting is to be used.")
	ErrExistedMeeting  = errors.New("The meeting has been existed.")
	ErrMeetingNotFound = errors.New("Cannot find the meeting.")

	ErrEmptyMeetingTitle = errors.New("Given meeting title cannot be empty.")

	// Err
	// Err
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
