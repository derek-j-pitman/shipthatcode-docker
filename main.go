package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PID Namespace Simulator
// Each namespace gets its own PID counter starting at 1.
// Commands:
//   NEWNS  -> allocate next ns id, init empty process table, print id
//   FORK <ns> <name>  -> spawn process; assign in-ns pid; print pid
//   EXIT <ns> <pid>  -> mark exited; print OK
//   PS <ns>  -> print '<pid> <name> <state>' sorted by pid

type Process struct {
	Name    string
	Running bool
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var out []string
	// TODO: declare your state structures here
	var nsTable [][]Process
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "NEWNS":
			// TODO: allocate next ns id, init empty process table, print id
			nsTable = append(nsTable, []Process{})
			out = append(out, strconv.Itoa(len(nsTable)))
		case "FORK":
			// TODO: spawn process; assign in-ns pid; print pid
			nsId, _ := strconv.Atoi(parts[1])
			nsTable[nsId-1] = append(nsTable[nsId-1], Process{parts[2], true})
			out = append(out, strconv.Itoa(len(nsTable[nsId-1])))
		case "EXIT":
			// TODO: mark exited; print OK
			nsId, _ := strconv.Atoi(parts[1])
			pid, _ := strconv.Atoi(parts[2])
			nsTable[nsId-1][pid-1].Running = false
			out = append(out, "OK")
		case "PS":
			// TODO: print '<pid> <name> <state>' sorted by pid
			nsId, _ := strconv.Atoi(parts[1])
			var processes []string
			for i, v := range nsTable[nsId-1] {
				var isRunning string
				if v.Running {
					isRunning = "running"
				} else {
					isRunning = "exited"
				}
				processes = append(processes, fmt.Sprintf("%d %s %s", i+1, v.Name, isRunning))
			}
			out = append(out, strings.Join(processes, "\n"))
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
