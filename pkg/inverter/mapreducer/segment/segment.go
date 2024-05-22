package segment

import (
	"sync"
)

type (
	Pair struct {
		Term string
		Doc  int
	}

	Fragment struct {
		Pairs []Pair
	}

	Segment struct {
		Section []*Fragment
	}


	SegmentManager struct {
		Segments []*Segment
		mu sync.Mutex
	}
)


func NewSegmentManager(splits int) *SegmentManager{
	manager := SegmentManager{
		Segments: make([]*Segment, splits),
	}
	return &manager
}

func (sm *SegmentManager) AddSegment(segment *Segment) {
	sm.mu.Lock()
	sm.Segments = append(sm.Segments, segment)
	sm.mu.Unlock()
}