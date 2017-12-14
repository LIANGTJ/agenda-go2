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
	"agendaHttp"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	// "time"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query for further use",
	Long: `query for further use and u need to input username, password.
	it will be better if email and phone is provider`,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Errorf("Error[query]： %v\n", err)
			}
		}()

		fmt.Println("query called")

		userFlagBool, _ := cmd.Flags().GetBool("user")
		meetingFlagBool, _ := cmd.Flags().GetBool("meeting")
		// startTimeFlag, _ := cmd.Flags().GetString("startTime")
		// endTimeFlag, _ := cmd.Flags().GetString("endTime")

		// startTime, _ := time.Parse("2006-01-02 15:04:05", startTimeFlag)
		// endTime, _ := time.Parse("2006-01-02 15:04:05", endTimeFlag)

		if userFlagBool && meetingFlagBool {

			panic(errors.New("The flag meeting(m) and user(u) can't appear at the same time"))

		} else if userFlagBool {

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Email", "Phone"})
			userList, err:= agendaHttp.QueryAccountAll()
			if err != nil {
				panic(err)
			}
			for _, user := range *userList {
				data := []string{user.Username, user.Email, user.Phone}
				table.Append(data)
			}
			table.Render()

		} else {
			fmt.Println("[queryMeetingList]: TODO")
			// table := tablewriter.NewWriter(os.Stdout)
			// table.SetHeader([]string{"Title", "Sponsor", "startTime", "EndTime", "Participators"})
			// currentUsername := agenda.LoginedUser().Name
			// meetingList := agenda.QueryMeetingByInterval(startTime, endTime, currentUsername)
			// for _, meeting := range meetingList {
			// 	participators := ""
			// 	for _, participator := range meeting.Participators {
			// 		participators = participators + " " + string(participator)
			// 	}
			// 	table.Append([]string{string(meeting.Title), string(meeting.Sponsor),
			// 		meeting.StartTime, meeting.EndTime,
			// 		participators})

			// 	table.Render()
			// }
		}

		// info := agenda.MakeUserInfo(agenda.Username(username), agenda.Auth(password), email, phone)

		// if err := agenda.RegisterUser(info); err != nil {
		// 	panic(err)
		// } else {
		// 	fmt.Print("query sucessfully!\n")
		// 	agenda.SaveAll()
		// }
	},
}

func init() {
	RootCmd.AddCommand(queryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	queryCmd.Flags().BoolP("user", "u", false, "query users")
	queryCmd.Flags().BoolP("meeting", "m", false, "query meeting")
	queryCmd.Flags().StringP("startTime", "s", "", "query startTime")
	queryCmd.Flags().StringP("endTime", "e", "", "query endTime")

}
