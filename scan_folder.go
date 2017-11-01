package main

import (
  "flag"
  "github.com/ajithnn/go-flow/components"
  "github.com/ajithnn/go-flow/utils"
  "github.com/ajithnn/go-flow/scanner"
  "github.com/golang/glog"
  "os"
  "encoding/json"
  "io/ioutil"
  "path"
  "strings"
  "time"
)


var pipeChannels = make(map[string](chan struct{}))
var channelTypes = make(map[string]string)
var processList = make(map[string]bool)
var filePathList []string

func init() {
  flag.Parse()
}

func main() {
  inputArgs := flag.Args()[0:]
  if len(inputArgs) != 2 {

    glog.V(2).Infof("Usage:")
    glog.V(2).Infof("go run scan_folder.go -logtostderr=true -v=2 <Inbox Path> <Comma separated whitelist of folders>")
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
  for {
    filepath := <-c
    if filepath == "__EOF" {
      glog.V(2).Infof("End of current Scan will continue after wait.")
      pushByPriority(filePathList)
      ch <- "__DONE"
    } else if len(filepath) > 0 {
      filePathList = append(filePathList,filepath)
    }
  }
}

func pushByPriority(fileList []string){
  filePathList = utils.GetPrioritizedList(fileList)
  filepath := ""
  for len(filePathList) > 0 {
    filepath,filePathList = filePathList[0], filePathList[1:]
    typeToProcess,typeName := getTypeFromFilePath(filepath)
    actualProcess(typeToProcess, typeName, filepath)
  }
}

func readConfigAndCreateChannels(){
  configFilePath := "./pipes.json"
  var tempPipe interface{}
  commonChannel := make(chan struct{},1)
  configFile, _ := ioutil.ReadFile(configFilePath)
  json.Unmarshal(configFile, &tempPipe)
  curPipe := tempPipe.(map[string]interface{})
  for k,v := range curPipe{
    tempType := v.(map[string]interface{})
    glog.V(2).Infof("String %s Media values %f",k,tempType["capacity"].(float64))
    if tempType["type"] == "separate"{
      pipeChannels[k] = make(chan struct{},int(tempType["capacity"].(float64)))
    }else{
      pipeChannels[k] =commonChannel
    }
    channelTypes[k] = tempType["type"].(string)
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

func actualProcess(processType components.Asset,typeName string, filepath string) {
  if _,ok := processList[filepath]; !ok {
    select {
    case pipeChannels[typeName] <- struct{}{}:
      processList[filepath] = true
      go processType.Process(filepath,func(){
        delete(processList,filepath)
        <-pipeChannels[typeName]
        glog.V(2).Infof("Released channel and cleared file hold.")
      })
    default:
      glog.V(2).Infof("All channels blocked for type %s",typeName)
    }
  }
}
