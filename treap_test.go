package gtreap

import (
	"bytes"
	"math/rand"
	"sort"
	"testing"
)

func stringCompare(a, b interface{}) int {
	return bytes.Compare([]byte(a.(string)), []byte(b.(string)))
}

func intCompare(a, b interface{}) int {
	return a.(int) - b.(int)
}

func TestTreap(t *testing.T) {
	x := NewTreap(stringCompare)
	if x == nil {
		t.Errorf("expected NewTreap to work")
	}

	tests := []struct {
		op  string
		val string
		pri int
		exp string
	}{
		{"get", "not-there", -1, "NIL"},
		{"ups", "a", 100, ""},
		{"get", "a", -1, "a"},
		{"ups", "b", 200, ""},
		{"get", "a", -1, "a"},
		{"get", "b", -1, "b"},
		{"ups", "c", 300, ""},
		{"get", "a", -1, "a"},
		{"get", "b", -1, "b"},
		{"get", "c", -1, "c"},
		{"get", "not-there", -1, "NIL"},
		{"ups", "a", 400, ""},
		{"get", "a", -1, "a"},
		{"get", "b", -1, "b"},
		{"get", "c", -1, "c"},
		{"get", "not-there", -1, "NIL"},
		{"del", "a", -1, ""},
		{"get", "a", -1, "NIL"},
		{"get", "b", -1, "b"},
		{"get", "c", -1, "c"},
		{"get", "not-there", -1, "NIL"},
		{"ups", "a", 10, ""},
		{"get", "a", -1, "a"},
		{"get", "b", -1, "b"},
		{"get", "c", -1, "c"},
		{"get", "not-there", -1, "NIL"},
		{"del", "a", -1, ""},
		{"del", "b", -1, ""},
		{"del", "c", -1, ""},
		{"get", "a", -1, "NIL"},
		{"get", "b", -1, "NIL"},
		{"get", "c", -1, "NIL"},
		{"get", "not-there", -1, "NIL"},
		{"del", "a", -1, ""},
		{"del", "b", -1, ""},
		{"del", "c", -1, ""},
		{"get", "a", -1, "NIL"},
		{"get", "b", -1, "NIL"},
		{"get", "c", -1, "NIL"},
		{"get", "not-there", -1, "NIL"},
		{"ups", "a", 10, ""},
		{"get", "a", -1, "a"},
		{"get", "b", -1, "NIL"},
		{"get", "c", -1, "NIL"},
		{"get", "not-there", -1, "NIL"},
		{"ups", "b", 1000, "b"},
		{"del", "b", -1, ""}, // cover join that is nil
		{"ups", "b", 20, "b"},
		{"ups", "c", 12, "c"},
		{"del", "b", -1, ""}, // cover join second return
		{"ups", "a", 5, "a"}, // cover upsert existing with lower priority
	}

	for testIdx, test := range tests {
		switch test.op {
		case "get":
			i := x.Get(test.val)
			if i != test.exp && !(i == nil && test.exp == "NIL") {
				t.Errorf("test: %v, on Get, expected: %v, got: %v", testIdx, test.exp, i)
			}
		case "ups":
			x = x.Upsert(test.val, test.pri)
		case "del":
			x = x.Delete(test.val)
		}
	}
}

func load(x *Treap, arr []string) *Treap {
	for i, s := range arr {
		x = x.Upsert(s, i)
	}
	return x
}

func visitExpect(t *testing.T, x *Treap, start Item, arr []string) {
	n := 0
	x.VisitAscend(start, func(i Item) bool {
		if i.(string) != arr[n] {
			t.Errorf("expected visit item: %v, saw: %v", arr[n], i)
		}
		n++
		return true
	})
	if n != len(arr) {
		t.Errorf("expected # visit callbacks: %v, saw: %v", len(arr), n)
	}
}

