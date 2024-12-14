package server

import (
	"encoding/json"
	"simplified-twitter/feed"
	"simplified-twitter/queue"
	"sync"
)

type Response struct {
	ID      int  `json:"id"`
	Success bool `json:"success"`
}

type FeedResponse struct {
	ID   int             `json:"id"`
	Feed []feed.FeedPost `json:"feed"`
}

type SharedContext struct {
	mutex  *sync.Mutex
	cond   *sync.Cond
	wg     *sync.WaitGroup
	finish bool
}

type Config struct {
	Encoder        *json.Encoder
	Decoder        *json.Decoder
	Mode           string
	ConsumersCount int
}

func consumer(config *Config, taskQueue *queue.LockFreeQueue, feed *feed.Feed, context *SharedContext) {
	defer context.wg.Done()
	for {
		context.mutex.Lock()
		for taskQueue.IsEmpty() && !context.finish {
			context.cond.Wait()
		}
		if context.finish && taskQueue.IsEmpty() {
			context.mutex.Unlock()
			break
		}
		request := taskQueue.Dequeue()
		context.mutex.Unlock()

		response := processRequest(request, feed)
		err := config.Encoder.Encode(response)
		if err != nil {
			break
		}
	}
}

func producer(config *Config, taskQueue *queue.LockFreeQueue, context *SharedContext) {
	defer context.wg.Done()
	for {
		var request queue.Request
		err := config.Decoder.Decode(&request)
		if err != nil {
			break
		}
		if request.Command == "DONE" {
			context.mutex.Lock()
			context.finish = true
			context.cond.Broadcast()
			context.mutex.Unlock()
			break
		}
		context.mutex.Lock()
		taskQueue.Enqueue(&request)
		context.cond.Signal()
		context.mutex.Unlock()
	}
}

func processRequest(request *queue.Request, feed *feed.Feed) interface{} {
	switch request.Command {
	case "ADD":
		(*feed).Add(request.Body, request.Timestamp)
		return Response{ID: request.ID, Success: true}
	case "REMOVE":
		success := (*feed).Remove(request.Timestamp)
		return Response{ID: request.ID, Success: success}
	case "CONTAINS":
		success := (*feed).Contains(request.Timestamp)
		return Response{ID: request.ID, Success: success}
	case "FEED":
		posts := (*feed).GetAllFeeds()
		return FeedResponse{ID: request.ID, Feed: posts}
	}
	return Response{ID: request.ID, Success: false}
}

func sequentialRun(config *Config) {
	feed := feed.NewFeed()
	for {
		var request queue.Request
		err := config.Decoder.Decode(&request)
		if err != nil {
			break
		}
		if request.Command == "DONE" {
			break
		}
		response := processRequest(&request, &feed)
		err = config.Encoder.Encode(response)
		if err != nil {
			break
		}
	}
}

func parallelRun(config *Config) {
	feed := feed.NewFeed()
	var wg sync.WaitGroup
	var mutex sync.Mutex
	cond := sync.NewCond(&mutex)
	context := SharedContext{mutex: &mutex, cond: cond, wg: &wg, finish: false}
	queue := queue.NewLockFreeQueue()
	for i := 0; i < config.ConsumersCount; i++ {
		wg.Add(1)
		go consumer(config, queue, &feed, &context)
	}
	wg.Add(1)
	go producer(config, queue, &context)
	wg.Wait()
}

func Run(config Config) {
	if config.Mode == "s" {
		sequentialRun(&config)
		return
	}
	parallelRun(&config)
}
