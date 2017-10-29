package model

import "util"

func UserDataRegisteredPath() string { return util.WorkingDir() + "user-registered.json" }
func UserDataPath() string           { return util.WorkingDir() + "user-data.json" }
func UserTestPath() string           { return util.WorkingDir() + "user-test.json" }

func MeetingDataPath() string { return util.WorkingDir() + "meeting-data.json" }
func MeetingTestPath() string { return util.WorkingDir() + "meeting-test.json" }

func AgendaConfigPath() string { return util.WorkingDir() + "config.json" }