func TestVisit(t *testing.T) {
	x := NewTreap(stringCompare)
	visitExpect(t, x, "a", []string{})

	x = load(x, []string{"e", "d", "c", "c", "a", "b", "a"})

	visitX := func() {
		visitExpect(t, x, "a", []string{"a", "b", "c", "d", "e"})
		visitExpect(t, x, nil, []string{"a", "b", "c", "d", "e"})
		visitExpect(t, x, "a1", []string{"b", "c", "d", "e"})
		visitExpect(t, x, "b", []string{"b", "c", "d", "e"})
		visitExpect(t, x, "b1", []string{"c", "d", "e"})
		visitExpect(t, x, "c", []string{"c", "d", "e"})
		visitExpect(t, x, "c1", []string{"d", "e"})
		visitExpect(t, x, "d", []string{"d", "e"})
		visitExpect(t, x, "d1", []string{"e"})
		visitExpect(t, x, "e", []string{"e"})
		visitExpect(t, x, "f", []string{})
	}
	visitX()

	var y *Treap
	y = x.Upsert("f", 1)
	y = y.Delete("a")
	y = y.Upsert("cc", 2)
	y = y.Delete("c")

	visitExpect(t, y, "a", []string{"b", "cc", "d", "e", "f"})
	visitExpect(t, y, nil, []string{"b", "cc", "d", "e", "f"})
	visitExpect(t, y, "a1", []string{"b", "cc", "d", "e", "f"})
	visitExpect(t, y, "b", []string{"b", "cc", "d", "e", "f"})
	visitExpect(t, y, "b1", []string{"cc", "d", "e", "f"})
	visitExpect(t, y, "c", []string{"cc", "d", "e", "f"})
	visitExpect(t, y, "c1", []string{"cc", "d", "e", "f"})
	visitExpect(t, y, "d", []string{"d", "e", "f"})
	visitExpect(t, y, "d1", []string{"e", "f"})
	visitExpect(t, y, "e", []string{"e", "f"})
	visitExpect(t, y, "f", []string{"f"})
	visitExpect(t, y, "z", []string{})

	// an uninitialized treap
	z := NewTreap(stringCompare)

	// a treap to force left traversal of min
	lmt := NewTreap(stringCompare)
	lmt = lmt.Upsert("b", 2)
	lmt = lmt.Upsert("a", 1)

	// The x treap should be unchanged.
	visitX()

	if x.Min() != "a" {
		t.Errorf("expected min of a")
	}
	if x.Max() != "e" {
		t.Errorf("expected max of d")
	}
	if y.Min() != "b" {
		t.Errorf("expected min of b")
	}
	if y.Max() != "f" {
		t.Errorf("expected max of f")
	}
	if z.Min() != nil {
		t.Errorf("expected min of nil")
	}
	if z.Max() != nil {
		t.Error("expected max of nil")
	}
	if lmt.Min() != "a" {
		t.Errorf("expected min of a")
	}
	if lmt.Max() != "b" {
		t.Errorf("expeced max of b")
	}
}

func visitExpectEndAtC(t *testing.T, x *Treap, start string, arr []string) {
	n := 0
	x.VisitAscend(start, func(i Item) bool {
		if stringCompare(i, "c") > 0 {
			return false
		}
		if i.(string) != arr[n] {
			t.Errorf("expected visit item: %v, saw: %v", arr[n], i)
		}
		n++
		return true
	})
	if n != len(arr) {
		t.Errorf("expected # visit callbacks: %v, saw: %v", len(arr), n)
	}
}

func TestVisitEndEarly(t *testing.T) {
	x := NewTreap(stringCompare)
	visitExpectEndAtC(t, x, "a", []string{})

	x = load(x, []string{"e", "d", "c", "c", "a", "b", "a", "e"})

	visitX := func() {
		visitExpectEndAtC(t, x, "a", []string{"a", "b", "c"})
		visitExpectEndAtC(t, x, "a1", []string{"b", "c"})
		visitExpectEndAtC(t, x, "b", []string{"b", "c"})
		visitExpectEndAtC(t, x, "b1", []string{"c"})
		visitExpectEndAtC(t, x, "c", []string{"c"})
		visitExpectEndAtC(t, x, "c1", []string{})
		visitExpectEndAtC(t, x, "d", []string{})
		visitExpectEndAtC(t, x, "d1", []string{})
		visitExpectEndAtC(t, x, "e", []string{})
		visitExpectEndAtC(t, x, "f", []string{})
	}
	visitX()

}

func visitDescendExpect(t *testing.T, x *Treap, start Item, arr []string) {
	n := 0
	x.VisitDescend(start, func(i Item) bool {
		if i.(string) != arr[n] {
			t.Errorf("expected visit item: %v, saw: %v (%v)", arr[n], i, arr)
		}
		n++
		return true
	})
	if n != len(arr) {
		t.Errorf("expected # visit callbacks: %v, saw: %v", len(arr), n)
	}
}

