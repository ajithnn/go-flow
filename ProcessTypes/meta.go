package ProcessTypes

import "fmt"

type Meta struct {
	metaPath string
}

func (m Meta) Process(filepath string, token chan struct{}) {
	fmt.Println("File path is Meta ", filepath)
	<-token
}
