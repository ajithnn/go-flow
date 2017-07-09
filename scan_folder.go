package main

import "fmt"
import "time"
import "sync"
import "./scanner"
import "os/exec"
import "os"
import "strings"

func main() {

  inputArgs := os.Args[1:]
  if len(inputArgs) != 2 {
    fmt.Println("Usage: main <scan_path> <comma separated whitelist>")
    fmt.Println("eg : main Inbox/ 'Media,Transcode'")
    os.Exit(1)
  }

  scanPath := inputArgs[0]
  whiteList := strings.Split(inputArgs[1],",")

  var wg sync.WaitGroup
  ch := make(chan string)
    w := scanner.FileScanner{ scanPath, 300.00, make(chan string), wg, whiteList }
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
    if msg == "__EOF" {
      fmt.Println("Read ", msg)
      ch <- "__DONE"
    } else if len(msg) > 0 {
      fmt.Println("Processing " + msg)
      actualProcess(msg)
    }
  }
}

func actualProcess(fullFilePath string) {
  cmd := exec.Command("ffmpeg", "-i", fullFilePath)
  stdoutStderr, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Println(err)
  }
  fmt.Printf("%s\n", stdoutStderr)
  return
}

