package components

import (
	"github.com/golang/glog"
	"os"
	"os/exec"
	"path"
)

type Media struct {
	mediaPath string
}

func (m Media) Process(filepath string, postProcess func()) {
	glog.V(2).Info("File path ", filepath, " Media file is being processed.")
	cmd := exec.Command("ffmpeg", "-y", "-i", filepath, path.Join("outbox", path.Base(filepath)))
	_, err := cmd.CombinedOutput()
	if err != nil {
		glog.V(2).Info("Processing failed for ", filepath, "Moving file to error folder.")
		glog.V(2).Info(err)
		err = os.Rename(filepath, path.Join("outbox", "errors", path.Base(filepath)))
		if err != nil {
			glog.V(2).Info("Error Movement failed ", err)
		}
	} else {
		glog.V(2).Info("Successfully complete processing for ", filepath)
	}
        postProcess()
}
