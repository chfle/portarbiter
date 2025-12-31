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

	seen := make(map[string]bool)

	for _, pid := range pids {

		// 1) Docker / Docker-Compose
		if dock, ok, err := resolve.ResolveDocker(pid); err == nil && ok {
			key := dock.Type().String() + ":" + dock.ID()
			if !seen[key] {
				fmt.Println(" ", dock.Describe())
				seen[key] = true
			}
			continue
		}

		// 2) systemd
		if sys, ok, err := resolve.ResolveSystemd(pid); err == nil && ok {
			key := "systemd:" + sys.ID()
			if !seen[key] {
				fmt.Println(" ", sys.Describe())
				seen[key] = true
			}
			continue
		}

		// 3) raw process
		proc, err := resolve.ResolveProcess(pid)
		if err != nil {
			fmt.Printf("  pid=%d error=%v\n", pid, err)
			continue
		}

		key := "proc:" + proc.ID()
		if !seen[key] {
			fmt.Println(" ", proc.Describe())
			seen[key] = true
		}
	}
}

