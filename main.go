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

type IdKey struct {
	Namespace int
	Key       string
}

type IPC struct {
	Mem    map[string]int
	Queues map[string][]string
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var out []string
	// TODO: declare your state structures here
	var namespaces []IPC
	var memIds []IdKey
	var queueIds []IdKey
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "NEWNS":
			// TODO: new ns; print id
			namespaces = append(namespaces, IPC{map[string]int{}, map[string][]string{}})
			out = append(out, strconv.Itoa(len(namespaces)))
		case "SHMGET":
			// TODO: OK <id> or EEXIST
			nsId, _ := strconv.Atoi(parts[1])
			if _, exists := namespaces[nsId-1].Mem[parts[2]]; exists {
				out = append(out, "EEXIST")
			} else {
				size, _ := strconv.Atoi(parts[3])
				namespaces[nsId-1].Mem[parts[2]] = size
				memIds = append(memIds, IdKey{nsId - 1, parts[2]})
				out = append(out, fmt.Sprintf("OK %d", len(memIds)))
			}
		case "SHMLIST":
			// TODO: '<id> <key> <size>' sorted by id
			nsId, _ := strconv.Atoi(parts[1])
			var list []string
			for i, k := range memIds {
				if k.Namespace == nsId-1 {
					list = append(list, fmt.Sprintf("%d %s %d", i+1, k.Key, namespaces[nsId-1].Mem[k.Key]))
				}
			}
			out = append(out, strings.Join(list, "\n"))
		case "MSGGET":
			// TODO: OK <id> or EEXIST
			nsId, _ := strconv.Atoi(parts[1])
			if _, exists := namespaces[nsId-1].Queues[parts[2]]; exists {
				out = append(out, "EEXIST")
			} else {
				namespaces[nsId-1].Queues[parts[2]] = []string{}
				queueIds = append(queueIds, IdKey{nsId - 1, parts[2]})
				out = append(out, fmt.Sprintf("OK %d", len(queueIds)))
			}
		case "MSGSEND":
			// TODO: enqueue; OK
			nsId, _ := strconv.Atoi(parts[1])
			qId, _ := strconv.Atoi(parts[2])
			queueKey := queueIds[qId-1].Key
			namespaces[nsId-1].Queues[queueKey] = append(namespaces[nsId-1].Queues[queueKey], parts[3])
			out = append(out, "OK")
		case "MSGRECV":
			// TODO: dequeue head; print msg or EMPTY
			nsId, _ := strconv.Atoi(parts[1])
			qId, _ := strconv.Atoi(parts[2])
			queueKey := queueIds[qId-1].Key
			if len(namespaces[nsId-1].Queues[queueKey]) > 0 {
				out = append(out, namespaces[nsId-1].Queues[queueKey][0])
				namespaces[nsId-1].Queues[queueKey] = namespaces[nsId-1].Queues[queueKey][1:]
			} else {
				out = append(out, "EMPTY")
			}
		case "PEEK":
			// TODO: 1 if visible in ns_b, else 0 (0 if ns_a != ns_b)
			ns1, _ := strconv.Atoi(parts[1])
			ns2, _ := strconv.Atoi(parts[2])
			if ns1 != ns2 {
				out = append(out, "0")
			} else {
				if _, found := namespaces[ns1-1].Mem[parts[3]]; found {
					out = append(out, "1")
				} else {
					out = append(out, "0")
				}
			}
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
