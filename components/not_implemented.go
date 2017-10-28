package components

import (
	"github.com/golang/glog"
	"sync"
)

type NotImplemented struct {
  path string
}

func (m NotImplemented) Process(filepath string, token chan struct{}, wg *sync.WaitGroup) {
  glog.V(2).Info("Type not implemented, Please configure in pipes.json and components/asset.go files for path %s",filepath)
  wg.Done()
  <-token
}
