package main

import (
  "clean_csv/go-humanize"
  "clean_csv/sanitize"
  "encoding/csv"
  "fmt"
  "io"
  "os"
  "path"
  "regexp"
  "sort"
  "strings"
)

const REGEXP_EMAIL = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

func main() {
  file, _ := os.Open("sample.csv")
  defer file.Close()

  reader := csv.NewReader(file)
  reader.Comma = ','

  var validEmails []string
  var invalidEmails []string
  var duplicatedEmails []string
  existEmails := make(map[string]bool)

  for {
    // read just one record, but we could ReadAll() as well
    record, err := reader.Read()
    // end-of-file is fitted into err
    if err == io.EOF {
      break
    } else if err != nil {
      fmt.Println("Error:", err)
      return
    }
    // record is an array of string so is directly printable
    email := clearEmail(record[0])
    if existEmails[email] {
      duplicatedEmails = append(duplicatedEmails, email)
    } else {
      if validEmail(email) {
        validEmails = append(validEmails, email)
      } else {
        invalidEmails = append(invalidEmails, email)
      }
      existEmails[email] = true
    }
  }

  path := "output"
  mkerr := os.MkdirAll(path, 0755)
  if mkerr != nil {
    fmt.Println("MkdirAll: %s %s", path, mkerr)
  }

  cvalid      := make(chan int)
  cinvalid    := make(chan int)
  cduplicated := make(chan int)

  go func() {
    writeToCsv(validEmails, path + "/valid.csv")
    cvalid <- 1
  }()

  go func() {
    writeToCsv(invalidEmails, path + "/invalid.csv")
    cinvalid <- 1
  }()

  go func() {
    writeToCsv(duplicatedEmails, path + "/duplicated.csv")
    cduplicated <- 1
  }()

  <-cvalid
  <-cinvalid
  <-cduplicated
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

func writeToCsv(emails []string, fileName string) {
  file, _ := os.Create(fileName)
  defer file.Close()


  sort.Strings(emails)

  writer := csv.NewWriter(file)
  for _, email := range emails {
    writer.Write([]string{email})
    writer.Flush()
  }

  file.Close()

  kind := strings.Replace(path.Base(fileName), path.Ext(fileName), "", -1)
  fmt.Println(kind + " emails => " + humanize.Comma(int64(len(emails))))
}
