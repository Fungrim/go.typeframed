package io

import (
	"testing"
	"bytes"
)

func TestBufferGrowth(t *testing.T) {
	b := newBufferWithLength(4)
	b.ensure(12)
	b.position = 2
	if cap(b.buf) != 20 {
		t.Errorf("wrong capacity %v", cap(b.buf))
	}
	if b.remaining() != 18 {
		t.Errorf("wrong remaining count %v", b.remaining())
	}
}

func TestWriteVarintAndGrow(t *testing.T) {
	b := newBufferWithLength(0)
	b.writeVarInt(4)
	if cap(b.buf) != 5 {
		t.Errorf("wrong capacity %v", cap(b.buf))
	}
	if b.position != 1 {
		t.Errorf("wrong position %v", b.position)
	}
}

func TestWriteArrayAndGrow(t *testing.T) {
	b := newBufferWithLength(0)
	b.writeVarInt(0) // grow 5
	buff := []byte{1, 2, 3, 4}
	b.write(buff)
	if cap(b.buf) != 5 {
		t.Errorf("wrong capacity %v", cap(b.buf))
	}
	if b.remaining() != 0 {
		t.Errorf("wrong remaining count %v", b.remaining())
	}
	if !bytes.Equal(b.buf, []byte{0, 1, 2, 3, 4}) {
		t.Fail()
	}
}