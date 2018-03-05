package flow

import (
	"encoding/json"
	scanner "github.com/ajithnn/sense"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type Stage interface {
	Process(string, interface{}, func())
}

type notImplemented struct {
}

func (n notImplemented) Process(filepath string, config interface{}, postProcess func()) {
	defer postProcess()
	glog.V(2).Info("Type not found , unable to process", filepath)
	return
}

type Flow struct {
	ScanPath           string
	PipePath           string
	WhiteList          string
	Timeout            float64
	ScanTimeout        time.Duration
	TypeMap            map[string]Stage
	GetPrioritizedList func([]string) []string
}

var pipeChannels = make(map[string](chan struct{}))
var channelTypes = make(map[string]string)
var processList = make(map[string]bool)
var typeConfig = make(map[string]interface{})
var filePathList []string
var FlowConfig Flow

func Trigger(config Flow) {
	FlowConfig = config
	readConfigAndCreateChannels()

	ch := make(chan string)
	w := scanner.FileScanner{FlowConfig.ScanPath, FlowConfig.Timeout, make(chan string), strings.Split(FlowConfig.WhiteList, ","), &processList}
	for {
		go process(w.OutChannel, ch)
		go w.Scan()
		end := <-ch
		if end == "__DONE" {
			glog.V(2).Infof("Waiting for next scan")
			glog.Flush()
			time.Sleep(FlowConfig.ScanTimeout)
		}
	}
}

func process(c <-chan string, ch chan<- string) {
	for {
		filepath := <-c
		if filepath == "__EOF" {
			glog.V(2).Infof("End of current Scan will continue after wait.")
			pushByPriority(filePathList)
			ch <- "__DONE"
		} else if len(filepath) > 0 {
			filePathList = append(filePathList, filepath)
		}
	}
}

func pushByPriority(fileList []string) {
	filePathList = FlowConfig.GetPrioritizedList(fileList)
	filepath := ""
	for len(filePathList) > 0 {
		filepath, filePathList = filePathList[0], filePathList[1:]
		typeToProcess, typeName := getTypeFromFilePath(filepath)
		actualProcess(typeToProcess, typeName, filepath)
	}
}

func readConfigAndCreateChannels() {
	configFilePath := FlowConfig.PipePath
	var tempPipe interface{}
	commonChannel := make(chan struct{}, 1)
	configFile, _ := ioutil.ReadFile(configFilePath)
	json.Unmarshal(configFile, &tempPipe)
	curPipe := tempPipe.(map[string]interface{})
	for k, v := range curPipe {
		tempType := v.(map[string]interface{})
		glog.V(2).Infof("Type and Concurrency : ", k, tempType["capacity"].(float64))
		if tempType["type"].(string) == "separate" {
			pipeChannels[k] = make(chan struct{}, int(tempType["capacity"].(float64)))
		} else {
			pipeChannels[k] = commonChannel
		}
		channelTypes[k] = tempType["type"].(string)
		typeConfig[k] = tempType["config"]
	}
}

func getTypeFromFilePath(filepath string) (Stage, string) {
	dir := path.Dir(filepath)
	for k := range pipeChannels {
		if strings.Contains(strings.ToLower(dir), strings.ToLower(k)) {
			return FlowConfig.TypeMap[k], k
		}
	}
	return notImplemented{}, ""
}

func actualProcess(processType Stage, typeName string, filepath string) {
	select {
	case pipeChannels[typeName] <- struct{}{}:
		if _, err := os.Stat(filepath); !os.IsNotExist(err) {
			processList[filepath] = true
			go processType.Process(filepath, typeConfig[typeName], func() {
				delete(processList, filepath)
				<-pipeChannels[typeName]
				glog.V(2).Infof("Released channel and cleared file hold.")
			})
		}
	default:
		glog.V(2).Infof("All channels blocked for type %s", typeName)
	}
}
