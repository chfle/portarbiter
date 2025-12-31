package main

import (
	"fmt"
	"os"
	"strconv"
	"portarbiter/internal/resolve"
	"portarbiter/internal/detect"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: portarbiter <port>")
		os.Exit(1)
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid port")
		os.Exit(1)
	}

	pids, err := detect.FindPIDsByPort(port)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for _, pid := range pids {
		owner, err := resolve.ResolveProcess(pid)
		if err != nil {
			fmt.Println("PID", pid, "error:", err)
			continue
		}
		fmt.Println(owner.Describe())
	}

}

