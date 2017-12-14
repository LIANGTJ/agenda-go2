package agendaHttp


var(
	protocol                      = "https://"
	host                          = "private-12576-agenda32.apiary-mock.com"
	registerPath                  = "/v1/users"
	loginPath                     = "/v1/session"
	logoutPath                    = "/v1/session"
	getUserByIDPath               = "/v1/user/"
	deleteUserByIDPath            = "/v1/user/"
	listFilteredMeetingPath       = "/v1/user/"
	deleteFilteredMeetingPath     = "/v1/user/"
	QueryAccountAllPath              = "/v1/users"
	getMeetingByIDPath            = "/v1/meeting/"
	deleteMeetingByIDPath         = "/v1/meeting/"
	modifyMeetingByIDPath         = "/v1/meeting/"
	listMeetingInTimeIntervalPath = "/v1/meetings?"
	createMeetingPath             = "/v1/meetings"
)

func RegisterURL() string {
	return protocol + host + registerPath
}

func LoginURL() string {
	return protocol + host + loginPath
}

func LogoutURL() string{
	return protocol + host + logoutPath
}
func GetUserByIDURL(ID string) string {
	return protocol + host + getUserByIDPath + ID
}

func DeleteUserByIDURL(ID string) string {
	return protocol + host + deleteUserByIDPath + ID
}

func ListFilteredMeetingURL(ID, filter string) string {
	return protocol + host + listFilteredMeetingPath + ID + "/meetings?Filter=" + filter 
}

func DeleteFilteredMeetingURL(ID, filter string) string {
	return protocol + host + deleteFilteredMeetingPath + ID + "/meetings?Filter=" + filter
}

func QueryAccountAllURL() string {
	return protocol + host + QueryAccountAllPath
}

func GetMeetingByIDURL(ID string) string{
	return protocol + host + getMeetingByIDPath + ID
}

func DeleteMeetingByIDURL(ID string) string {
	return protocol + host + deleteMeetingByIDPath + ID
}

func ModifyMeetingByIDURL(ID string) string {
	return protocol + host + modifyMeetingByIDPath + ID
}

func ListMeetingInTimeIntervalURL(startTime, endTime string) string{
	return protocol + host + listMeetingInTimeIntervalPath + "starttime=" + startTime + "&endtime" + endTime
}
func CreateMeetingURL() string {
	return protocol + host + createMeetingPath
}