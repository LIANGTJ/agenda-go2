package main

// import "errors"
import (
	"encoding/json"
	"entity"
	"fmt"
	"model"
	"os"
)

var (
	counter = 0
)

func count() {
	counter += 1
}

func main() {
	path := model.UserDataPath()
	fin, err := os.Open(path)

	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(fin)

	ul, err := entity.DeserializeUserList(decoder)

	// ul := entity.NewUserList()
	// userLoader := new(model.UserLoader)
	// if err := userLoader.Load(ul); err != nil {
	// 	log.Println(err)
	// }
	fmt.Printf("+v: %+v\n", ul)

	// // ul.Users["a"] = NewUser(entity.UserInfo{"a", "a", "a"})
	// ul.Users["a"] = entity.NewUser(entity.UserInfo{})
	// ul.Users["aa"] = entity.NewUser(entity.UserInfo{"aa", "aa", "aa", "123"})
	// ul.Users["b"] = entity.NewUser(entity.UserInfo{"b", "b", "b", "123"})
	// ul.Users["bb"] = entity.NewUser(entity.UserInfo{"bb", "bb", "bb", "123"})

	if err := ul.ForEach(func(key entity.Username) error {
		println(key, ul.Users[key])
		defer count()
		return nil
	}); err != nil {
		panic(err)
	}

	println(ul.Size())
	{
		u := ul.Get("bb")
		println(u)

		println(ul.Size())
	}
	u, _ := ul.PickOut("bb")

	fmt.Println("ul.Size(): ", ul.Size())

	fmt.Printf("+v: %+v\n", u)

	// os.MkdirAll(util.WorkingDir(), 0777)
	// file, err := os.Create(model.UserDataPath())
	// if err != nil {
	// 	panic(err)
	// }
	// encoder := json.NewEncoder(file)
	// ul.Serialize(encoder)

}
