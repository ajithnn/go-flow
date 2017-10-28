package components

import "sync"

type Asset interface {
	Process(string, chan struct{}, *sync.WaitGroup)
}

var TypeMap =  map[string]Asset{
  "Media": Media{},
  "Meta": Meta{},
}
