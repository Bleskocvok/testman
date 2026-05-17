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

func main() {
	author := flag.String("email", "", "Email of author")
	day := flag.String("day", "", "Optional time. Can be YYYY-MM-DD, today, yesterday, week, month.")

	flag.Parse()

	if *author == "" {
		*author = os.Getenv("GITLIST_AUTHOR")
	}

	if *day == "" {
		*day = "today"
	}

	if *day == "" {
		*day = "today"
	}

	if *day == "today" {
		fmt.Println("SADaposkpjdsakldjhasdklj")
		*day = time.Now().Format("2006-01-02")
	}

	if *day == "yesterday" {
		*day = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	fmt.Println("Author: " + *author)
	fmt.Println("Day: " + *day)

	changes, err := getChanges(*author, *day)

	if err != nil {
		log.Fatalln(err)
	}

	totalAdded, totalRemoved := 0, 0
	for _, ch := range changes {
		totalAdded += ch.Added
		totalRemoved += ch.Removed
	}
	fmt.Println("+", totalAdded, "-", totalRemoved, "=", totalAdded - totalRemoved)
}
