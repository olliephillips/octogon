package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/ssh"
)

var (
	pwd          *string
	host         *string
	ip           *string
	user         *string
	folder       *string
	remotefolder *string
	del          *bool
	session      *ssh.Session
	dir          string
)

// defaults which can be overridden on startup
const (
	HOST         = "octopi.local:22" //-hp flag
	IP           = ""                //-ip flag
	USER         = "pi"              //-u flag
	FOLDER       = ""                //-f flag
	REMOTEFOLDER = ""                //-r flag
	DELETE       = false             //-d flag
	PWD          = ""                //-p flag
)

// keeps track of copying operations
var tracked []string
var mtx sync.RWMutex

func main() {
	// process command line args for overrides to default
	pwd = flag.String("p", PWD, "Password, required")
	host = flag.String("hp", HOST, "Remote Hostname and Port to connect on")
	ip = flag.String("ip", IP, "Remote IP address and Port to connect on")
	user = flag.String("u", USER, "User account to connect with")
	folder = flag.String("f", FOLDER, "Absolute path to local folder to monitor. Default is current folder")
	remotefolder = flag.String("r", REMOTEFOLDER, "Remote folder to send files to. Default is none")
	del = flag.Bool("d", DELETE, "Delete the file after sending")
	flag.Parse()

	// check password supplied
	if *pwd == "" {
		logger("password is required e.g 'octogon -p mypassword'", true)
	}

	// check folder is absolute path
	if *folder != "" {
		if !filepath.IsAbs(*folder) {
			logger("folder path needs to be absolute path to the folder", true)
		}
	}

	// start monitor on folder (current or specified with in folder)
	startMonitor()
}

// helper for getting current directory
func getCurrentDir() (string, error) {
	directory, err := os.Getwd()
	if err != nil {
		return directory, err
	}
	return directory, nil
}

// helper for consistent formatting of log output
func logger(msg interface{}, fatal bool) {
	pretext := "Octogon:"
	if !fatal {
		log.Println(pretext, msg)
		return
	}
	log.Fatalln(pretext, msg)
}

// helper to track file so don't try transfer if already in progress
func checkAndTrack(file string) bool {
	ok := true
	for _, v := range tracked {
		if v == file {
			ok = false
			return ok
		}
	}
	mtx.Lock()
	tracked = append(tracked, file)
	mtx.Unlock()
	return ok
}

// helper to untrack file once copy operation complete
func untrack(file string) {
	index := 0
	for i, v := range tracked {
		if v == file {
			index = i
			break
		}
	}
	mtx.Lock()
	tracked = append(tracked[:index], tracked[index+1:]...)
	mtx.Unlock()
}

// helper to check if delete flag passed and delete from local file system if required
func checkAndDelete(file string) {
	if *del {
		err := os.Remove(file)
		if err != nil {
			logger("problem deleting the file", false)
		}
		logger("local copy of "+filepath.Base(file)+" deleted", false)
	}
}
