package scanner

import "fmt"
import "path"
import "path/filepath"
import "strings"
import "os"
import "runtime"
import "time"
import "sync"

type FileScanner struct{
  Path string
  StableTimeout float64
  OutChannel chan string
  SyncLock sync.WaitGroup
  Whitelist []string
}

type fileScan interface{
  Scan()
  isLock(string, os.FileInfo)
  isWhiteListed(string, string)
}

func (w FileScanner) Scan() {
  err := filepath.Walk(w.Path, process_files(w))
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println("Waiting")
  w.SyncLock.Wait()
  w.OutChannel <- "__EOF"
}

func (w FileScanner) isLock(pth string,f os.FileInfo) {
  if osType() == "windows" {
    err := os.Rename(pth,pth)
    if err != nil {
      fmt.Println("File still locked", err)
      w.OutChannel <- ""
      w.SyncLock.Done()
    } else {
      fmt.Println("sent", pth)
      w.OutChannel <- pth
      w.SyncLock.Done()
    }
  } else {
    mod := f.ModTime()
    if time.Now().Sub(mod).Seconds() > w.StableTimeout {
      fmt.Println("sent", pth)
      w.OutChannel <- pth
      w.SyncLock.Done()
    } else {
      fmt.Println("Locked", pth)
      w.OutChannel <- ""
      w.SyncLock.Done()
    }
  }
}

func osType() string {
  return runtime.GOOS
}

func (w FileScanner) isWhiteListed(basePath string,curFilePath string) bool {
  for _,folder := range w.Whitelist {
    if strings.Contains(curFilePath,path.Join(basePath,folder)) {
      return true
    }
  }
  return false
}

func process_files(w FileScanner) filepath.WalkFunc {
  return func(pth string, info os.FileInfo, err error) error {
    if !info.IsDir() {
      if w.isWhiteListed(w.Path,pth) {
        w.SyncLock.Add(1)
        go w.isLock(pth, info)
      } else{
        fmt.Println("Path not in whitelist ", pth)
      }
    }
    return nil
  }
}
