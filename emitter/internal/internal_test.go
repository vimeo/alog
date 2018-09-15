package internal

import (
	"bytes"
	"math"
	"testing"
	"unsafe"
)

func TestPool(t *testing.T) {
	b1 := GetBuffer()
	bp1 := uintptr(unsafe.Pointer(b1))
	b2 := GetBuffer()
	bp2 := uintptr(unsafe.Pointer(b2))
	PutBuffer(b1)
	b := GetBuffer()
	bp := uintptr(unsafe.Pointer(b))
	if bp != bp1 {
		t.Errorf("%v != %v", bp1, bp2)
	}
	PutBuffer(b2)
	b = GetBuffer()
	bp = uintptr(unsafe.Pointer(b))
	if bp != bp2 {
		t.Errorf("%v != %v", bp1, bp2)
	}
}

func TestItoa(t *testing.T) {
	b := &bytes.Buffer{}
	Itoa(b, math.MaxUint32)
	got := b.String()
	want := "4294967295"
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v\n", got, want)
	}
}

func TestSerializedWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewSerializedWriter(b)
	w.Write([]byte("The quick brown fox"))
	got := b.String()
	want := "The quick brown fox"
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v\n", got, want)
	}
}
