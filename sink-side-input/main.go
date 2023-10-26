package main

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/numaproj/numaflow-go/pkg/sideinput"
	"log"
	"os"
	"path"
	"sync"

	"github.com/go-redis/redis/v8"
	sinksdk "github.com/numaproj/numaflow-go/pkg/sinker"
)

var sideInputName = "myticker"
var sideInputContent string
var sideInputMutex sync.Mutex

// This redis UDSink is created for numaflow e2e tests. This handle function assumes that
// a redis instance listening on address redis:6379 has already be up and running.
func handle(ctx context.Context, datumStreamCh <-chan sinksdk.Datum) sinksdk.Responses {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	result := sinksdk.ResponsesBuilder()
	for d := range datumStreamCh {
		_ = d.EventTime()
		_ = d.Watermark()

		// We use redis hashes to store messages.
		// The name of a hash is pipelineName:sinkName.
		// Each field of a hash is the content of a message and value of the field is the no. of occurrences of the message.
		hkey := fmt.Sprintf("%s:%s", os.Getenv("NUMAFLOW_PIPELINE_NAME"), os.Getenv("NUMAFLOW_VERTEX_NAME"))
		sideInputMutex.Lock()
		content := sideInputContent

		sideInputMutex.Unlock()
		log.Printf("Incremented by 1 the no. of occurrences of %s under hash key %s\n", content, hkey)
		log.Printf(" %s Side input Sideinputcontent----", sideInputContent)
		log.Printf(" %s Side input content----", content)

		err := client.HIncrBy(ctx, hkey, content, 1).Err()
		if err != nil {
			log.Println("Set Error - ", err)
		} else {
			log.Printf("Incremented by 1 the no. of occurrences of %s under hash key %s\n", content, hkey)
		}

		id := d.ID()
		result = result.Append(sinksdk.ResponseOK(id))
	}
	return result
}

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

	// Start a goroutine to listen for events from the watcher
	go fileWatcher(watcher, sideInputName)
	err = sinksdk.NewServer(sinksdk.SinkerFunc(handle)).Start(context.Background())
	if err != nil {
		log.Fatal(err)

	}
}

func fileWatcher(watcher *fsnotify.Watcher, sideInputName string) {
	log.Println("Watching for changes in side input file: ", sideinput.DirPath)
	p := path.Join(sideinput.DirPath, sideInputName)
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("watcher.Events channel closed")
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create && event.Name == p {
				log.Println("Side input file has been created:", event.Name)
				b, err := os.ReadFile(p)
				if err != nil {
					log.Println("Failed to read side input file: ", err)
				}
				// Store the file content in the global variable and protect with mutex
				sideInputMutex.Lock()
				sideInputContent = string(b)
				sideInputMutex.Unlock()
				// Perform some operation here, can update the value in a cache/variable
				log.Println("File contents:----------------------- ", sideInputContent)
				log.Println("File contents Original-- :----------------------- ", string(b))
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Println("watcher.Errors channel closed")
				return
			}
			log.Println("error:", err)
		}
	}
}
