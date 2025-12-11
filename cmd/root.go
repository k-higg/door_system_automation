/*
Copyright © 2025 Kenny
*/
package cmd

import (
	"Door_System_User_Automation/utils"
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

var mobileCmd = &cobra.Command{
	Use:   "mobile",
	Short: "Add user's mobile number and email, and then enable mobile credentials",
	Long: `Add user's mobile number and email, and then enable mobile credentials. 
		Outputs two csvs that ned to be imported in the correct order so S2 can properly enable the credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		wisDF := utils.CreateDataFrame(WisFile)
		utils.NormalizeUID(wisDF)
		s2DF := utils.CreateDataFrame(S2File)

		utils.MergeMobile(s2DF, wisDF)
	},
}

var genderCmd = &cobra.Command{
	Use:   "gender",
	Short: "Add user's gender and updates their access levels based on that",
	Long: `Add user's gender and updates their access levels based on that... 
		Males will get access to a specified building for laundry access, 
		Females will get access to a different specified building for laundry access.
		Outputs one file to be imported to S2.`,
	Run: func(cmd *cobra.Command, args []string) {
		wisDF := utils.CreateDataFrame(WisFile)
		utils.NormalizeUID(wisDF)
		s2DF := utils.CreateDataFrame(S2File)

		utils.MergeGender(s2DF, wisDF)
	},
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
	rootCmd.AddCommand(mobileCmd)
	rootCmd.AddCommand(genderCmd)
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Door_System_User_Automation.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("help", "h", false, "Help message for help")

	mobileCmd.Flags().StringVarP(&WisFile, "wis", "w", "", "Path to WIS file (CSV)")
	mobileCmd.Flags().StringVarP(&S2File, "s2", "s", "", "Path to S2 file (CSV)")

	genderCmd.Flags().StringVarP(&WisFile, "wis", "w", "", "Path to WIS file (CSV)")
	genderCmd.Flags().StringVarP(&S2File, "s2", "s", "", "Path to S2 file (CSV)")
}
