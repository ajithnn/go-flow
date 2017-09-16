package main

import "fmt"

import "time"
import "path"
import "./scanner"
import "os"
import "./components"
import "strings"
import "sync"

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

	ch := make(chan string)
	w := scanner.FileScanner{scanPath, 30.00, make(chan string), whiteList}
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
	var wg sync.WaitGroup
	for {
		filepath := <-c
		if filepath == "__EOF" {
			fmt.Println("End of current Scan")
			wg.Wait()
			ch <- "__DONE"
		} else if len(filepath) > 0 {
			typeToProcess := getTypeFromFilePath(filepath)
			actualProcess(typeToProcess, filepath, token, &wg)
		}
	}
}

func getTypeFromFilePath(filepath string) components.Asset {
	dir := path.Dir(filepath)
	if strings.Contains(dir, "media") {
		return components.Media{}
	} else {
		return components.Meta{}
	}
}

func actualProcess(processType components.Asset, filepath string, token chan struct{}, wg *sync.WaitGroup) {
	token <- struct{}{}
	wg.Add(1)
	go processType.Process(filepath, token, wg)
}
