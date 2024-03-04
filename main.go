/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/DanielHakim98/destroyer-of-worlds/cmd"
)

func main() {
	f, err := os.Create("destroyer-of-worlds.pprof")
	if err != nil {
		log.Fatal(f)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	cmd.Execute()
}
