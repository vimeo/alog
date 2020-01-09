package internal

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestItoa(t *testing.T) {
	b := &bytes.Buffer{}
	Itoa(b, math.MaxUint32)
	got := b.String()
	want := "4294967295"
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v\n", got, want)
	}
}

type nonConcurrentWriter struct {
	w   io.Writer
	cas uint64
}

func (w *nonConcurrentWriter) Write(b []byte) (int, error) {
	if !atomic.CompareAndSwapUint64(&w.cas, 0, 1) {
		panic("unsynchronized entry")
	}
	defer atomic.CompareAndSwapUint64(&w.cas, 1, 0)

	time.Sleep(10 * time.Millisecond)

	return w.w.Write(b)
}

func writeStuffConcurrently(w io.Writer, size int, count int) (err error) {
	var wg sync.WaitGroup
	b := make([]byte, size)
	wg.Add(count)
	var m sync.Mutex
	setError := func(p interface{}) {
		m.Lock()
		defer m.Unlock()
		err = fmt.Errorf("%v", p)
	}
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			defer func() {
				if p := recover(); p != nil {
					setError(p)
				}
			}()
			w.Write(b)
		}()
	}
	wg.Wait()

	return
}

func TestSerializedWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewSerializedWriter(&nonConcurrentWriter{w: b})

	size := 4
	count := 5
	err := writeStuffConcurrently(w, size, count)
	if err != nil {
		t.Error(err.Error())
	}

	got := b.String()
	want := size * count
	if len(got) != want {
		t.Errorf("got:\n%d\nwant:\n%d\n", len(got), want)
	}
}

func TestNonConcurrentWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := &nonConcurrentWriter{w: b}

	size := 4
	count := 5
	err := writeStuffConcurrently(w, size, count)
	if err == nil {
		t.Errorf("nonConcurrentWriter did not detect concurrent write")
	}
}
