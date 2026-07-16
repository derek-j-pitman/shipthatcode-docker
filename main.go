package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PID Namespace Simulator
// Each namespace gets its own PID counter starting at 1.
// Commands:
//   NEWNS  -> allocate next ns id, init empty process table, print id
//   FORK <ns> <name>  -> spawn process; assign in-ns pid; print pid
//   EXIT <ns> <pid>  -> mark exited; print OK
//   PS <ns>  -> print '<pid> <name> <state>' sorted by pid

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var out []string
	// TODO: declare your state structures here
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" { continue }
		parts := strings.Fields(line)
		switch parts[0] {
		case "NEWNS":
			// TODO: allocate next ns id, init empty process table, print id
		case "FORK":
			// TODO: spawn process; assign in-ns pid; print pid
		case "EXIT":
			// TODO: mark exited; print OK
		case "PS":
			// TODO: print '<pid> <name> <state>' sorted by pid
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
