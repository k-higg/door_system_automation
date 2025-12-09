/*
Copyright © 2025 Kenny
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Automation",
	Short: "User configuration automation for S2 Lenel",
	Long: `User configuration automation for S2 Lenel, using data exported from WIS and S2 lenel
we can update the correct user by matching the User ID and update any information that is either
missing or is incorrect.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Defaults
var WisFile string = "resources/StudentWorkList.csv"
var S2File string = "resources/people.csv"

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Door_System_User_Automation.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("help", "h", false, "Help message for help")
	rootCmd.Flags().StringVarP(&WisFile, "wis", "w", "", "Path to WIS file (CSV)")
	rootCmd.Flags().StringVarP(&S2File, "s2", "s", "", "Path to S2 file (CSV)")
}
