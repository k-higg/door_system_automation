/*
Copyright © 2025 Kenny
*/
package main

import (
	"os"

	"Door_System_User_Automation/cmd"
)

func main() {
	cmd.Execute()

	err := os.RemoveAll("temp")
	if err != nil {
		panic(err)
	}
}
