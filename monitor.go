package main

import (
	"path/filepath"

	fsnotify "gopkg.in/fsnotify.v1"
)

// set up directory monitoring
func startMonitor() {
	monitor, err := fsnotify.NewWatcher()
	if err != nil {
		logger(err, true)
	}
	defer monitor.Close()

	// which path
	if *folder != "" {
		dir = *folder
	} else {
		dir, err = getCurrentDir()
		if err != nil {
			logger(err, true)
		}
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-monitor.Events:
				if filepath.Ext(event.Name) == ".gcode" || filepath.Ext(event.Name) == ".stl" {
					if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
						msg := "file modified " + filepath.Base(event.Name)
						logger(msg, false)

						// check if not in progress and send
						if checkAndTrack(event.Name) {
							go sendSCP(event.Name)
						} else {
							logger("file transfer already in progress", false)
						}
					}
				}
			case err := <-monitor.Errors:
				logger(err, false)
			}
		}
	}()

	// add the path to our monitor
	if err = monitor.Add(dir); err != nil {
		logger(err, true)
	}
	logger("monitoring changes", false)
	<-done
}