func TestVisitDescend(t *testing.T) {
	x := NewTreap(stringCompare)
	visitDescendExpect(t, x, "a", []string{})

	x = load(x, []string{"e", "d", "c", "c", "a", "b", "a"})

	visitX := func() {
		visitDescendExpect(t, x, "a", []string{"a"})
		visitDescendExpect(t, x, "a1", []string{"a"})
		visitDescendExpect(t, x, "b", []string{"b", "a"})
		visitDescendExpect(t, x, "b1", []string{"b", "a"})
		visitDescendExpect(t, x, "c", []string{"c", "b", "a"})
		visitDescendExpect(t, x, "c1", []string{"c", "b", "a"})
		visitDescendExpect(t, x, "d", []string{"d", "c", "b", "a"})
		visitDescendExpect(t, x, "d1", []string{"d", "c", "b", "a"})
		visitDescendExpect(t, x, "e", []string{"e", "d", "c", "b", "a"})
		visitDescendExpect(t, x, "f", []string{"e", "d", "c", "b", "a"})
		visitDescendExpect(t, x, nil, []string{"e", "d", "c", "b", "a"})
	}
	visitX()

	var y *Treap
	y = x.Upsert("f", 1)
	y = y.Delete("a")
	y = y.Upsert("cc", 2)
	y = y.Delete("c")

	visitDescendExpect(t, y, "a", []string{})
	visitDescendExpect(t, y, "a1", []string{})
	visitDescendExpect(t, y, "b", []string{"b"})
	visitDescendExpect(t, y, "b1", []string{"b"})
	visitDescendExpect(t, y, "c", []string{"b"})
	visitDescendExpect(t, y, "c1", []string{"b"})
	visitDescendExpect(t, y, "cd", []string{"cc", "b"})
	visitDescendExpect(t, y, "d", []string{"d", "cc", "b"})
	visitDescendExpect(t, y, "d1", []string{"d", "cc", "b"})
	visitDescendExpect(t, y, "e", []string{"e", "d", "cc", "b"})
	visitDescendExpect(t, y, "f", []string{"f", "e", "d", "cc", "b"})
	visitDescendExpect(t, y, "z", []string{"f", "e", "d", "cc", "b"})
	visitDescendExpect(t, y, nil, []string{"f", "e", "d", "cc", "b"})

	// The x treap should be unchanged.
	visitX()

	if x.Min() != "a" {
		t.Errorf("expected min of a")
	}
	if x.Max() != "e" {
		t.Errorf("expected max of d")
	}
	if y.Min() != "b" {
		t.Errorf("expected min of b")
	}
	if y.Max() != "f" {
		t.Errorf("expected max of f")
	}
}

func TestTreapCount(t *testing.T) {
	x := NewTreap(intCompare)
	r := rand.New(rand.NewSource(1984))
	vals := make([]int, 1000)

	for i := 0; i < len(vals); i++ {
		vals[i] = r.Int()
		x = x.Upsert(vals[i], r.Int())
		if x.Len() != i+1 {
			t.Errorf("expected count after upsert to be %d but was %d", i+1, x.Len())
		}
	}

	// store a copy of the full treap before we start deleting things from it, we'll use this later
	fullTreap := x

	for i := len(vals); i > 0; i-- {
		// pick an index between 0 and i (i is sliding down so we will only ever pick values
		// we've never used before, since all the values we have used will be above i)
		idx := r.Intn(i)
		// now delete the selected value from the treap
		x = x.Delete(vals[idx])
		// this should cause the length of the treap to decrease by one (which means it should match i-1)
		if x.Len() != i-1 {
			t.Errorf("expected count after delete to be %d but was %d", i, x.Len())
		}
		// here we move (by swapping) the selected value to the end of the slice so that
		// it won't be picked again
		vals[idx], vals[i-1] = vals[i-1], vals[idx]
	}

	// now we sort the values so they are in the same order they would be in the treap. This is important
	// so that when we start splitting the treap we know how many items are to the left and right of a given
	// value. We also restore x to the full treap before we deleted everything.
	sort.Ints(vals)
	x = fullTreap

	for i := 0; i < len(vals); i++ {
		// randomly pick an index and split the treap on the value at that index
		idx := r.Intn(len(vals))
		left, middle, right := x.Split(vals[idx])
		// at this point because vals has been sorted, we know that the left treap should contain
		// all values less than vals[idx] (which is idx values) and the right treap should contain
		// all values greater than vals[idx] (which is len(vals)-idx-1), and that the middle treap
		// contains one value.
		if left.Len() != idx {
			t.Errorf("expected left count after split to be %d but was %d", idx, left.Len())
		}
		if middle.Len() != 1 {
			t.Errorf("expected middle count after split to be 1 but was %d", middle.Len())
		}
		if right.Len() != len(vals)-idx-1 {
			t.Errorf("expected right count after split to be %d but was %d", len(vals)-idx-1, right.Len())
		}
	}

	// finally, do some splits on values not in the treap. Note: the random seed above was selected such that
	// this section will not generate any duplicate values that are already in the treap.
	for i := 0; i < 50; i++ {
		// first we get a random value that we're going to split on, and then search through our values list
		// to find the index where our random value would go. This tells us how many values in the treap are
		// less than the random value and how many are greater (and thus what the left and right lens should be)
		val := r.Int()
		var idx int
		for idx = len(vals) - 1; idx > 0; idx-- {
			if vals[idx] < val {
				break
			}
		}

		// vals[idx] is now the largest value in vals that is less than our randomly selected value, which
		// means there are idx+1 values in the treap that are less than val, and len(vals)-idx-1 values
		// which are greater than it.
		left, middle, right := x.Split(val)
		if middle != nil {
			t.Errorf("expected middle to be nil but wasn't")
		}
		if left.Len() != idx+1 {
			t.Errorf("expected left count after split to be %d but was %d", idx, left.Len())
		}
		if right.Len() != len(vals)-idx-1 {
			t.Errorf("expected right count after split to be %d but was %d", len(vals)-idx-1, right.Len())
		}
	}
}
