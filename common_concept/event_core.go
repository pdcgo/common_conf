package common_concept

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
)

type Subscriber struct {
	Ctx    context.Context
	Key    string
	Chan   chan interface{}
	Cancel func() error
}

func (sub *Subscriber) Emit(data interface{}) {
	select {
	case sub.Chan <- data:
		return
	case <-sub.Ctx.Done():
		return
	}
}

func NewSubscriber() *Subscriber {
	key := uuid.New()
	return &Subscriber{
		Key:  key.String(),
		Chan: make(chan interface{}),
	}
}

type CloseEvent struct{}

type CoreEvent struct {
	sync.Mutex
	Input      chan interface{}
	Subscriber []*Subscriber
}

func (ev *CoreEvent) Emit(event interface{}) {
	go func() {
		ev.Input <- event
	}()

}

func (ev *CoreEvent) CreateSubscriber() (sub *Subscriber) {
	ev.Lock()
	defer ev.Unlock()

	ctx, cancelCtx := context.WithCancel(context.TODO())

	newsub := NewSubscriber()
	newsub.Ctx = ctx
	ev.Subscriber = append(ev.Subscriber, newsub)
	newsub.Cancel = func() error {
		ev.Lock()
		defer ev.Unlock()

		newsubs := []*Subscriber{}

		for _, sub := range ev.Subscriber {
			if newsub.Key == sub.Key {
				cancelCtx()
				continue
			}

			newsubs = append(newsubs, sub)
		}
		ev.Subscriber = newsubs
		return nil
	}

	return newsub
}

func (ev *CoreEvent) GetEvent() <-chan interface{} {
	sub := ev.CreateSubscriber()
	return sub.Chan
}

func (ev *CoreEvent) Close() {

	for _, sub := range ev.Subscriber {
		sub.Cancel()
	}
}

func NewCoreEvent() *CoreEvent {
	log.Println("creating channel event")

	core := CoreEvent{
		Input:      make(chan interface{}),
		Subscriber: []*Subscriber{},
	}

	go func() {
		c := 0
		for item := range core.Input {

			func() {
				core.Lock()
				defer core.Unlock()
				itemw := item
				for _, channel := range core.Subscriber {
					subs := channel
					go subs.Emit(itemw)

				}
			}()
			c += 1
		}

		log.Println("close core event")
	}()

	return &core
}
