/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/DanielHakim98/destroyer-of-worlds/cmd"
)

func main() {
	profilingPath := os.Getenv("PPROF_PATH")
	if profilingPath != "" {
		parent := filepath.Dir(profilingPath)
		if _, err := os.Stat(parent); os.IsNotExist(err) {
			err := os.MkdirAll(parent, os.ModePerm)
			if err != nil {
				fmt.Println("Error: Failed to create directory from PPROF_PATH")
				os.Exit(1)
			}
		}

		f, err := os.Create(profilingPath)
		if err != nil {
			fmt.Println("Error: Failed to create file from PPROF_PATH")
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	cmd.Execute()
}
