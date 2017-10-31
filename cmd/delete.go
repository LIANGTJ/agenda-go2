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
	"agenda"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete for further use",
	Long: `delete for further use and u need to input username, password.
	it will be better if email and phone is provider`,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Error[delete]： %v\n", err)
			}
		}()

		fmt.Println("delete called")

		meetingBoolFlag, _ := cmd.Flags().GetBool("meeting")
		userBoolFlag, _ := cmd.Flags().GetBool("user")
		titleFlag, _ := cmd.Flags().GetString("title")
		if meetingBoolFlag && userBoolFlag {
			panic(errors.New("The flag meeting(m) and user(u) can't appear at the same time"))

		} else if userBoolFlag {

			username := agenda.LoginedUser().Name
			agenda.ClearAllMeeting()
			agenda.LogOut(username)
			agenda.CancelAccount(username)

		} else {

			title := agenda.MeetingTitle(titleFlag)
			if err := agenda.QuitMeeting(title); err != nil {
				panic(err)
			}
			if err := agenda.CancelMeeting(title); err != nil {
				panic(err)
			}

		}
		agenda.SaveAll()

	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	deleteCmd.Flags().BoolP("meeting", "m", false, "delete info for meeting")
	deleteCmd.Flags().BoolP("user", "u", false, "delete info for user")
	deleteCmd.Flags().StringP("title", "t", "", "delete info for email")

}
