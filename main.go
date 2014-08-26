package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

type AppCounters struct {
	Name  []byte
	Pid   int
	User  uint64
	Sys   uint64
	Vsize uint64
	Rss   uint64
}

func (a AppCounters) String() string {
	return fmt.Sprintf("%s[%d]:\n  User = %d, Sys = %d\n  Vsize = %d, RSS = %d",
		string(a.Name[:]), a.Pid, a.User, a.Sys, a.Vsize, a.Rss)
}

func main() {
	fileInfos, err := ioutil.ReadDir("/proc/")
	if err != nil {
		log.Fatal(err)
	}

	for _, info := range fileInfos {
		if info.IsDir() {
			var pid int
			_, err := fmt.Sscanf(info.Name(), "%d", &pid)
			if err != nil {
				break
			}

			readPidStats(pid)
		}
	}
}

func readPidStats(pid int) {
	var counters AppCounters

	if buf, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/statm", pid)); err == nil {
		fmt.Sscanf(string(buf), "%d %d", &counters.Vsize, &counters.Rss)
	}

	if buf, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid)); err == nil {
		var (
			i int
			c byte
		)
		fmt.Sscanf(string(buf), "%d %s %c %d %d %d %d %d %d %d %d %d %d %d %d",
			&i, // placeholder
			&counters.Name,

			&c, &i, &i, &i, &i, &i, &i, &i, &i, &i, &i, // placeholders

			&counters.User, &counters.Sys)
	}

	if len(counters.Name) > 1 {
		counters.Name = counters.Name[1 : len(counters.Name)-1]
	}
	counters.Pid = pid
	fmt.Println(counters)
}
