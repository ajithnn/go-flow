package components

import "fmt"
import "io/ioutil"
import "sync"

type Meta struct {
	metaPath string
}

func (m Meta) Process(filepath string, token chan struct{}, wg *sync.WaitGroup) {
	fmt.Println("Processing Meta file ", filepath)
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading meta file ", filepath)
	} else {
		fmt.Println("Length of meta file ", filepath, " is ", len(dat))
	}
	wg.Done()
	<-token
}
