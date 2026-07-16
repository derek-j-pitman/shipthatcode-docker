package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Mount Namespace Simulator
// Each ns has its own mount table; new ns inherits parent ns 0's snapshot.
// Commands:
//   NEWNS  -> create new ns; snapshot ns 0's mounts; print id
//   MOUNT <ns> <src> <dst>  -> add mount; print OK
//   UMOUNT <ns> <dst>  -> remove mount; print OK
//   LISTMOUNTS <ns>  -> print '<src> on <dst>' sorted by dst, or (empty)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var out []string
	// TODO: declare your state structures here
	var mounts []map[string]string
	mounts = append(mounts, make(map[string]string))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "NEWNS":
			// TODO: create new ns; snapshot ns 0's mounts; print id
			newNs := make(map[string]string)
			maps.Copy(newNs, mounts[0])
			mounts = append(mounts, newNs)
			out = append(out, strconv.Itoa(len(mounts)-1))
		case "MOUNT":
			// TODO: add mount; print OK
			nsId, _ := strconv.Atoi(parts[1])
			srcMount := parts[2]
			dstMount := parts[3]
			mounts[nsId][dstMount] = srcMount
			out = append(out, "OK")
		case "UMOUNT":
			// TODO: remove mount; print OK
			nsId, _ := strconv.Atoi(parts[1])
			dstMount := parts[2]
			delete(mounts[nsId], dstMount)
			out = append(out, "OK")
		case "LISTMOUNTS":
			// TODO: print '<src> on <dst>' sorted by dst, or (empty)
			nsId, _ := strconv.Atoi(parts[1])
			if len(mounts[nsId]) == 0 {
				out = append(out, "(empty)")
			} else {
				var points []string
				var keys []string
				for k := range mounts[nsId] {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					points = append(points, fmt.Sprintf("%s on %s", mounts[nsId][k], k))
				}
				out = append(out, strings.Join(points, "\n"))
			}
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
