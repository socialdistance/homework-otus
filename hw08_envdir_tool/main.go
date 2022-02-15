package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		return
	}

	dir := args[1]
	cmds := args[2:]
	envs, err := ReadDir(dir)
	if err != nil {
		log.Panicf("Error:%v", err)
	}

	run := RunCmd(cmds, envs)

	os.Exit(run)
}
