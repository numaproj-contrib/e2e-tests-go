package main

import (
	"context"
	sideinputsdk "github.com/numaproj/numaflow-go/pkg/sideinput"
	"log"
)

var counter = 0

// handle is the side input handler function.
func handle(_ context.Context) sideinputsdk.Message {
	// randomly drop side input message. Note that the side input message is not retried.
	// NoBroadcastMessage() is used to drop the message and not to
	// broadcast it to other side input vertices.
	counter = (counter + 1) % 10
	if counter%2 == 0 {
		return sideinputsdk.BroadcastMessage([]byte(`test-data`))
	}
	// BroadcastMessage() is used to broadcast the message with the given value to other side input vertices.
	// val must be converted to []byte.
	return sideinputsdk.BroadcastMessage([]byte(`test-data`))
}

func main() {
	// Start the side input server.
	err := sideinputsdk.NewSideInputServer(sideinputsdk.RetrieveFunc(handle)).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start side input server: ", err)
	}
}
