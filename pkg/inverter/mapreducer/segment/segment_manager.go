package segment

import "sync"



type SegmentManager struct {
	Segments []*Segment
	mu sync.Mutex
}


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

/* suppose that segments for the array are defined as follows:
array = [ [frag11] [frag12] [frag13] frag[14] frag[15] ]
		| [frag21] [frag22] [frag23] frag[24] frag[25] |
		| [frag31] [frag32] [frag33] frag[34] frag[35] |
		[ [frag41] [frag42] [frag43] frag[44] frag[45] ]


now, calling getVerticalSegment returns an 1d array of fragments of column "num"


example: GetVerticalSegment(2) would return [[frag12] [frag22] [frag32] [frag42]]

this actually maintains hierarchial implementation so that inverter workers could have a single 
Segment to parse instead of list of segments that not all of them are used

*/
func (sm *SegmentManager) GetVerticalSegment(num int) *Segment{
if num >= len(sm.Segments[0].Fragments) {
	return nil
}
segment := Segment{}
for i := 0; i < len(sm.Segments); i++ {
	segment.Fragments = append(segment.Fragments, sm.Segments[i].Fragments[num]) 
}
return &segment
}