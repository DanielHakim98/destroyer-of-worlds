/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/DanielHakim98/destroyer-of-worlds/core"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var n int
var rootCmd = &cobra.Command{
	Use:   "destroyer-of-worlds",
	Short: "Let's destroy the worlds",
	Long: `Well actually, it's just a fancy name for load tester.

	I mean, have you ever think about break one's server by DDosing it because you kinda hate those guy or that corporate.
	Well I have great news for you, This is the tool you can use to accomplish your goal`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("No 'url' is given")
			os.Exit(1)
		}

		url := args[0]
		loadTester := core.NewFetcher(url, n)
		loadTester.Run()
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

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.destroyer-of-worlds.yaml)")
	rootCmd.PersistentFlags().IntVarP(&n, "requests", "n", 1, "The total requests to be sent. Default is 1")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
