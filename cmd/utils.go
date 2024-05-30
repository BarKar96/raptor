package cmd

import (
	"os"
	"runtime"
	"strconv"
)

const requestedCPU = "REQUESTED_CPU"

func SetGoMaxProcs() {
	if procs := os.Getenv(requestedCPU); len(procs) > 0 {
		maxProcs, err := strconv.Atoi(procs)
		if err != nil {
			panic(err)
		}
		runtime.GOMAXPROCS(maxProcs * 2)

	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())
		os.Setenv("GOMAXPROCS", strconv.Itoa(runtime.NumCPU()))
	}
}
