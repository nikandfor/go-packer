package packer

import (
	"bytes"
	"encoding/gob"
	"math/rand"
	"testing"
)

func TestInt(t *testing.T) {
	var v int

	for i := 0; i < 10000; i++ {
		v = rand.Int()

		data, err := Marshal(v)
		if err != nil {
			t.Errorf("marshal error for value %v: %v", v, err)
		}

		var res int
		err = Unmarshal(data, &res)
		if err != nil {
			t.Errorf("unmarshal error for value %v: %v. data: %d [%x]", v, err, len(data), data)
		}

		if v != res {
			t.Errorf("marshal/unmarshal error %v become %v. data: %d [%x]", v, res, len(data), data)
		}
	}
}

func BenchmarkPackerIntEncode(t *testing.B) {
	var buf bytes.Buffer
	w := NewEncoder(&buf)
	r := NewDecoder(&buf)

	var v int = 100000000000000000
	var res int

	err := w.Encode(v)
	if err != nil {
		t.Errorf("encode %v error: %v", v, err)
	}

	s := buf.Len()
	//	t.Logf("packed size: %v", s)

	t.ReportAllocs()
	t.SetBytes(int64(s))

	err = r.Decode(&res)
	if err != nil {
		t.Errorf("encode %v error: %v", v, err)
	}

	if v != res {
		t.Errorf("encode/decode error %v become %v", v, res)
	}

	for i := 0; i < t.N; i++ {
		_ = w.Encode(v)
		_ = r.Decode(&res)
		buf.Reset()
	}
}

func BenchmarkPackerIntMarshal(t *testing.B) {
	var v int = 100000000000000000
	var res int

	data, err := Marshal(v)
	if err != nil {
		t.Errorf("encode %v error: %v", v, err)
	}

	s := len(data)
	//	t.Logf("packed size: %v", s)

	t.ReportAllocs()
	t.SetBytes(int64(s))

	err = Unmarshal(data, &res)
	if err != nil {
		t.Errorf("encode %v error: %v", v, err)
	}

	if v != res {
		t.Errorf("encode/decode error %v become %v", v, res)
	}

	for i := 0; i < t.N; i++ {
		_, _ = Marshal(v)
		_ = Unmarshal(data, &res)
	}
}

func BenchmarkGobInt(t *testing.B) {
	var buf bytes.Buffer
	w := gob.NewEncoder(&buf)
	r := gob.NewDecoder(&buf)

	var v int = 100000000000000000
	var res int

	err := w.Encode(v)
	if err != nil {
		t.Errorf("encode %v error: %v", v, err)
	}

	s := buf.Len()
	//	t.Logf("packed size: %v", s)

	t.ReportAllocs()
	t.SetBytes(int64(s))

	err = r.Decode(&res)
	if err != nil {
		t.Errorf("encode %v error: %v", v, err)
	}

	if v != res {
		t.Errorf("encode/decode error %v become %v", v, res)
	}

	for i := 0; i < t.N; i++ {
		_ = w.Encode(v)
		_ = r.Decode(&res)
		buf.Reset()
	}
}
