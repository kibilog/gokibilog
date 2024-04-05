package gokibilog

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var apiToken = flag.String("apiToken", "", "api token")
var logId = flag.String("logId", "", "log id")

func Test_client_Send(t *testing.T) {
	flag.Parse()
	if *apiToken == "" || *logId == "" {
		log.Printf("Skipped client test. Need set -apiToken and -logId")
		return
	}

	t.Run("levels", func(t *testing.T) {
		levels := []MessageLevel{LevelDebug, LevelInfo, LevelNotice, LevelWarning, LevelError, LevelCritical, LevelAlert, LevelEmergency}
		k := GetInstance()
		k.SetAuthToken(*apiToken)

		for _, level := range levels {

			l, _ := NewLogPool(*logId)
			k.AddLogPool(l)

			m, _ := NewMessage("level test", level)
			l.AddMessage(m)

			err := k.SendMessages()
			if len(err) != 0 {
				t.Errorf("SendMessages(), levels, Error: %+v", err)
			}
		}

	})

	t.Run("partition", func(t *testing.T) {
		k := GetInstance()
		k.SetAuthToken(*apiToken)
		l, _ := NewLogPool(*logId)
		k.AddLogPool(l)

		genPart := func() string {
			prefix := strconv.Itoa(rand.Intn(8888) + 1000)
			postfix := strconv.Itoa(rand.Intn(8888) + 1000)
			return fmt.Sprintf("%sf48e-4047-461f-bf97-06a5cd5c%s", prefix, postfix)
		}

		partitions := []string{genPart(), genPart()}
		for _, part := range partitions {

			for _, v := range []string{"1", "2"} {
				m, _ := NewMessage(fmt.Sprintf("partition %v message %v", part, v), LevelDebug)
				_ = m.SetPartition(part)
				l.AddMessage(m)
			}

		}

		err := k.SendMessages()
		if len(err) != 0 {
			t.Errorf("SendMessages(), partition,  Error: %+v", err)
		}
	})

	t.Run("param map", func(t *testing.T) {
		k := GetInstance()
		k.SetAuthToken(*apiToken)
		l, _ := NewLogPool(*logId)
		k.AddLogPool(l)

		m, _ := NewMessage("params map", LevelDebug)
		var param any
		json.Unmarshal([]byte(`{
            "orderId": 123456,
            "status": {
                "past": "sent",
                "current": "delivered"
            }
        }`), &param)
		m.SetParams(param)

		l.AddMessage(m)

		err := k.SendMessages()
		if len(err) != 0 {
			t.Errorf("SendMessages(), param map,  Error: %+v", err)
		}
	})

	t.Run("param struct", func(t *testing.T) {
		k := GetInstance()
		k.SetAuthToken(*apiToken)
		l, _ := NewLogPool(*logId)
		k.AddLogPool(l)

		m, _ := NewMessage("params struct", LevelDebug)
		param := struct {
			OrderId int
			Status  struct {
				Past    string
				Current string
			}
		}{
			OrderId: 123456,
			Status: struct {
				Past    string
				Current string
			}{Past: "sent", Current: "delivered"},
		}
		m.SetParams(param)

		l.AddMessage(m)

		err := k.SendMessages()
		if len(err) != 0 {
			t.Errorf("SendMessages(), param struct,  Error: %+v", err)
		}
	})

	t.Run("param array", func(t *testing.T) {
		k := GetInstance()
		k.SetAuthToken(*apiToken)
		l, _ := NewLogPool(*logId)
		k.AddLogPool(l)

		m, _ := NewMessage("params array", LevelDebug)
		param := []string{"sent", "delivered"}
		m.SetParams(param)

		l.AddMessage(m)

		err := k.SendMessages()
		if len(err) != 0 {
			t.Errorf("SendMessages(), param array,  Error: %+v", err)
		}
	})

	t.Run("createdAt", func(t *testing.T) {
		k := GetInstance()
		k.SetAuthToken(*apiToken)
		l, _ := NewLogPool(*logId)
		k.AddLogPool(l)

		genPart := func() string {
			prefix := strconv.Itoa(rand.Intn(8888) + 1000)
			postfix := strconv.Itoa(rand.Intn(8888) + 1000)
			return fmt.Sprintf("%sf48e-4047-461f-bf97-06a5cd5c%s", prefix, postfix)
		}
		part := genPart()

		m, _ := NewMessage("createdAt now", LevelDebug)
		m.SetPartition(part)
		l.AddMessage(m)

		m, _ = NewMessage("createdAt + 15 sec", LevelDebug)
		m.SetPartition(part)
		m.SetCreatedAt(time.Now().Add(15 * time.Second))
		l.AddMessage(m)

		err := k.SendMessages()
		if len(err) != 0 {
			t.Errorf("SendMessages(), param array,  Error: %+v", err)
		}
	})
}
