package multipattern

import "testing"

func TestMultipattern(t *testing.T) {
	mp := New()
	mp.On("ab", func(strings []string, _ interface{}) interface{} {
		t.Log(strings)
		return 1
	}).On("b(c)", func(strings []string, _ interface{}) interface{} {
		t.Log(strings)
		return 2
	}).On("(d)(e)", func(strings []string, _ interface{}) interface{} {
		t.Log(strings)
		return 3
	})

	if mp.Dispatch("ab", nil).(int) != 1 {
		t.Fatal("Expected match")
	}

	if mp.Dispatch("bc", nil).(int) != 2 {
		t.Fatal("Expected match")
	}

	if mp.Dispatch("de", nil).(int) != 3 {
		t.Fatal("Expected match")
	}

	if mp.Dispatch("fg", nil) != nil {
		t.Fatal("Expected no match")
	}
}
