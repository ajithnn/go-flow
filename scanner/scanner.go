package scanner

import "fmt"
import "path/filepath"
import "os"
import "runtime"
import "time"
import "sync"

type FileScanner struct{
  Path string
  StableTimeout float64
  OutChannel chan string
  SyncLock sync.WaitGroup
}

type fileScan interface{
  Scan()
  isLock(string, os.FileInfo)
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

func (w FileScanner) isLock(path string,f os.FileInfo) {
  name := f.Name()
  if osType() == "windows" {
    err := os.Rename(path,path)
    if err != nil {
      fmt.Println("File still locked", err)
      w.OutChannel <- ""
      w.SyncLock.Done()
    } else {
      fmt.Println("sent", name)
      w.OutChannel <- name
      w.SyncLock.Done()
    }
  } else {
    mod := f.ModTime()
    if time.Now().Sub(mod).Seconds() > w.StableTimeout {
      fmt.Println("sent", name)
      w.OutChannel <- name
      w.SyncLock.Done()
    } else {
      fmt.Println("Locked", name)
      w.OutChannel <- ""
      w.SyncLock.Done()
    }
  }
}

func osType() string {
  return runtime.GOOS
}

func process_files(w FileScanner) filepath.WalkFunc {
  return func(path string, info os.FileInfo, err error) error {
    if !info.IsDir() {
        w.SyncLock.Add(1)
        go w.isLock(path, info)
    }
    return nil
  }
}
