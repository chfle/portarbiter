package main

import (
	"fmt"
	"os"
	"strconv"

	"portarbiter/internal/detect"
	"portarbiter/internal/resolve"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: portarbiter <port>")
		os.Exit(1)
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil || port <= 0 || port > 65535 {
		fmt.Println("Invalid port")
		os.Exit(1)
	}

	pids, err := detect.FindPIDsByPort(port)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Port %d is used by:\n", port)

	for _, pid := range pids {

		// 1) systemd has priority
		if sys, ok, err := resolve.ResolveSystemd(pid); err == nil && ok {
			fmt.Println(" ", sys.Describe())
			continue
		}

		// 2) fallback: raw process
		proc, err := resolve.ResolveProcess(pid)
		if err != nil {
			fmt.Printf("  pid=%d error=%v\n", pid, err)
			continue
		}

		fmt.Println(" ", proc.Describe())
	}
}

