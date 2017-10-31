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

// import (
// 	"fmt"
// 	"agenda"
// 	"entity"
// 	"github.com/spf13/cobra"
// )

// // addCmd represents the register command
// var addCmd = &cobra.Command{
// 	Use:   "logIn",
// 	Short: "add a detailed meeting to a calendar",
// 	Long: `add a detailed event to a calendar,
// 		   u need to input title, start_time, end_time and participators`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("register called")

// 		username, _ := cmd.Flags().GetString("username")
// 		password, _ := cmd.Flags().GetString("password")

// 		fmt.Println("register called by " + username)
// 		fmt.Println("register with info password: " + password)
// 		fmt.Println("register with info email: " + email)
// 		fmt.Println("register with info phone: " + phone)

// 		userInfo := new UserInfo{Username(username), auth.Auth(password),email, phone}
// 		user := newUser(userInfo)
// 		RegisterUser(user)

// 	},
// }

// func init() {
// 	RootCmd.AddCommand(addCmd)

// 	// Here you will define your flags and configuration settings.

// 	// Cobra supports Persistent Flags which will work for this command
// 	// and all subcommands, e.g.:
// 	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

// 	// Cobra supports local flags which will only run when this command
// 	// is called directly, e.g.:
// 	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

// 	addCmd.Flags().StringP("title", "t", "Anonymous", "meeting info for username")
// 	addCmd.Flags().StringP("start", "s", "", "meeting info for start_time")
// 	addCmd.Flags().StringP("end", "e", "", "meeting info for end_time")
// 	addCmd.Flags().StringP("participator", "p", "", "meeting info for participator")

// }
