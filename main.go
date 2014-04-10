package main

import (
  "fmt"
  "os"
)

func main() {
  finfo, _ := os.Stat("output")
  fmt.Println(finfo.Mode())
}
