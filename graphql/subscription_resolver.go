package graphql

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

type client struct {
	Message chan *Message
}

var clients = map[string]*client{}

var mutex = &sync.Mutex{}

func Broadcaster() {
	for {
		select {
		case <-time.After(1 * time.Second):
			for _, v := range clients {
				v.Message <- &Message{Text: "hello" + time.Now().String()}
			}
		}
	}
}

func (r *subscriptionResolver) Messages(ctx context.Context) (<-chan *Message, error) {
	id := randString(8)

	events := make(chan *Message, 1)

	client := &client{Message: events}
	clients[id] = client

	go func() {
		<-ctx.Done()
		mutex.Lock()
		delete(clients, id)
		mutex.Unlock()
	}()

	return events, nil
}

func randString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
