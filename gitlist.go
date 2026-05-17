package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "2006-01-02"

type Changes struct {
	Added int
	Removed int
	File string
}

func getChanges(author string, day string) ([]Changes, error) {
	cmd := exec.Command("git",
						"log",
						"--all",
						"--author=" + author,
						"--since=" + day + " 00:00:00",
						"--until=" + day + " 23:59:59",
						"--numstat",
						"--pretty=tformat:short")
	output, e := cmd.Output()

	if e != nil {
		return []Changes{}, e
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	changes := make([]Changes, 0)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if len(fields) < 3 {
			continue
		}

		a, err1 := strconv.Atoi(fields[0])
		r, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			continue
		}

		changes = append(changes, Changes{a, r, fields[2]})
	}
	return changes, nil
}

func printNumStat(prefix string, added, removed int) {
	fmt.Printf("%s +%4d -%4d = %4d\n", prefix, added, removed, added - removed)
}

func printDayChanges(author, day string) (int, int) {
	changes, err := getChanges(author, day)

	if err != nil {
		log.Fatalln(err)
	}

	totalAdded, totalRemoved := 0, 0
	for _, ch := range changes {
		totalAdded += ch.Added
		totalRemoved += ch.Removed
	}
// 	fmt.Println("Day:", day, " +", totalAdded, "-", totalRemoved, "=", totalAdded - totalRemoved)
	printNumStat("Day: " + day, totalAdded, totalRemoved)

	return totalAdded, totalRemoved
}

func printLongerPeriod(author string, days int) {
	added, removed := 0, 0
	for i := range days {
		day := time.Now().AddDate(0, 0, -i).Format(DateFormat)
		a, r := printDayChanges(author, day)
		added += a
		removed += r
	}
	printNumStat("Total:", added, removed)
}

func valiDate(str string) bool {
	_, e := time.Parse(DateFormat, str)
	return e == nil
}

func main() {
	author := flag.String("author", "", "Email of author")
	day := flag.String("time", "", "Optional time. Can be YYYY-MM-DD, today, yesterday, week, month.")

	flag.Parse()

	if *author == "" {
		*author = os.Getenv("GITLIST_AUTHOR")
	}

	if *day == "" {
		*day = "today"
	}

	fmt.Println("Author: " + *author)

	switch *day {
		case "today":
			printDayChanges(*author, time.Now().Format(DateFormat))

		case "yesterday":
			printDayChanges(*author, time.Now().AddDate(0, 0, -1).Format(DateFormat))

		case "week":
			printLongerPeriod(*author, 7)

		case "month":
			printLongerPeriod(*author, 30)

		case "year":
			printLongerPeriod(*author, 365)

		default:
			if !valiDate(*day) {
				log.Fatalln("Invalid date: ", *day)
			}
			printDayChanges(*author, *day)
	}
}
