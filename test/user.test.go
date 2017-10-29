package main

// import "errors"
import (
	"encoding/json"
	"entity"
	"fmt"
	"log"
	"model"
	"os"
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
	fin, err := os.Open(model.UserDataPath())
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(fin)

	// u0, err := entity.LoadedUser(decoder)
	// fmt.Printf("+v: %+v\n", u0)

	ul := entity.NewUserList()
	// ul, err := entity.LoadedUserList(decoder)
	ul.LoadFrom(decoder)
	fmt.Printf("+v: %+v\n", ul)

	// ul.Add(entity.NewUser(entity.UserInfo{Name: "a", Auth: "a", Mail: "a"}))
	ul.Users["a"] = entity.NewUser(entity.UserInfo{Name: "a", Auth: "a", Mail: "a"})
	ul.Add(entity.NewUser(entity.UserInfo{}))
	ul.Add(entity.NewUser(entity.UserInfo{"aa", "aa", "aa", "123"}))
	ul.Add(entity.NewUser(entity.UserInfo{"b", "b", "b", "123"}))
	ul.Add(entity.NewUser(entity.UserInfo{"bb", "bb", "bb", "123"}))

	if err := ul.ForEach(func(u *entity.User) error {
		println(counter, u.Name, ul.Users[u.Name])
		defer count()
		return nil
	}); err != nil {
		panic(err)
	}

	{
		oldSize := ul.Size()
		u := ul.Ref("bb")
		logf("ul.Size(): %v ---> %v, u: %+v", oldSize, ul.Size(), u)
	}
	{
		oldSize := ul.Size()
		u, _ := ul.PickOut("bb")
		logf("ul.Size(): %v ---> %v, u: %+v", oldSize, ul.Size(), u)
	}

	// fout, err := os.Create(model.UserTestPath())
	// if err != nil {
	// 	panic(err)
	// }
	// encoder := json.NewEncoder(fout)
	// u.Save(encoder)

	// os.MkdirAll(util.WorkingDir(), 0777)
	fout, err := os.Create(model.UserDataPath())
	if err != nil {
		log.Println(err)
	}
	encoder := json.NewEncoder(fout)
	if err := ul.Save(encoder); err != nil {
		logln(err)
	}

}
