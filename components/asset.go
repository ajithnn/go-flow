package components

import "sync"

type Asset interface {
	Process(string, chan struct{}, *sync.WaitGroup)
}
