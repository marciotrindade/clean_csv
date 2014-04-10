package main

import (
  "encoding/csv"
  "fmt"
  "golang-book/go-humanize"
  "golang-book/sanitize"
  "io"
  "os"
  "regexp"
  "sort"
  "strings"
)

const REGEXP_EMAIL = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`

func ClearEmail(email string) string {
  email = strings.TrimSpace(email)
  email = strings.ToLower(email)
  email = sanitize.Accents(email)
  return email
}

func ValidEmail(email string) bool {
  match, _ := regexp.MatchString(REGEXP_EMAIL, email)
  return match
}

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
    email := ClearEmail(record[0])
    if existEmails[email] {
      duplicatedEmails = append(duplicatedEmails, email)
    } else {
      if ValidEmail(email) {
        validEmails = append(validEmails, email)
      } else {
        invalidEmails = append(invalidEmails, email)
      }
      existEmails[email] = true
    }
  }

  sort.Strings(validEmails)
  sort.Strings(invalidEmails)
  sort.Strings(duplicatedEmails)

  path := "output"
  mkerr := os.MkdirAll(path, 0755)
  if mkerr != nil {
    fmt.Println("MkdirAll: %s %s", path, mkerr)
  }

  WriteToCsv(validEmails, path + "/valid.csv")
  WriteToCsv(invalidEmails, path + "/invalid.csv")
  WriteToCsv(duplicatedEmails, path + "/duplicated.csv")

  countValid := int64(len(validEmails))
  countInvalid := int64(len(invalidEmails))
  countDuplicated := int64(len(duplicatedEmails))

  fmt.Println("valid emails", humanize.Comma(countValid))
  fmt.Println("invalid emails", humanize.Comma(countInvalid))
  fmt.Println("duplicated emails", humanize.Comma(countDuplicated))
}

func WriteToCsv(emails []string, fileName string) {
  file, _ := os.Create(fileName)
  defer file.Close()

  writer := csv.NewWriter(file)
  for _, email := range emails {
    writer.Write([]string{email})
    writer.Flush()
  }

  file.Close()
}

