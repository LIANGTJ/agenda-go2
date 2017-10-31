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
	"gihub.com/olekukonko/tablewriter"
	"agenda"
)

// searchCmd represents the register command
var searchCmd = &cobra.Command{
	Use:   "register",
	Short: "register for further use",
	Long: `register for further use and u need to input username, password.
	it will be better if email and phone is provider`,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Errorf("Error[register]： %v\n", err)
			}
		}()

		fmt.Println("register called")

		userFlagBool, _ := cmd.Flags().GetBool("user")
		meetingFlagBool, _ := cmd.Flags().GetString("meeting")
		
		if(userFlagBool && meetingFlagBool) {
			panic(errors.New("The flag meeting(m) and user(u) can't appear at the same time"))
		} else if(userFlagBool) {
			// usernameItem := "username"
			// emailItem := "email"
			// phoneItem := "phone"
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Email", "Phone"})

			userInfoList := agenda.QueryAccountAll()
			for _, userInfo range userInfoList {
				data := []string{userInfo.username, userInfo.email, userInfo.phone}
				table.append(data)
			}
			table.Render()
		} else if{
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Title", "Sponsor", "startTime", "EndTime", "Participators"})

			 agenda.QueryMeetingAll
			for _, data := range data {
				table.Append()
			}
		}

		// fmt.Println("search called by " + userFlagBool)
		// fmt.Println("register with info password: " + meetingFlagBool)

		info := agenda.MakeUserInfo(agenda.Username(username), agenda.Auth(password), email, phone)

		if err := agenda.RegisterUser(info); err != nil {
			panic(err)
		} else {
			fmt.Print("register sucessfully!\n")
			agenda.SaveAll()
		}
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	searchCmd.Flags().BoolP("user", "u", false, "search users")
	searchCmd.Flags().BoolP("meeting", "m", false, "search meeting")

}
