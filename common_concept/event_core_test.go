package common_concept_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/pdcgo/common_conf/common_concept"
	"github.com/stretchr/testify/assert"
)

type AccountItem struct {
	Username  string
	Connected bool
}

func TestCloseEvent(t *testing.T) {
	event := common_concept.NewCoreEvent()

	found := false

	go func() {
		time.Sleep(time.Second * 2)
		event.Close()

	}()

	timeout := time.After(time.Second * 5)

	sub := event.CreateSubscriber()
	ctx := sub.Ctx

	select {
	case <-ctx.Done():
		found = true
	case event := <-sub.Chan:
		t.Log("get event", event)
		t.Log(reflect.TypeOf(event))

	case <-timeout:
		break
	}

	assert.True(t, found)
}

func TestEventCo(t *testing.T) {
	event := common_concept.NewCoreEvent()

	found := false

	status := &AccountItem{
		Username:  "asdasdasd",
		Connected: true,
	}

	go func() {
		time.Sleep(time.Second * 2)
		event.Emit(status)

	}()

	timeout := time.After(time.Second * 5)

	ch := event.GetEvent()

	select {
	case event := <-ch:
		t.Log("get event", event)
		t.Log(reflect.TypeOf(event))
		switch ev := event.(type) {
		case *AccountItem:
			found = true
			t.Log("masuk case", ev)
		default:
			t.Log(reflect.TypeOf(ev))
		}
	case <-timeout:
		break
	}

	assert.True(t, found)
}
