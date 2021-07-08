package algorithm

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type TimeIntervalSet struct {
	UnboundedBelow bool
	ChangePoints   []time.Time
}

type inOrder []time.Time

func (io inOrder) Len() int           { return len(io) }
func (io inOrder) Less(i, j int) bool { return io[i].Before(io[j]) }
func (io inOrder) Swap(i, j int)      { io[i], io[j] = io[j], io[i] }

// NewTimeInterval returns a TimeIntervalSet for the left-closed, right-open range [start, end)
func NewTimeInterval(start time.Time, end time.Time) TimeIntervalSet {
	if !start.Before(end) {
		start, end = end, start
	}
	return TimeIntervalSet{false, []time.Time{start, end}}
}

func (tis TimeIntervalSet) Complement() TimeIntervalSet {
	return TimeIntervalSet{
		UnboundedBelow: !tis.UnboundedBelow,
		ChangePoints:   append([]time.Time{}, tis.ChangePoints...),
	}
}

func (tis TimeIntervalSet) String() string {
	sb := strings.Builder{}
	if tis.UnboundedBelow {
		sb.WriteString("(-inf")
	}
	for i := 0; i < len(tis.ChangePoints); i += 2 {
		if tis.UnboundedBelow {
			sb.WriteString(fmt.Sprintf(",%v)[%v", tis.ChangePoints[i], tis.ChangePoints[i+1]))
		} else {
			sb.WriteString(fmt.Sprintf("[%v,%v)", tis.ChangePoints[i], tis.ChangePoints[i+1]))
		}
	}
	if len(tis.ChangePoints)&1 == 0 {
		if tis.UnboundedBelow {
			sb.WriteString(",inf)")
		}
	} else {
		if tis.UnboundedBelow {
			sb.WriteString(fmt.Sprintf(", %v)", tis.ChangePoints[len(tis.ChangePoints)-1]))
		} else {
			sb.WriteString(fmt.Sprintf("[%v, inf)", tis.ChangePoints[len(tis.ChangePoints)-1]))
		}
	}
	return sb.String()
}

var EmptyTimeIntervalSet = TimeIntervalSet{}
var FullTimeIntervalSet = TimeIntervalSet{UnboundedBelow: true}

func (tis TimeIntervalSet) TotalDuration() time.Duration {
	if tis.UnboundedBelow || len(tis.ChangePoints)&1 == 1 {
		// infinite set
		return time.Duration(math.MaxInt64)
	}
	var durationAccumulator time.Duration
	for i := 0; i < len(tis.ChangePoints); i += 2 {
		durationAccumulator += tis.ChangePoints[i+1].Sub(tis.ChangePoints[i])
	}
	return durationAccumulator
}

func (tis TimeIntervalSet) Contains(point time.Time) bool {
	n := sort.Search(len(tis.ChangePoints), func(idx int) bool {
		return tis.ChangePoints[idx].After(point)
	})
	return tis.UnboundedBelow == (n&1 == 0)
}

// MergeTimeIntervalSets creates a new TimeIntervalSet from two other TimeIntervalSets.  The predicate function is
// called on various test points derived from the input sets and should return true if the testpoint should be
// in the resulting set.
func MergeTimeIntervalSets(left TimeIntervalSet, right TimeIntervalSet, predicate func(testpoint time.Time) bool) TimeIntervalSet {
	result := TimeIntervalSet{}

	var candidateChangePoints []time.Time
	candidateChangePoints = append(candidateChangePoints, left.ChangePoints...)
	candidateChangePoints = append(candidateChangePoints, right.ChangePoints...)

	if len(candidateChangePoints) == 0 {
		result.UnboundedBelow = predicate(time.Time{})
		return result
	}

	sort.Sort(inOrder(candidateChangePoints))

	// uniquify
	out := 1
	for i := 1; i < len(candidateChangePoints); i++ {
		if !candidateChangePoints[out-1].Equal(candidateChangePoints[i]) {
			candidateChangePoints[out] = candidateChangePoints[i]
			out++
		}
	}
	candidateChangePoints = candidateChangePoints[:out]

	including := predicate(candidateChangePoints[0].Add(time.Duration(-1)))
	result.UnboundedBelow = including
	for _, cp := range candidateChangePoints {
		if including != predicate(cp) {
			result.ChangePoints = append(result.ChangePoints, cp)
			including = !including
		}
	}
	return result
}

func UnionTimeIntervalSets(sets ...TimeIntervalSet) TimeIntervalSet {
	result := EmptyTimeIntervalSet
	for _, set := range sets {
		result = MergeTimeIntervalSets(result, set, func(testpoint time.Time) bool { return result.Contains(testpoint) || set.Contains(testpoint) })
	}
	return result
}

func IntersectTimeIntervalSets(sets ...TimeIntervalSet) TimeIntervalSet {
	result := FullTimeIntervalSet
	for _, set := range sets {
		result = MergeTimeIntervalSets(result, set, func(testpoint time.Time) bool { return result.Contains(testpoint) && set.Contains(testpoint) })
	}
	return result
}
