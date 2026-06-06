package service

import (
	"sync"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
)

type GenerationEventBus struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan domain.GenerationEvent]struct{}
}

func NewGenerationEventBus() *GenerationEventBus {
	return &GenerationEventBus{subscribers: make(map[string]map[chan domain.GenerationEvent]struct{})}
}

func (b *GenerationEventBus) Subscribe(projectID, jobID string) (<-chan domain.GenerationEvent, func()) {
	key := b.key(projectID, jobID)
	ch := make(chan domain.GenerationEvent, 32)

	b.mu.Lock()
	if _, ok := b.subscribers[key]; !ok {
		b.subscribers[key] = make(map[chan domain.GenerationEvent]struct{})
	}
	b.subscribers[key][ch] = struct{}{}
	b.mu.Unlock()

	cancel := func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if subs, ok := b.subscribers[key]; ok {
			if _, ok := subs[ch]; ok {
				delete(subs, ch)
				close(ch)
			}
			if len(subs) == 0 {
				delete(b.subscribers, key)
			}
		}
	}

	return ch, cancel
}

func (b *GenerationEventBus) Publish(event domain.GenerationEvent) {
	key := b.key(event.ProjectID, event.JobID)

	b.mu.RLock()
	subs := b.subscribers[key]
	for ch := range subs {
		select {
		case ch <- event:
		default:
		}
	}
	b.mu.RUnlock()
}

func (b *GenerationEventBus) Close(projectID, jobID string) {
	key := b.key(projectID, jobID)
	b.mu.Lock()
	defer b.mu.Unlock()
	if subs, ok := b.subscribers[key]; ok {
		for ch := range subs {
			close(ch)
			delete(subs, ch)
		}
		delete(b.subscribers, key)
	}
}

func (b *GenerationEventBus) key(projectID, jobID string) string {
	return projectID + ":" + jobID
}
