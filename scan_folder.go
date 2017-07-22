package main

import "fmt"
import "time"
import "sync"
import "path"
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

	token := make(chan struct{}, 2)
	scanPath := inputArgs[0]
	whiteList := strings.Split(inputArgs[1], ",")

	var wg sync.WaitGroup
	ch := make(chan string)
	w := scanner.FileScanner{scanPath, 30.00, make(chan string), wg, whiteList}
	for {
		go process(w.OutChannel, ch, token)
		go w.Scan()
		end := <-ch
		if end == "__DONE" {
			fmt.Println("Waiting for next scan")
			time.Sleep(time.Second * 30)
		}
	}
}

func process(c <-chan string, ch chan<- string, token chan struct{}) {
	for {
		msg := <-c
		fmt.Println("Read", msg)
		if msg == "__EOF" {
			fmt.Println("Read ", msg)
			ch <- "__DONE"
		} else if len(msg) > 0 {
			fmt.Println("Processing " + msg)
			go actualProcess(msg, token)
		}
	}
}

func actualProcess(fullFilePath string, token chan struct{}) {
	token <- struct{}{}
	fmt.Println("token acquired for " + fullFilePath)
	cmd := exec.Command("ffmpeg", "-y", "-i", fullFilePath, path.Join("outbox", path.Base(fullFilePath)))
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("token reverted for " + fullFilePath)
	<-token
	return
}
