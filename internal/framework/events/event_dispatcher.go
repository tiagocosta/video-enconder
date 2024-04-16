package events

import (
	"sync"
)

type EventDispatcher struct {
	Handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		Handlers: make(map[string][]EventHandlerInterface),
	}
}

func (em *EventDispatcher) Dispatch(event EventInterface) error {
	if handlers, ok := em.Handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

func (em *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := em.Handlers[eventName]; ok {
		for _, h := range em.Handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	em.Handlers[eventName] = append(em.Handlers[eventName], handler)
	return nil
}

func (em *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := em.Handlers[eventName]; ok {
		for _, h := range em.Handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (em *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, ok := em.Handlers[eventName]; ok {
		for i, h := range em.Handlers[eventName] {
			if h == handler {
				em.Handlers[eventName] = append(em.Handlers[eventName][:i], em.Handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (em *EventDispatcher) Clear() {
	em.Handlers = make(map[string][]EventHandlerInterface)
}
