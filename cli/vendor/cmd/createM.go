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
	"fmt"
	"status"
	"github.com/spf13/cobra"
	"time"
)

// addCmd represents the create Meeting command
var createMCmd = &cobra.Command{
	Use:   "createM",
	Short: "create Meeting for further use",
	Long: `create Meeting for further use and u need to input username, password.
	it will be better if email and phone is provider`,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Error[create Meeting]： %v\n", err)
			}
		}()
		
		fmt.Println("create Meeting called")

		startTimeFlag, _ := cmd.Flags().GetString("startTime")
		endTimeFlag, _ := cmd.Flags().GetString("endTime")
		participatorsFlag, _ := cmd.Flags().GetStringSlice("participator")
		titleFlag, _ := cmd.Flags().GetString("title")

		startTime, _ := time.Parse("2006-01-02 15:04:05", startTimeFlag)
		endTime, _ := time.Parse("2006-01-02 15:04:05", endTimeFlag)
		participators := make([]string, 0, 0)
		participators = append(participators, status.LoginedUser())
		for _, participator := range participatorsFlag {
			participators = append(participators, participator)
		}


		_, err := agendaHttp.CreateMeeting(titleFlag, participators, startTime, endTime)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("sucessfully create meeting")
		}
		// meetingInfo := agenda.MakeMeetingInfo(agenda.MeetingTitle(titleFlag),
		// 	agenda.LoginedUser().Name, participators,
		// 	startTime, endTime)

		// if _, err := agenda.SponsorMeeting(meetingInfo); err != nil {
		// 	panic(err)
		// }
		// // agenda.SaveAll()
		// fmt.Print("sucessfully create meeting\n")

	},
}

func init() {
	RootCmd.AddCommand(createMCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	createMCmd.Flags().StringP("startTime", "s", "", "create Meeting info for startTime")
	createMCmd.Flags().StringP("endTime", "e", "", "create Meeting info for endTime")
	createMCmd.Flags().StringP("participator", "p", "", "create Meeting info for participator")
	createMCmd.Flags().StringP("title", "t", "", "create Meeting info for title")

}
