package components

import (
  "github.com/golang/glog"
  "io/ioutil"
  "os"
  "path"
)

type Meta struct {
  metaPath string
}

func (m Meta) Process(filepath string, postProcess func()) {
  defer postProcess()
  glog.V(2).Info("Processing Meta file ", filepath)
  dat, err := ioutil.ReadFile(filepath)
  if err != nil {
    glog.V(2).Info("Error reading meta file ", filepath)
  } else {
    glog.V(2).Info("Length of meta file ", filepath, " is ", len(dat))
    err = os.Rename(filepath,path.Join("outbox","meta",path.Base(filepath)))
    if err != nil{
      glog.V(2).Info("Error moving meta ", filepath," error is ",err)
    }
  }
  return
}
