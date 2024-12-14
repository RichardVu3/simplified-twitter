package feed

import (
	"simplified-twitter/lock"
)

type Feed interface {
	Add(body string, timestamp float64)
	Remove(timestamp float64) bool
	Contains(timestamp float64) bool
	GetAllFeeds() []FeedPost
}

type feed struct {
	start *post
	lock  lock.RWLock
}

type post struct {
	body      string
	timestamp float64
	next      *post
}

func newPost(body string, timestamp float64, next *post) *post {
	return &post{body, timestamp, next}
}

func NewFeed() Feed {
	return &feed{start: nil, lock: lock.NewRWMutex()}
}

func (f *feed) Add(body string, timestamp float64) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.start == nil {
		f.start = newPost(body, timestamp, nil)
		return
	}
	if f.start.timestamp < timestamp {
		f.start = newPost(body, timestamp, f.start)
		return
	}
	for p := f.start; ; p = p.next {
		if p.next == nil || p.next.timestamp < timestamp {
			p.next = newPost(body, timestamp, p.next)
			return
		}
	}
}

func (f *feed) Remove(timestamp float64) bool {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.start == nil {
		return false
	}
	if f.start.timestamp == timestamp {
		f.start = f.start.next
		return true
	}
	for p := f.start; p.next != nil; p = p.next {
		if p.next.timestamp == timestamp {
			p.next = p.next.next
			return true
		}
	}
	return false
}

func (f *feed) Contains(timestamp float64) bool {
	f.lock.RLock()
	defer f.lock.RUnlock()

	for p := f.start; p != nil; p = p.next {
		if p.timestamp == timestamp {
			return true
		}
	}
	return false
}

type FeedPost struct {
	Body      string  `json:"body"`
	Timestamp float64 `json:"timestamp"`
}

func (f *feed) GetAllFeeds() []FeedPost {
	f.lock.RLock()
	defer f.lock.RUnlock()

	var posts []FeedPost
	for p := f.start; p != nil; p = p.next {
		posts = append(posts, FeedPost{Body: p.body, Timestamp: p.timestamp})
	}
	return posts
}
