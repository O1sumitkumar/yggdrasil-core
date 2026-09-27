package events

import "testing"

func TestBusPublishSubscribe(t *testing.T) {
	bus := NewBus(8)
	id, ch := bus.Subscribe()
	defer bus.Unsubscribe(id)

	bus.Publish(New(TaskCreated, map[string]any{"task_id": "t1"}))

	evt := <-ch
	if evt.Type != TaskCreated {
		t.Fatalf("expected %s, got %s", TaskCreated, evt.Type)
	}
	if evt.Payload["task_id"] != "t1" {
		t.Fatalf("unexpected payload: %#v", evt.Payload)
	}
	if evt.ID == "" {
		t.Fatal("expected event ID")
	}
}

func TestBusUnsubscribe(t *testing.T) {
	bus := NewBus(2)
	id, ch := bus.Subscribe()
	bus.Unsubscribe(id)
	_, ok := <-ch
	if ok {
		t.Fatal("expected closed channel")
	}
}
