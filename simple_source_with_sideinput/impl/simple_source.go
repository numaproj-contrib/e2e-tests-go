package impl

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/numaproj/numaflow-go/pkg/sideinput"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	sourcesdk "github.com/numaproj/numaflow-go/pkg/sourcer"
)

var sideInputMutex sync.Mutex

// A global channel that will recev the data as soon as the file watcher sends it
var globalChan = make(chan string)
var contentValue string

// SimpleSource is a simple source implementation.
type SimpleSource struct {
	readIdx  int64
	toAckSet map[int64]struct{}
	lock     *sync.Mutex
}

func NewSimpleSource() *SimpleSource {
	return &SimpleSource{
		readIdx:  0,
		toAckSet: make(map[int64]struct{}),
		lock:     new(sync.Mutex),
	}
}

func (s *SimpleSource) Pending(_ context.Context) int64 {
	// The simple source always returns zero to indicate there is no pending record.
	return 0
}

func (s *SimpleSource) Read(_ context.Context, readRequest sourcesdk.ReadRequest, messageCh chan<- sourcesdk.Message) {
	// Handle the timeout specification in the read request.
	ctx, cancel := context.WithTimeout(context.Background(), readRequest.TimeOut())
	defer cancel()

	// Read the data from the source and send the data to the message channel.
	for i := 0; uint64(i) < readRequest.Count(); i++ {
		select {
		case <-ctx.Done():
			// If the context is done, the read request is timed out.
			return
		case value := <-globalChan:
			log.Println("Value received in global chan444444--------", value)
			s.lock.Lock()
			// Otherwise, we read the data from the source and send the data to the message channel.
			offsetValue := serializeOffset(s.readIdx)
			// save the value in the global variable and return that variable
			messageCh <- sourcesdk.NewMessage(
				[]byte(value),
				sourcesdk.NewOffset(offsetValue, "0"),
				time.Now())
			// Mark the offset as to be acked, and increment the read index.
			s.toAckSet[s.readIdx] = struct{}{}
			s.readIdx++
			s.lock.Unlock()

		}
	}

}

func (s *SimpleSource) Ack(_ context.Context, request sourcesdk.AckRequest) {
	for _, offset := range request.Offsets() {
		delete(s.toAckSet, deserializeOffset(offset.Value()))
	}
}

func serializeOffset(idx int64) []byte {
	return []byte(strconv.FormatInt(idx, 10))
}

func deserializeOffset(offset []byte) int64 {
	idx, _ := strconv.ParseInt(string(offset), 10, 64)
	return idx
}

func FileWatcher(watcher *fsnotify.Watcher, sideInputName string) {
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
				contentValue = string(b)
				globalChan <- contentValue
				sideInputMutex.Unlock()
				// Perform some operation here, can update the value in a cache/variable
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
