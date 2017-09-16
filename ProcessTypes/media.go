package ProcessTypes

import "fmt"

type Media struct {
	mediaPath string
}

func (m Media) Process(filepath string, token chan struct{}) {
	fmt.Println("File path is Media ", filepath)
	<-token
}
