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
	"net/http"
	"io/ioutil"
	"bytes"
	// "agenda"
	"encoding/json"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register for further use",
	Long: `register for further use and u need to input username, password.
	it will be better if email and phone is provider`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("registere called")
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Error[registerd4]： %v\n", err)
			}
		}()

		fmt.Println("registere called")
		fmt.Println("register called--2")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		email, _ := cmd.Flags().GetString("email")
		phone, _ := cmd.Flags().GetString("phone")
		fmt.Println("register called2")
		type User struct {
			Username string
			Password string
			Email string
			Phone string
		}
		fmt.Println("register called2")
		client := &http.Client {}
		var registerURL = "https://private-12576-agenda32.apiary-mock.com/v1/users"
		u := User {
			username,
			password,
			email,
			phone,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(u)
		req, err := http.NewRequest(http.MethodPost, registerURL,buf)
		if err != nil {
			fmt.Print("[error1]")
			panic(err)
		}
		req.Header.Add("Content-Type","application/json")
		resp, err := client.Do(req)//发送请求
		if err != nil {
			fmt.Println("Fatal error ", err.Error())
		}
		defer resp.Body.Close()//一定要关闭resp.Body
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Fatal error ", err.Error())
		}
	
		fmt.Println(string(content))
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
