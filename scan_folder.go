package main

import (
	"flag"
	"github.com/ajithnn/go-flow/components"
	"github.com/ajithnn/go-flow/scanner"
	"github.com/golang/glog"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

func init() {
	flag.Parse()
}

func main() {
	inputArgs := flag.Args()[0:]
	if len(inputArgs) != 2 {
		glog.V(2).Infof("Usage: main <scan_path> <comma separated whitelist>")
		glog.V(2).Infof("eg : main Inbox/ 'Media,Transcode'")
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
			glog.V(2).Infof("Waiting for next scan")
			glog.Flush()
			time.Sleep(time.Second * 30)
		}
	}
}

func process(c <-chan string, ch chan<- string, token chan struct{}) {
	var wg sync.WaitGroup
	for {
		filepath := <-c
		if filepath == "__EOF" {
			glog.V(2).Infof("End of current Scan")
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
