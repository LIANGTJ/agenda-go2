package util

// Identifier as a unique identifier, like ID
type Identifier interface {
	Empty() bool
	Vaild() bool
}

// // CHECK: EmptyIdentifier were unexported before planning to use Username/MeetingTitle as a stand-alone type.
// var EmptyIdentifier = *new(Identifier)

// func (n Identifier) Empty() bool {
// 	return n == EmptyIdentifier
// }
