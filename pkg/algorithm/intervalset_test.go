package algorithm

import (
	"testing"
	"time"
)

func timeList(count int) []time.Time {
	var result []time.Time
	result = append(result, time.Now())
	for i := 1; i < count; i++ {
		result = append(result, result[i-1].Add(time.Second))
	}
	return result
}

func mustContain(t *testing.T, tis TimeIntervalSet, point time.Time, name string) {
	if !tis.Contains(point) {
		t.Fatal("failed containment test ", name)
	}
}

func mustNotContain(t *testing.T, tis TimeIntervalSet, point time.Time, name string) {
	if tis.Contains(point) {
		t.Fatal("failed containment test ", name)
	}
}

func Test_TimeIntervalSetContains(t *testing.T) {
	times := timeList(5)

	is := NewTimeInterval(times[1], times[3])
	mustNotContain(t, is, times[0], "before")
	mustContain(t, is, times[1], "start")
	mustContain(t, is, times[2], "mid")
	mustNotContain(t, is, times[3], "end")
	mustNotContain(t, is, times[4], "after")

	// Negate set
	is.UnboundedBelow = true
	mustContain(t, is, times[0], "inv before")
	mustNotContain(t, is, times[1], "inv start")
	mustNotContain(t, is, times[2], "inv mid")
	mustContain(t, is, times[3], "inv end")
	mustContain(t, is, times[4], "inv after")
}

func Test_UnionTimeIntervalSets(t *testing.T) {
	times := timeList(9)
	// non overlapping
	is := UnionTimeIntervalSets(NewTimeInterval(times[1], times[3]), NewTimeInterval(times[5], times[7]))
	mustNotContain(t, is, times[0], "U1 before")
	mustContain(t, is, times[1], "U1 start left")
	mustContain(t, is, times[2], "U1 mid left")
	mustNotContain(t, is, times[3], "U1 end left")
	mustNotContain(t, is, times[4], "U1 between")
	mustContain(t, is, times[5], "U1 start right")
	mustContain(t, is, times[6], "U1 mid right")
	mustNotContain(t, is, times[7], "U1 end right")
	mustNotContain(t, is, times[8], "U1 after'")

	if is.TotalDuration() != 4*time.Second {
		t.Fatal("U1 Expected 4 seconds")
	}

	// overlapping, no shared endpoints
	is = UnionTimeIntervalSets(NewTimeInterval(times[1], times[3]), NewTimeInterval(times[2], times[4]))
	mustNotContain(t, is, times[0], "U2 before")
	mustContain(t, is, times[1], "U2 start left")
	mustContain(t, is, times[2], "U2 mid left / start right")
	mustContain(t, is, times[3], "U2 end left / mid right")
	mustNotContain(t, is, times[4], "U2 end right")
	mustNotContain(t, is, times[5], "U2 after")

	if is.TotalDuration() != 3*time.Second {
		t.Fatal("U2 Expected 3 seconds")
	}

	// overlapping, shared endpoint
	is = UnionTimeIntervalSets(NewTimeInterval(times[1], times[3]), NewTimeInterval(times[3], times[5]))
	mustNotContain(t, is, times[0], "U3 before")
	mustContain(t, is, times[1], "U3 start left")
	mustContain(t, is, times[2], "U3 mid left")
	mustContain(t, is, times[3], "U3 end left / start right")
	mustContain(t, is, times[4], "U3 mid right")
	mustNotContain(t, is, times[5], "U3 end right")
	mustNotContain(t, is, times[6], "U3 after")

	if is.TotalDuration() != 4*time.Second {
		t.Fatal("U3 Expected 4 seconds")
	}

	is = UnionTimeIntervalSets(EmptyTimeIntervalSet, EmptyTimeIntervalSet)
	if is.TotalDuration() != 0 {
		t.Fatal("U4 Expected 0 seconds")
	}
}

func TestIntersectTimeIntervalSets(t *testing.T) {
	times := timeList(9)

	// non overlapping
	is := IntersectTimeIntervalSets(NewTimeInterval(times[1], times[3]), NewTimeInterval(times[5], times[7]))
	mustNotContain(t, is, times[0], "I1 before")
	mustNotContain(t, is, times[1], "I1 start left")
	mustNotContain(t, is, times[2], "I1 mid left")
	mustNotContain(t, is, times[3], "I1 end left")
	mustNotContain(t, is, times[4], "I1 between")
	mustNotContain(t, is, times[5], "I1 start right")
	mustNotContain(t, is, times[6], "I1 mid right")
	mustNotContain(t, is, times[7], "I1 end right")
	mustNotContain(t, is, times[8], "I1 after'")

	if is.TotalDuration() != 0 {
		t.Fatal("I1 Expected 0 seconds")
	}

	// overlapping, no shared endpoints
	is = IntersectTimeIntervalSets(NewTimeInterval(times[1], times[3]), NewTimeInterval(times[2], times[4]))
	mustNotContain(t, is, times[0], "I2 before")
	mustNotContain(t, is, times[1], "I2 start left")
	mustContain(t, is, times[2], "I2 mid left / start right")
	mustNotContain(t, is, times[3], "I2 end left / mid right")
	mustNotContain(t, is, times[4], "I2 end right")
	mustNotContain(t, is, times[5], "I2 after")

	if is.TotalDuration() != 1*time.Second {
		t.Fatal("I2 Expected 1 seconds")
	}

	// overlapping, shared endpoint
	is = IntersectTimeIntervalSets(NewTimeInterval(times[1], times[3]), NewTimeInterval(times[3], times[5]))
	mustNotContain(t, is, times[0], "I3 before")
	mustNotContain(t, is, times[1], "I3 start left")
	mustNotContain(t, is, times[2], "I3 mid left")
	mustNotContain(t, is, times[3], "I3 end left / start right")
	mustNotContain(t, is, times[4], "I3 mid right")
	mustNotContain(t, is, times[5], "I3 end right")
	mustNotContain(t, is, times[6], "I3 after")

	if is.TotalDuration() != 0 {
		t.Fatal("I3 Expected 0 seconds")
	}
}

func TestString(t *testing.T) {
	is := NewTimeInterval(time.Time{}, time.Time{}.Add(1*time.Second))
	if "[0001-01-01 00:00:00 +0000 UTC,0001-01-01 00:00:01 +0000 UTC)" != is.String() {
		t.Fatal("Expected different value in TestString")
	}
}
