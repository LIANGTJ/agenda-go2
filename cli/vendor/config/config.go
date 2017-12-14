package config
import (
    // "io"
    "log"
    "os"
    "time"
)

var logToConsoleMode = true

func LogToConsoleMode() bool { return logToConsoleMode }

var neededFilePaths = []string {
    UserLoginedStatusPath(),
    LogPath(),

}

func NeededFilePath() []string {
    return neededFilePaths
}

func WorkingDir() string {
    // location , existed := os.LookupEnv("HOME")

    // if !existed {
    //     location = "."
    // }
    location := "."
    workDir := location + "/agenda.d"
    return workDir
}

func UserLoginedStatusPath() string {
    return WorkingDir() + "/curUser.txt"
}

func LogPath() string {
    return WorkingDir() + "agenda_" + time.Now().Format("20170102_15") + ".log"
}

func ensurePathNeededExist() {
    if err :=  os.MkdirAll(WorkingDir(),0777); err != nil {
        log.Fatal(err)
    }

    for _, path := range NeededFilePath() {
        if _, err := os.Stat(path); os.IsNotExist(err) {
            fd, err := os.Create(path)
            defer fd.Close()
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}

func init() {
    ensurePathNeededExist()
}
