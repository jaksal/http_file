package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	conf *Config

	statePath string
	syncPath  string
)

func main() {
	// init log
	ex, _ := os.Executable()
	fpLog, err := os.OpenFile(filepath.Dir(ex)+"/watcher.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer fpLog.Close()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(io.MultiWriter(fpLog, os.Stdout))

	if err := initDB(filepath.Dir(ex) + "/data.db"); err != nil {
		log.Fatal(err)
	}

	conf, err = loadConfig(filepath.Dir(ex) + "/conf.json")
	if err != nil {
		log.Fatal(err)
	}

	statePath = filepath.Dir(ex) + "/state.txt"
	syncPath = filepath.Dir(ex) + "/sync.txt"

	log.Println("start..", conf)

	var wg sync.WaitGroup

	// start async uploader
	wg.Add(1)
	startUploader(1000, &wg)

	if conf.Check != "" {
		// check upload path
		checkUploadFiles(conf.Path + "/" + conf.Check)

	} else {
		// start init folder exist file
		checkExistFile(conf.Path)

		// start watch folder()
		watchFolder(conf.Path)

	}

	closeUploader()
	wg.Wait()
	log.Printf("Stopped.\n")
}
