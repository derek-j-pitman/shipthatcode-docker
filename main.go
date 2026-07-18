package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// IPC Namespace Isolation
// Per-ns SysV shm and message queues. Ids are globally unique, keys are per-ns.
// Commands:
//   NEWNS  -> new ns; print id
//   SHMGET <ns> <key> <size>  -> OK <id> or EEXIST
//   SHMLIST <ns>  -> '<id> <key> <size>' sorted by id
//   MSGGET <ns> <key>  -> OK <id> or EEXIST
//   MSGSEND <ns> <qid> <msg>  -> enqueue; OK
//   MSGRECV <ns> <qid>  -> dequeue head; print msg or EMPTY
//   PEEK <ns_a> <ns_b> <key>  -> 1 if visible in ns_b, else 0 (0 if ns_a != ns_b)

type MsgQueue struct {
	Key      string
	Messages []string
}

type SHM struct {
	Key  string
	Size int
}

type IPC struct {
	Mem    []SHM
	Queues []MsgQueue
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var out []string
	// TODO: declare your state structures here
	var namespaces []IPC
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "NEWNS":
			// TODO: new ns; print id
			namespaces = append(namespaces, IPC{[]SHM{}, []MsgQueue{}})
			out = append(out, strconv.Itoa(len(namespaces)))
		case "SHMGET":
			// TODO: OK <id> or EEXIST
			nsId, _ := strconv.Atoi(parts[1])
			found := false
			for _, k := range namespaces[nsId-1].Mem {
				if k.Key == parts[2] {
					found = true
					break
				}
			}
			if found {
				out = append(out, "EEXIST")
			} else {
				size, _ := strconv.Atoi(parts[3])
				namespaces[nsId-1].Mem = append(namespaces[nsId-1].Mem, SHM{parts[2], size})
				out = append(out, fmt.Sprintf("OK %d", len(namespaces[nsId-1].Mem)))
			}
		case "SHMLIST":
			// TODO: '<id> <key> <size>' sorted by id
			nsId, _ := strconv.Atoi(parts[1])
			var list []string
			for i, k := range namespaces[nsId-1].Mem {
				list = append(list, fmt.Sprintf("%d %s %d", i+1, k.Key, k.Size))
			}
			out = append(out, strings.Join(list, "\n"))
		case "MSGGET":
			// TODO: OK <id> or EEXIST
			nsId, _ := strconv.Atoi(parts[1])
			found := false
			for _, k := range namespaces[nsId-1].Queues {
				if k.Key == parts[2] {
					found = true
					break
				}
			}
			if found {
				out = append(out, "EEXIST")
			} else {
				namespaces[nsId-1].Queues = append(namespaces[nsId-1].Queues, MsgQueue{parts[2], []string{}})
				out = append(out, fmt.Sprintf("OK %d", len(namespaces[nsId-1].Queues)))
			}
		case "MSGSEND":
			// TODO: enqueue; OK
			nsId, _ := strconv.Atoi(parts[1])
			qId, _ := strconv.Atoi(parts[2])
			namespaces[nsId-1].Queues[qId].Messages = append(namespaces[nsId-1].Queues[qId].Messages, parts[3])
			out = append(out, "OK")
		case "MSGRECV":
			// TODO: dequeue head; print msg or EMPTY
			nsId, _ := strconv.Atoi(parts[1])
			qId, _ := strconv.Atoi(parts[2])
			out = append(out, namespaces[nsId-1].Queues[qId].Messages[0])
			namespaces[nsId-1].Queues[qId].Messages = namespaces[nsId-1].Queues[qId].Messages[1:]
		case "PEEK":
			// TODO: 1 if visible in ns_b, else 0 (0 if ns_a != ns_b)
			ns1, _ := strconv.Atoi(parts[1])
			ns2, _ := strconv.Atoi(parts[2])
			if ns1 != ns2 {
				out = append(out, "0")
			} else {
				found := false
				for _, k := range namespaces[ns1-1].Mem {
					if k.Key == parts[2] {
						found = true
						break
					}
				}
				if found {
					out = append(out, "1")
				} else {
					out = append(out, "0")
				}
			}
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
