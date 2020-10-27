package minmaxcnt

import "testing"

// nolint
var actions = []struct {
	key    interface{}
	action string
	count  uint
}{
	{nil, "max", 0},
	{nil, "min", 0},
	{"cat", "increment", 1},
	{"dog", "increment", 1},
	{"cat", "decrement", 1},
	{"dog", "decrement", 1},
	{"foo", "increment", 3},
	{"bar", "increment", 4},
	{"bar", "max", 4},
	{"foo", "min", 3},
	{"foo", "increment", 2},
	{"foo", "max", 5},
	{"bar", "min", 4},
	{"foo", "decrement", 3},
	{"bar", "max", 4},
	{"foo", "min", 2},
	{"foo", "decrement", 1},
	{"bar", "max", 4},
	{"foo", "min", 1},
	{"foo", "decrement", 1},
	{"bar", "max", 4},
	{"bar", "min", 4},
	{"doesntexist", "decrement", 4},
	{"bar", "count", 4},
	{"doesntexist", "count", 0},
}

func TestMinMaxCnt(t *testing.T) {
	d := New()

	for _, a := range actions {
		switch a.action {
		case "increment":
			for i := uint(0); i < a.count; i++ {
				d.Increment(a.key)
			}
		case "decrement":
			for i := uint(0); i < a.count; i++ {
				d.Decrement(a.key)
			}
		case "min":
			if key, count := d.Min(); key != a.key || count != a.count {
				t.Errorf("Expected Min (%v,%v) got (%v,%v)", a.key, a.count, key, count)
			}
		case "max":
			if key, count := d.Max(); key != a.key || count != a.count {
				t.Errorf("Expected Max (%v,%v) got (%v,%v)", a.key, a.count, key, count)
			}
		case "count":
			if count := d.Count(a.key); count != a.count {
				t.Errorf("Expected Count %d, got Count %d", a.count, count)
			}
		}
	}
}

func BenchmarkIncrement(b *testing.B) {
	d := New()
	size := len(actions)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.Increment(actions[i%size].key)
	}
}

func BenchmarkDecrement(b *testing.B) {
	d := New()
	size := len(actions)

	for i := 0; i < 50_000; i++ {
		d.Increment(actions[i%size].key)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.Decrement(actions[i%size].key)
	}
}

func BenchmarkCount(b *testing.B) {
	d := New()
	size := len(actions)

	for i := 0; i < 50_000; i++ {
		d.Increment(actions[i%size].key)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.Count(actions[i%size].key)
	}
}

func BenchmarkMax(b *testing.B) {
	d := New()
	size := len(actions)

	for i := 0; i < 50_000; i++ {
		d.Increment(actions[i%size].key)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.Max()
	}
}

func BenchmarkMin(b *testing.B) {
	d := New()
	size := len(actions)

	for i := 0; i < 50_000; i++ {
		d.Increment(actions[i%size].key)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.Min()
	}
}
