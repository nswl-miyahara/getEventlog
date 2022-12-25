package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var r = regexp.MustCompile(" +")
var r2 = regexp.MustCompile(":[0-9]+$")

var cmd = `Get-EventLog System -After (Get-Date).AddDays(-7) | ? { $_.InstanceId -in (6001,6002,7001,7002) } | select @{n='Message'; e={if (($_ | select -ExpandProperty InstanceId) % 2 -eq 1) {"Logon"} else {"Logoff"} }},TimeGenerated`

var output = "./output.tsv"

type Event struct {
	logon  string
	logoff string
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	output = os.Getenv("OUTPUT")
}
func main() {
	out, err := exec.Command("powershell", cmd).Output()
	if err != nil {
		fmt.Println(err)
	}
	m := bindMap(string(out))

	file, err := os.Create(output)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t'
	defer writer.Flush()

	keys := getKeys(m)
	sort.Strings(keys)

	writer.Write([]string{"date", "logon", "logoff"})
	for _, key := range keys {
		writer.Write([]string{key, m[key].logon, m[key].logoff})
	}
}

func getKeys(m map[string]Event) []string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func bindMap(out string) map[string]Event {
	logs := strings.Split(out, "\r\n")
	m := map[string]Event{}
	for i, log := range logs {
		if i < 3 || log == "" {
			continue
		}

		columns := strings.Split(r.ReplaceAllString(log, " "), " ")
		time := r2.ReplaceAllString(columns[2], "")

		var event Event
		if _, ok := m[columns[1]]; !ok {
			if columns[0] == "Logon" {
				event = Event{time, ""}
			} else {
				event = Event{"", time}
			}
			m[columns[1]] = event
			continue
		}

		if columns[0] == "Logon" {
			v := m[columns[1]]
			if v.logon != "" {
				old, _ := strconv.Atoi(strings.Replace(v.logon, ":", "", -1))
				new, _ := strconv.Atoi(strings.Replace(time, ":", "", -1))
				if old < new {
					continue
				}
			}
			v.logon = time
			m[columns[1]] = v
		} else {
			v := m[columns[1]]
			if v.logoff != "" {
				old, _ := strconv.Atoi(strings.Replace(v.logoff, ":", "", -1))
				new, _ := strconv.Atoi(strings.Replace(time, ":", "", -1))
				if old > new {
					continue
				}
			}
			v.logoff = time
			m[columns[1]] = v
		}
	}
	return m
}
