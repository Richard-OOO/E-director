package service

import (
	"sync"

	"github.com/Richard-OOO/E-director/apps/backend/internal/domain"
)

const generationEventHistoryLimit = 64

type GenerationEventBus struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan domain.GenerationEvent]struct{}
	history     map[string][]domain.GenerationEvent
}

func NewGenerationEventBus() *GenerationEventBus {
	return &GenerationEventBus{
		subscribers: make(map[string]map[chan domain.GenerationEvent]struct{}),
		history:     make(map[string][]domain.GenerationEvent),
	}
}

func (b *GenerationEventBus) Subscribe(projectID, jobID string) (<-chan domain.GenerationEvent, func()) {
	key := b.key(projectID, jobID)
	ch := make(chan domain.GenerationEvent, 32)

	b.mu.Lock()
	if _, ok := b.subscribers[key]; !ok {
		b.subscribers[key] = make(map[chan domain.GenerationEvent]struct{})
	}
	b.subscribers[key][ch] = struct{}{}
	for _, event := range b.history[key] {
		ch <- event
	}
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

	b.mu.Lock()
	b.history[key] = append(b.history[key], event)
	if len(b.history[key]) > generationEventHistoryLimit {
		b.history[key] = b.history[key][len(b.history[key])-generationEventHistoryLimit:]
	}
	for ch := range b.subscribers[key] {
		select {
		case ch <- event:
		default:
		}
	}
	b.mu.Unlock()
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
