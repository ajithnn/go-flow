package scanner

import "fmt"
import "path"
import "path/filepath"
import "strings"
import "os"
import "runtime"
import "time"

type FileScanner struct {
	Path          string
	StableTimeout float64
	OutChannel    chan string
	Whitelist     []string
}

func (w FileScanner) Scan() {
	fmt.Println(" ", w.Path)
	err := filepath.Walk(w.Path, process_files(w))
	if err != nil {
		fmt.Println(err)
	}
	w.OutChannel <- "__EOF"
}

func (w FileScanner) isLock(pth string, f os.FileInfo) bool {
	if osType() == "windows" {
		err := os.Rename(pth, pth)
		if err != nil {
			fmt.Println("File still locked", err)
			return false
		}
		return true
	} else {
		mod := f.ModTime()
		if time.Now().Sub(mod).Seconds() < w.StableTimeout {
			fmt.Println("Locked", pth)
			return false
		}
		return true
	}
}

func osType() string {
	return runtime.GOOS
}

func (w FileScanner) isWhiteListed(basePath string, curFilePath string) bool {
	for _, folder := range w.Whitelist {
		if strings.Contains(curFilePath, path.Join(basePath, folder)) {
			return true
		}
	}
	return false
}

func process_files(w FileScanner) filepath.WalkFunc {
	return func(pth string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if w.isWhiteListed(w.Path, pth) && w.isLock(pth, info) {
				w.OutChannel <- pth
			} else {
				fmt.Println("Path ", pth, " not in whitelist")
			}
		}
		return nil
	}
}
