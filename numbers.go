package main

import "fmt"

func main() {
  for i := 1; i <= 10; i++ {
    var text string
    if i % 2 == 0 {
      text = "even"
    } else {
      text = "odd"
    }

    fmt.Println(i, text)
  }

}
