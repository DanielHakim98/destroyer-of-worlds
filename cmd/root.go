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
var (
	_number     int
	_url        string
	_concurrent int
)

var rootCmd = &cobra.Command{
	Use:   "destroyer-of-worlds",
	Short: "Let's destroy the worlds",
	Long: `Well actually, it's just a fancy name for load tester.

	I mean, have you ever think about break one's server by DDosing it because you kinda hate those guy or that corporate.
	Well I have great news for you, This is the tool you can use to accomplish your goal`,
	Run: func(cmd *cobra.Command, args []string) {
		if _url == "" {
			fmt.Fprintln(os.Stderr, cmd.Help())
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Error: URL not provided. Please specify a URL using the --url flag.")
			os.Exit(1)
		}

		if _concurrent <= 0 {
			fmt.Fprintln(os.Stderr, cmd.Help())
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Error: Concurrent worker cannot be zero or negative. Please specify value more than 0")
			os.Exit(1)
		}

		loadTester := core.NewFetcher(_url, _number, _concurrent)
		loadTester.Run()
		loadTester.Summary()
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
	rootCmd.PersistentFlags().IntVarP(&_number, "requests", "n", 1, "The total requests to be sent. Default is 1")
	rootCmd.PersistentFlags().StringVarP(&_url, "url", "u", "", "URL to be tested.")
	rootCmd.PersistentFlags().IntVarP(&_concurrent, "concurrent", "c", 1, "maximum concurent request. Default is 1")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
