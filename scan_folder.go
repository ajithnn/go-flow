package main

import "fmt"

import "time"
import "path"
import "./scanner"
import "os"
import "./ProcessTypes"
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
	for {
		filepath := <-c
		if filepath == "__EOF" {
			fmt.Println("End of current Scan")
			ch <- "__DONE"
		} else if len(filepath) > 0 {
			fmt.Println("Scanned ", filepath)
			typeToProcess := getTypeFromFilePath(filepath)
			actualProcess(typeToProcess, filepath, token)
		}
	}
}

func getTypeFromFilePath(filepath string) interface{} {
	dir := path.Dir(filepath)
	if strings.Contains(dir, "media") {
		return ProcessTypes.Media{}
	} else {
		return ProcessTypes.Meta{}
	}
}

func actualProcess(processType interface{}, filepath string, token chan struct{}) {
	token <- struct{}{}
	switch val := processType.(type) {
	case ProcessTypes.Media:
		go val.Process(filepath, token)
	case ProcessTypes.Meta:
		go val.Process(filepath, token)
	}
}
