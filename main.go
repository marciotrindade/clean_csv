package main

import (
	"clean_csv/go-humanize"
	"clean_csv/sanitize"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const REGEXP_EMAIL = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You need to set the filename\nexample: ./clean_csv sample.csv")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.TrimLeadingSpace = true
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading all lines: %v", err)
		return
	}

	var headers []string
	var validEmails [][]string
	var invalidEmails [][]string
	var duplicatedEmails [][]string
	existEmails := make(map[string]bool)

	for i, line := range lines {
		if i == 0 {
			headers = line
			continue
		}

		email := clearEmail(line[0])
		line[0] = email
		if existEmails[email] {
			duplicatedEmails = append(duplicatedEmails, line)
		} else {
			if validEmail(email) {
				validEmails = append(validEmails, line)
			} else {
				invalidEmails = append(invalidEmails, line)
			}
			existEmails[email] = true
		}
	}

	path := path.Dir(os.Args[1]) + "/output"
	mkerr := os.MkdirAll(path, 0755)
	if mkerr != nil {
		fmt.Println("MkdirAll: %s %s", path, mkerr)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	go func() {
		writeToCsv(headers, validEmails, path+"/valid.csv")
		waitGroup.Done()
	}()

	go func() {
		writeToCsv(headers, invalidEmails, path+"/invalid.csv")
		waitGroup.Done()
	}()

	go func() {
		writeToCsv(headers, duplicatedEmails, path+"/duplicated.csv")
		waitGroup.Done()
	}()

	waitGroup.Wait()
}

func clearEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	email = sanitize.Accents(email)
	return email
}

func validEmail(email string) bool {
	match, _ := regexp.MatchString(REGEXP_EMAIL, email)
	return match
}

type sortByEmail [][]string

func (p sortByEmail) Len() int           { return len(p) }
func (p sortByEmail) Less(i, j int) bool { return p[i][0] < p[j][0] }
func (p sortByEmail) Swap(i, j int)      { p[i][0], p[j][0] = p[j][0], p[i][0] }

func writeToCsv(headers []string, lines [][]string, fileName string) {
	file, _ := os.Create(fileName)
	defer file.Close()

	sort.Sort(sortByEmail(lines))

	var rows [][]string
	rows = append(rows, headers)
	rows = append(rows, lines...)

	writer := csv.NewWriter(file)
	writer.WriteAll(rows)
	writer.Flush()

	file.Close()

	kind := strings.Replace(path.Base(fileName), path.Ext(fileName), "", -1)
	fmt.Println(kind + " emails => " + humanize.Comma(int64(len(lines))))
}
