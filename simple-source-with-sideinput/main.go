package main

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/numaproj/numaflow-go/pkg/sideinput"
	"log"

	"github.com/numaproj/numaflow-go/pkg/sourcer"

	"github.com/numaproj/numaflow-go/pkg/sourcer/examples/simple_source/impl"
)

var sideInputName = "myticker"

func main() {
	// Create a new fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	// Add a path to the watcher
	err = watcher.Add(sideinput.DirPath)
	if err != nil {
		log.Fatal(err)
	}
	go impl.FileWatcher(watcher, sideInputName)

	simpleSource := impl.NewSimpleSource()
	err = sourcer.NewServer(simpleSource).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start source server : ", err)
	}
}
