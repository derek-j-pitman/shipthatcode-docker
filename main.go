package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// User Namespace Translator
// Bidirectional uid_map: in-ns id <-> host id, via ranges (in, host, length).
// Commands:
//   NEWNS  -> create new userns; print id (1,2,...)
//   MAP <ns> <in_id> <host_id> <length>  -> append range; OK or 'ERR overlap'
//   TRANSLATE <ns> <in_id>  -> print host uid or 'unmapped'
//   WHO <ns> <host_id>  -> print in-ns uid or 'unmapped'

type UidRange struct {
	InStart, HostStart, Length int
}

func isOverlap(lower int, lowerLen int, upper int, upperLen int) bool {
	return upper < lower+lowerLen || lower > upper+upperLen
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var out []string
	// TODO: declare your state structures here
	var namespaces [][]UidRange
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "NEWNS":
			// TODO: create new userns; print id (1,2,...)
			namespaces = append(namespaces, []UidRange{})
			out = append(out, strconv.Itoa(len(namespaces)))
		case "MAP":
			// TODO: append range; OK or 'ERR overlap'
			// For a new interval, overlaps exist if:
			// new id <= existing id + length
			// new id + length >= existing id
			// This could happen for both in-ns *and* host ranges
			nsId, _ := strconv.Atoi(parts[1])
			inId, _ := strconv.Atoi(parts[2])
			hostId, _ := strconv.Atoi(parts[3])
			length, _ := strconv.Atoi(parts[4])
			newRange := UidRange{inId, hostId, length}
			overlapFound := false
			for _, r := range namespaces[nsId-1] {
				var lo, loLen, hi, hiLen int
				if r.InStart < newRange.InStart {
					lo = r.InStart
					loLen = r.Length
					hi = newRange.InStart
					hiLen = newRange.Length
				} else {
					lo = newRange.InStart
					loLen = newRange.Length
					hi = r.InStart
					hiLen = r.Length
				}
				overlapFound = overlapFound || isOverlap(lo, loLen, hi, hiLen)
				if r.HostStart < newRange.HostStart {
					lo = r.HostStart
					loLen = r.Length
					hi = newRange.HostStart
					hiLen = newRange.Length
				} else {
					lo = newRange.HostStart
					loLen = newRange.Length
					hi = r.HostStart
					hiLen = r.Length
				}
				overlapFound = overlapFound || isOverlap(lo, loLen, hi, hiLen)
			}
			if overlapFound {
				out = append(out, "ERR overlap")
			} else {
				namespaces[nsId-1] = append(namespaces[nsId-1], newRange)
				out = append(out, "OK")
			}
		case "TRANSLATE":
			// TODO: print host uid or 'unmapped'
			nsId, _ := strconv.Atoi(parts[1])
			inId, _ := strconv.Atoi(parts[2])
			res := "unmapped"
			for _, r := range namespaces[nsId-1] {
				if inId >= r.InStart && inId < r.InStart+r.Length {
					res = strconv.Itoa(r.HostStart + inId)
					break
				}
			}
			out = append(out, res)
		case "WHO":
			// TODO: print in-ns uid or 'unmapped'
			nsId, _ := strconv.Atoi(parts[1])
			hostId, _ := strconv.Atoi(parts[2])
			res := "unmapped"
			for _, r := range namespaces[nsId-1] {
				if hostId >= r.HostStart && hostId < r.HostStart+r.Length {
					res = strconv.Itoa(hostId - r.HostStart)
					break
				}
			}
			out = append(out, res)
		}
	}
	fmt.Println(strings.Join(out, "\n"))
}
