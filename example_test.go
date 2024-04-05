package gokibilog_test

import (
	"github.com/kibilog/gokibilog"
	"log"
	"net/http"
)

func Example_process() {
	// Create kibilog instance and set auth token (api key)
	kibilog := gokibilog.GetInstance()
	kibilog.SetAuthToken("01htapnjvw83bz7xjhgcdwtry4")

	// Let's create a LogPool into which we will add messages
	logPool, err := gokibilog.NewLogPool("01htapms8kf6wyngde3mvyjn8x")
	if err != nil {
		log.Fatalf("Error creating logpool: %v", err)
		return
	}

	// Attach the LogPool to the instance kibilog
	kibilog.AddLogPool(logPool)

	// Let's do... Well, I don't know... Request!
	resp, err := http.Get("https://example.re/get-status-404")
	if err != nil {
		log.Fatalf("Error getting status: %v", err)
		return
	}

	// ... and write the result in a message
	var message *gokibilog.Message
	if resp.StatusCode == http.StatusOK {
		message, err = gokibilog.NewMessage("We received a response with the status 200", gokibilog.LevelInfo)
		if err != nil {
			log.Fatalf("Error creating message: %v", err)
			return
		}
	} else {
		message, err = gokibilog.NewMessage("Oh no, the status code is different from 200!", gokibilog.LevelError)
		if err != nil {
			log.Fatalf("Error creating message: %v", err)
			return
		}
	}

	// Adding a message to the LogPool
	logPool.AddMessage(message)

	// It is worth understanding that we can write a lot of messages to the LogPool.
	// We can also create several logpools at once, write them to kibilog.
	// This will allow us to record messages in several logs at once.

	// When all the necessary messages are recorded in the LogPool,
	// it will be time to send them to Kibilog.com
	// The SendMessages call sends all the accumulated messages to Kibilog.com by erasing them from storage.
	errs := kibilog.SendMessages()
	if len(errs) > 0 {
		log.Fatalf("Error sending messages: %v", errs)
	}

	// By the way, we can get the previously recorded LogPool by LogID.
	// Just in case you want to make a global variable.
	logPool, err = kibilog.GetLogPoolById("01htapms8kf6wyngde3mvyjn8x")
}
