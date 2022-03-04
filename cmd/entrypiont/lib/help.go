package lib

import (
	"errors"
	"k8s.io/klog/v2"
	"os"
	"time"
)

func CheckFlags() {
	if waitFile == "" || out == "" || command == "" {
		klog.Exit("flag error")
	}
}

// CheckWaitFile check waitFile exist
func CheckWaitFile() {
	for {
		if _, err := os.Stat(waitFile); err == nil { // file exist return
			return
		} else if errors.Is(err, os.ErrNotExist) { // file is not exist continue
			time.Sleep(time.Millisecond * 20)
			continue
		} else {
			klog.Exit(err)
		}

	}
}
