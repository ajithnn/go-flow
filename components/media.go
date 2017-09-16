package components

import "fmt"
import "path"
import "os"
import "os/exec"
import "sync"

type Media struct {
	mediaPath string
}

func (m Media) Process(filepath string, token chan struct{}, wg *sync.WaitGroup) {
	fmt.Println("File path ", filepath, " Media file is being processed.")
	cmd := exec.Command("ffmpeg", "-y", "-i", filepath, path.Join("outbox", path.Base(filepath)))
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Processing failed for ", filepath, "Moving file to error folder.")
		fmt.Println(err)
		err = os.Rename(filepath, path.Join("outbox", "errors", path.Base(filepath)))
		if err != nil {
			fmt.Println("Error Movement failed ", err)
		}
	} else {
		fmt.Println("Successfully complete processing for ", filepath)
	}
	wg.Done()
	<-token
}
