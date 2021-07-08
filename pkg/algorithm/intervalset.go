package algorithm

import (
	"math"
	"sort"
	"time"
)

type TimeIntervalSet struct {
	containsNegInf bool
	changePoints []time.Time
}

type inOrder []time.Time
func (io inOrder) Len() int {return len(io)}
func (io inOrder) Less(i, j int) bool {return io[i].Before(io[j])}
func (io inOrder) Swap(i, j int) {io[i], io[j] = io[j], io[i]}

// NewTimeInterval returns a TimeIntervalSet for the left-closed, right-open range [start, end)
func NewTimeInterval(start time.Time, end time.Time) TimeIntervalSet {
	if !start.Before(end) {
		start, end = end, start
	}
	return TimeIntervalSet{false, []time.Time{start, end}}
}

var EmptyTimeIntervalSet=TimeIntervalSet{}
var FullTimeIntervalSet=TimeIntervalSet{containsNegInf: true}

func (tis TimeIntervalSet) TotalDuration() time.Duration {
	if tis.containsNegInf || len(tis.changePoints)&1==1 {
		// infinite set
		return time.Duration(math.MaxInt64)
	}
	var durationAccumulator time.Duration
	for i:=0;i<len(tis.changePoints);i+=2 {
		durationAccumulator+=tis.changePoints[i+1].Sub(tis.changePoints[i])
	}
	return durationAccumulator
}

func (tis TimeIntervalSet) Contains(point time.Time) bool {
	n:=sort.Search(len(tis.changePoints), func(idx int) bool {
		return tis.changePoints[idx].After(point)
	})
	return tis.containsNegInf==(n&1==0)
}

// MergeTimeIntervalSets creates a new TimeIntervalSet from two other TimeIntervalSets.  The predicate function is
// called on various test points derived from the input sets and should return true if the testpoint should be
// in the resulting set.
func MergeTimeIntervalSets(left TimeIntervalSet, right TimeIntervalSet, predicate func(testpoint time.Time) bool) TimeIntervalSet {
	var candidateChangePoints []time.Time
	candidateChangePoints = append(candidateChangePoints, left.changePoints...)
	candidateChangePoints = append(candidateChangePoints, right.changePoints...)
	sort.Sort(inOrder(candidateChangePoints))
	// uniquify
	out:=1
	for i:=1; i<len(candidateChangePoints); i++ {
		if !candidateChangePoints[out-1].Equal(candidateChangePoints[i]) {
			candidateChangePoints[out]= candidateChangePoints[i]
			out++
		}
	}
	candidateChangePoints = candidateChangePoints[:out]

	result := TimeIntervalSet{}
	if len(candidateChangePoints) == 0 {
		result.containsNegInf=predicate(time.Time{})
		return result
	}
	including:=predicate(candidateChangePoints[0].Add(time.Duration(-1)))
	result.containsNegInf=including
	for _, cp := range candidateChangePoints {
		if including != predicate(cp) {
			result.changePoints = append(result.changePoints, cp)
			including = !including
		}
	}
	return result
}

func UnionTimeIntervalSets(sets... TimeIntervalSet) TimeIntervalSet {
	result := EmptyTimeIntervalSet
	for _, set := range sets {
		result = MergeTimeIntervalSets(result, set, func(testpoint time.Time) bool {return result.Contains(testpoint) || set.Contains(testpoint)})
	}
	return result
}

func IntersectTimeIntervalSets(sets... TimeIntervalSet) TimeIntervalSet {
	result := FullTimeIntervalSet
	for _, set := range sets {
		result = MergeTimeIntervalSets(result, set, func(testpoint time.Time) bool {return result.Contains(testpoint) && set.Contains(testpoint)})
	}
	return result
}