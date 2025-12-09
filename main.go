/*
Copyright © 2025 Kenny
*/
package main

import (
	"os"

	"Door_System_User_Automation/cmd"
	"Door_System_User_Automation/utils"
)

func main() {
	cmd.Execute()

	wisDF := utils.CreateDataFrame(cmd.WisFile)
	s2DF := utils.CreateDataFrame(cmd.S2File)

	utils.MergeAndExport(s2DF, wisDF)

	err := os.RemoveAll("temp")
	if err != nil {
		panic(err)
	}
}
