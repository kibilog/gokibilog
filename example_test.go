package gokibilog_test

import (
	"encoding/json"
	"fmt"
	"github.com/kibilog/gokibilog"
	"io"
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

func Example_partition() {
	kibilog := gokibilog.GetInstance()
	kibilog.SetAuthToken("01htapnjvw83bz7xjhgcdwtry4")

	logPool, err := gokibilog.NewLogPool("01htapms8kf6wyngde3mvyjn8x")
	if err != nil {
		log.Fatalf("Error creating logpool: %v", err)
		return
	}
	kibilog.AddLogPool(logPool)

	// We need a partition in order to collect messages within a logical group.
	// For example, if you need to create a message history with a change in a specific order.
	// The partition must be a UUID. We recommend using v4 for random events.
	// But if you need to get a stable UUID based on the identifier (in the example with
	// the order, the identifier may be the order id), then v3 or v5 will do.
	// We do not impose restrictions on the UUID version, but the partition must be a UUID.

	// Let's get json with comments on posts from a public service and group the messages into partitions.
	resp, err := http.Get("https://jsonplaceholder.typicode.com/comments")
	if err != nil {
		log.Fatalf("Error getting status: %v", err)
		return
	}
	defer resp.Body.Close()

	type Comment struct {
		PostId int    `json:"postId"`
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Body   string `json:"body"`
	}
	var aComments []Comment

	body, err := io.ReadAll(resp.Body)
	json.Unmarshal(body, &aComments)

	for _, comment := range aComments {
		// Since there are many packages that can make a UUID, we will not focus on how to get a UUID.
		// Let's imagine that we have a GetUUIDByPostId method that returns a string with a UUID.
		uuid := GetUUIDByPostId(comment.PostId)
		m, err := gokibilog.NewMessage(
			fmt.Sprintf("New comment from %s to %d post", comment.Email, comment.PostId),
			gokibilog.LevelInfo)
		if err != nil {
			log.Fatalf("Error getting status: %v", err)
		}
		m.SetPartition(uuid)

		logPool.AddMessage(m)
	}

	// After sending the messages, we can see in Kibilog.com how did they group themselves by partitions.
	kibilog.SendMessages()
}

func Example_params() {
	// You can add parameters to messages.
	// For example, we write all incoming and outgoing requests.
	// In this case, it would be useful for us to record the route
	// that was accessed, the method, the incoming request and the response.
	kibilog := gokibilog.GetInstance()
	kibilog.SetAuthToken("01htapnjvw83bz7xjhgcdwtry4")

	logPool, err := gokibilog.NewLogPool("01htapms8kf6wyngde3mvyjn8x")
	if err != nil {
		log.Fatalf("Error creating logpool: %v", err)
		return
	}
	kibilog.AddLogPool(logPool)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/user/create", func(w http.ResponseWriter, r *http.Request) {
		// Let's record the incoming request
		requestUUID := SomeFuncGetUUIDV4()
		m, _ := gokibilog.NewMessage("Incoming request", gokibilog.LevelInfo)
		m.SetParams(map[string]any{
			"uri":    r.URL.Path,
			"method": r.Method,
		})
		m.SetPartition(requestUUID)
		logPool.AddMessage(m)

		// Let's start processing the request
		type Person struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		var p Person
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			// Oops... Something went wrong.
			m, _ = gokibilog.NewMessage("Response", gokibilog.LevelError)
			m.SetParams(map[string]any{
				"error":              err.Error(),
				"httpStatusResponse": http.StatusBadRequest,
			})
			m.SetPartition(requestUUID)
			logPool.AddMessage(m)
			kibilog.SendMessages()

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// We now have a user entity. Let's add it.
		m, _ = gokibilog.NewMessage("Request params", gokibilog.LevelInfo)
		m.SetParams(map[string]any{
			"_requestBody": p,
		})
		m.SetPartition(requestUUID)
		logPool.AddMessage(m)

		// Let's prepare a response
		type Response struct {
			Done int `json:"done"`
			Id   int `json:"id"`
		}
		resp := Response{
			Done: 1,
			Id:   SomeFuncCreateUser(p),
		}

		respB, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(respB))

		m, _ = gokibilog.NewMessage("Response", gokibilog.LevelInfo)
		m.SetParams(map[string]any{
			"_response": resp,
		})
		m.SetPartition(requestUUID)
		logPool.AddMessage(m)
		kibilog.SendMessages()
	})
}
