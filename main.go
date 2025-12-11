/*
Copyright © 2025 Kenny
*/
package main

import (
	"fmt"
	"os"

	"Door_System_User_Automation/cmd"
)

func main() {
	cmd.Execute()

	err := os.RemoveAll("temp")
	if err != nil {
		fmt.Println("Error removing temp folder..")
	}
}
