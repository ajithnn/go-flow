package components

import (
	"github.com/golang/glog"
)

type NotImplemented struct {
  path string
}

func (m NotImplemented) Process(filepath string,postProcess func()) {
  glog.V(2).Info("Type not implemented, Please configure in pipes.json and components/asset.go files for path %s",filepath)
  postProcess()
}
