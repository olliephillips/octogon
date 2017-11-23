package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

// send file using SSH secure copy
func sendSCP(tfile string) {

	logger("starting SSH session", false)

	// config
	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{
			ssh.Password(*pwd),
		},
		Timeout: 0 * time.Second,
	}

	client := scp.NewClient(*host, config)

	err := client.Connect()
	if err != nil {
		logger("couldn't establisch a connection to the remote server ", true)
	}

	// open file
	f, err := os.Open(tfile)
	if err != nil {
		logger(err, false)
	}

	// clean up
	defer client.Session.Close()
	defer f.Close()

	remotePath := ".octoprint/watched/"

	// handle remote path flag if given
	if *remotefolder != "" {
		remotePath = ".octoprint/uploads/"
		remotePath += *remotefolder + "/"
	}

	remotePath += filepath.Base(tfile)
	log.Println(remotePath)
	logger("copying "+filepath.Base(tfile), false)
	client.CopyFromFile(*f, remotePath, "0655")
	logger("SSH session ended", false)

	// stop tracking file
	untrack(tfile)

	// should file be deleted from local file system
	checkAndDelete(tfile)
}
