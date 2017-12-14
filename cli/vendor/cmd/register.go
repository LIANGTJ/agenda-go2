// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	// "log"
	"agendaHttp"
	// "net/http"
	// "io/ioutil"
	// "bytes"
	// "os"
	// "encoding/json"
	// log "util/logger"
	// "entity"
	// "status"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register for further use",
	Long: `register for further use and u need to input username, password.
	it will be better if email and phone is provider`,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Error[registerd]： %v\n", err)
			}
		}()

		fmt.Println("register called")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		email, _ := cmd.Flags().GetString("email")
		phone, _ := cmd.Flags().GetString("phone")
		
		// user := entity.NewUser(username, password, email, phone)
		// if user.Invalid() {
		// 	panic("[error]: user regiestered invalid")
		// }
		user, err := agendaHttp.Register(username, password, email, phone)
		if err != nil {
			panic(err)
		}
		
		
		fmt.Println("register successfully", user.Username)
		


		
		
		

		// // client := &http.Client {}
		// // var u interface{} = user
		// buf := new(bytes.Buffer)
		// json.NewEncoder(buf).Encode(user)
		// registerURL := agendaHttp.RegisterURL()
		// // req, err := http.NewRequest(http.MethodPost, registerURL, buf)
		// resp, err := http.Post(registerURL, "application/json", buf)
		// if err != nil {
		// 	panic(err)
		// }
		// // req.Header.Add("Content-Type","application/json")
		// // resp, err := client.Do(req)
		// // if err != nil {
		// // 	panic(err)
		// // }
		// if resp.Status[0] != '2' {
		// 	agendaHttp.ErrorHandle(resp)
		// 	fmt.Println("register falied", user.Username)
		// 	// fmt.Println(resp.)
		// }
		// defer resp.Body.Close() //一定要关闭resp.Body
		
		// resp.Write(os.Stdout)
	},
}

func init() {
	RootCmd.AddCommand(registerCmd)


	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	registerCmd.Flags().StringP("username", "u", "Anonymous", "register info for username")
	registerCmd.Flags().StringP("password", "p", "", "register info for password")
	registerCmd.Flags().StringP("email", "e", "", "register info for email")
	registerCmd.Flags().StringP("phone", "t", "", "register info for phone")

}
