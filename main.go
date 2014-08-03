package main

import (
	"encoding/csv"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/kennygrant/sanitize"
)

var regexpEmail = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

// main function that initialize the program
func main() {
	// log an error if there is less than 2 args
	// the first is the name of program
	if len(os.Args) < 2 {
		log.Fatal("You need to set the filename\nexample: ./clean_csv sample.csv")
	}

	// set fileName from args of program
	fileName := os.Args[1]

	lines := ReadCsv(fileName)

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

		email := ClearEmail(line[0])
		line[0] = email
		if existEmails[email] {
			duplicatedEmails = append(duplicatedEmails, line)
		} else {
			if ValidEmail(email) {
				validEmails = append(validEmails, line)
			} else {
				invalidEmails = append(invalidEmails, line)
			}
			existEmails[email] = true
		}
	}

	path := path.Dir(os.Args[1]) + "/output"

	CreateFolder(path)

	var wg sync.WaitGroup
	wg.Add(3)

	go writeToCsv(headers, validEmails, path+"/valid.csv", &wg)
	go writeToCsv(headers, invalidEmails, path+"/invalid.csv", &wg)
	go writeToCsv(headers, duplicatedEmails, path+"/duplicated.csv", &wg)

	wg.Wait()
}

// ReadCsv is a function that open the file,
// read al the content and return it, if an error
// occurs it's log it and close the program
func ReadCsv(fileName string) [][]string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.TrimLeadingSpace = true

	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading all lines: ", err)
	}
	return lines
}

// ClearEmail is a function that remove spaces,
// transform in lower case and remove accents
// of string that should be an email
func ClearEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	email = sanitize.Accents(email)
	return email
}

// ValidEmail is a function that return is an email
// is valid or not
func ValidEmail(email string) bool {
	match, _ := regexp.MatchString(regexpEmail, email)
	return match
}

// CreateFolder is a function to create a
// folder with spefic permission if an error
// occurs it's log it and close the program
func CreateFolder(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatal("Error to create", path, err)
	}
}

func writeToCsv(headers []string, lines [][]string, fileName string, wg *sync.WaitGroup) {
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
	log.Printf("%s emails => %d\n", kind, len(lines))
	wg.Done()
}

type sortByEmail [][]string

func (p sortByEmail) Len() int           { return len(p) }
func (p sortByEmail) Less(i, j int) bool { return p[i][0] < p[j][0] }
func (p sortByEmail) Swap(i, j int)      { p[i][0], p[j][0] = p[j][0], p[i][0] }
