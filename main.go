package main

import "fmt"
import "time"
import "sync"
import "./scanner"
import "os/exec"

func main() {
  var wg sync.WaitGroup
  ch := make(chan string)
  w := scanner.FileScanner{"./Inbox/", 300.00, make(chan string), wg}
  for {
    go process(w.OutChannel,ch)
    go w.Scan()
    end := <-ch
    if end == "__DONE" {
      fmt.Println("Waiting for next scan")
      time.Sleep(time.Second * 30)
    }
  }
}

func process(c <-chan string, ch chan<- string) {
  for {
    msg := <-c
    fmt.Println("Read", msg)
    if len(msg) > 0 {
      fp := "./Inbox/" + msg
      fmt.Println("Processing " + fp)
      cmd := exec.Command("ffmpeg", "-i", fp)
      stdoutStderr, err := cmd.CombinedOutput()
      if err != nil {
        fmt.Println(err)
      }
      fmt.Printf("%s\n", stdoutStderr)
    }
    if msg == "__EOF" {
      fmt.Println("Read ", msg)
      ch <- "__DONE"
    }
  }
}

