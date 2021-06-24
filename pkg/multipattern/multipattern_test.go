package multipattern

import "testing"

func TestMultipattern(t *testing.T) {
	mp := New()
	mp.On("ab", func(strings []string) interface{} {
		t.Log(strings)
		return 1
	}).On("b(c)", func(strings []string) interface{} {
		t.Log(strings)
		return 2
	}).On("(d)(e)", func(strings []string) interface{} {
		t.Log(strings)
		return 3
	})

	if mp.Dispatch("ab").(int) != 1 {
		t.Fatal("Expected match")
	}

	if mp.Dispatch("bc").(int) != 2 {
		t.Fatal("Expected match")
	}

	if mp.Dispatch("de").(int) != 3 {
		t.Fatal("Expected match")
	}

	if mp.Dispatch("fg") != nil {
		t.Fatal("Expected no match")
	}
}
