package main

// import "errors"
import (
	"encoding/json"
	"entity"
	"log"
	"model"
	"os"
	"time"
	"util"
)

var logln = util.Log
var logf = util.Logf

var (
	counter = 0
)

func count() {
	counter += 1
}

func main() {
	fin, err := os.Open(model.MeetingDataPath())
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(fin)

	// u0 := entity.LoadedMeeting(decoder)
	// logf("+v: %+v\n", u0)

	ml := entity.NewMeetingList()
	// ml := entity.LoadedMeetingList(decoder)
	ml.LoadFrom(decoder)
	logf("+v: %+v\n", ml)

	t := time.Now()
	u := entity.NewUser(entity.UserInfo{Name: "Sponsor"})
	// ml.Add(entity.NewMeeting(entity.MeetingInfo{Title: "a", StartTime: "20171023_1723", EndTime: "20171023_1724"}))
	ml.Meetings["a"] = entity.NewMeeting(entity.MeetingInfo{Title: "a", Sponsor: u, StartTime: t, EndTime: t})
	ml.Add(entity.NewMeeting(entity.MeetingInfo{}))
	ml.Add(entity.NewMeeting(entity.MeetingInfo{Title: "aa", Sponsor: u, StartTime: t.AddDate(1, 1, 1), EndTime: t}))
	ml.Add(entity.NewMeeting(entity.MeetingInfo{Title: "b", Sponsor: u, StartTime: t.AddDate(2, 2, 2), EndTime: t}))
	ml.Add(entity.NewMeeting(entity.MeetingInfo{Title: "bb", Sponsor: u, StartTime: t.AddDate(2, 2, 2), EndTime: t}))

	if err := ml.ForEach(func(key entity.MeetingTitle) error {
		logln(counter, key, ml.Meetings[key])
		defer count()
		return nil
	}); err != nil {
		panic(err)
	}

	{
		oldSize := ml.Size()
		m := ml.Ref("bb")
		logf("ml.Size(): %v ---> %v, m: %+v", oldSize, ml.Size(), m)
	}
	{
		oldSize := ml.Size()
		m, _ := ml.PickOut("bb")
		logf("ml.Size(): %v ---> %v, m: %+v", oldSize, ml.Size(), m)

		fout, err := os.Create(model.MeetingTestPath())
		if err != nil {
			panic(err)
		}
		encoder := json.NewEncoder(fout)
		m.Save(encoder)
	}

	// os.MkdirAll(util.WorkingDir(), 0777)
	fout, err := os.Create(model.MeetingDataPath())
	if err != nil {
		log.Println(err)
	}
	encoder := json.NewEncoder(fout)
	if encoder != nil {
		err := ml.Save(encoder)
		if err != nil {
			log.Println(err)
		}
	}

	logln("correct.")
}
