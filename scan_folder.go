package main

import (
  "flag"
  "github.com/ajithnn/go-flow/components"
  "github.com/ajithnn/go-flow/scanner"
  "github.com/golang/glog"
  "os"
  "encoding/json"
  "io/ioutil"
  "path"
  "strings"
  "sync"
  "time"
)


var pipeChannels = make(map[string](chan struct{}))
var channelTypes = make(map[string]string)
var channelState = make(map[string]int)

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
  readConfigAndCreateChannels()
  scanPath := inputArgs[0]
  whiteList := strings.Split(inputArgs[1], ",")

  ch := make(chan string)
  w := scanner.FileScanner{scanPath, 30.00, make(chan string), whiteList}
  for {
    go process(w.OutChannel, ch)
    go w.Scan()
    end := <-ch
    if end == "__DONE" {
      glog.V(2).Infof("Waiting for next scan")
      glog.Flush()
      time.Sleep(time.Second * 30)
    }
  }
}

func process(c <-chan string, ch chan<- string) {
  var wg sync.WaitGroup
  for {
    filepath := <-c
    if filepath == "__EOF" {
      glog.V(2).Infof("End of current Scan")
      wg.Wait()
      ch <- "__DONE"
    } else if len(filepath) > 0 {
      typeToProcess,typeName := getTypeFromFilePath(filepath)
      actualProcess(typeToProcess, typeName, filepath, &wg)
    }
  }
}

func readConfigAndCreateChannels(){
  configFilePath := "./pipes.json"
  var tempPipe interface{}
  configFile, _ := ioutil.ReadFile(configFilePath)
  json.Unmarshal(configFile, &tempPipe)
  curPipe := tempPipe.(map[string]interface{})
  for k,v := range curPipe{
    tempType := v.(map[string]interface{})
    glog.V(2).Infof("String %s Media values %f",k,tempType["capacity"].(float64))
    pipeChannels[k] = make(chan struct{},int(tempType["capacity"].(float64)))
    channelTypes[k] = tempType["type"].(string)
    channelState[k] = 0
  }
}

func getTypeFromFilePath(filepath string) (components.Asset,string) {
  dir := path.Dir(filepath)
  for k := range pipeChannels{
    if strings.Contains(strings.ToLower(dir),strings.ToLower(k)){
      return components.TypeMap[k],k
    }
  }
  return components.NotImplemented{},""
}

func allowCommonchannels() bool {
  ret := true
  for pipeName,_ := range pipeChannels{
    if channelTypes[pipeName] == "separate"{
      if channelState[pipeName] > 0{
        ret = false
        break
      }
    }
  }
  return ret
}

func actualProcess(processType components.Asset,typeName string, filepath string, wg *sync.WaitGroup) {
  for pipeName,pipeChannel := range pipeChannels{
    if pipeName == typeName && channelTypes[pipeName] == "separate" {
      select {
      case pipeChannel <- struct{}{}:
        wg.Add(1)
        go processType.Process(filepath,func(){
          <-pipeChannel
          channelState[typeName]--
          wg.Done()
        })
        channelState[typeName]++
      default:
        glog.V(2).Infof("All channels blocked for type %s",typeName)
      }
    }
  }
  if allowCommonchannels(){
    BreakToHere:
    for pipeName,pipeChannel := range pipeChannels{
      if pipeName == typeName && channelTypes[pipeName] == "common" {
        select {
        case pipeChannel <- struct{}{}:
          wg.Add(1)
          go processType.Process(filepath,func(){
            <-pipeChannel
            channelState[typeName]--
            wg.Done()
          })
          channelState[typeName]++
          break BreakToHere
        default:
          glog.V(2).Infof("All channels blocked for type %s",typeName)
        }
      }
    }
  }
}

