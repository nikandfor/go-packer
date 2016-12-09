package packer

import (
	"bytes"
	"encoding/gob"
	"math/rand"
	"reflect"
	"testing"
)

func TestInt(t *testing.T) {
	var v int

	p := CreatePacker(v)

	for i := 0; i < 10000; i++ {
		v = rand.Int()

		size := p.Size(reflect.ValueOf(v))

		data, err := Marshal(v)
		if err != nil {
			t.Errorf("marshal error for value %v: %v", v, err)
		}

		if size != len(data) {
			t.Errorf("Size (%v) != len(data) (%v). data: [%x]", size, len(data), data)
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

func BenchmarkGob___Int(t *testing.B) {
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

func TestBytesShort(t *testing.T) {
	var v = []byte("QWEQEQqwdw001nfcudbaqwlkenli31lni")

	p := CreatePacker(v)

	data, err := Marshal(v)
	if err != nil {
		t.Errorf("marshal error for value %v: %v", v, err)
	}

	size := p.Size(reflect.ValueOf(v))
	if size != len(data) {
		t.Errorf("Size (%v) != len(data) (%v). data: [%x]", size, len(data), data)
	}

	var res []byte
	err = Unmarshal(data, &res)
	if err != nil {
		t.Errorf("unmarshal error for value %v: %v. data: %d [%x]", v, err, len(data), data)
	}

	if !bytes.Equal(v, res) {
		t.Errorf("marshal/unmarshal error %x become %x. data: %d [%x]", v, res, len(data), data)
	}
}

func TestBytesLong(t *testing.T) {
	var v = []byte("QWEQEQqwdw001nfcudbawwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwqelqwqwqowv0293022972vccccijks83792qcccckwlcnqwkenlqwc.c1092cqneocq2389fbwb9923928hp9dh8oqwg872gf38g2o28gqo8qwlkenli31lni")

	data, err := Marshal(v)
	if err != nil {
		t.Errorf("marshal error for value %v: %v", v, err)
	}

	var res []byte
	err = Unmarshal(data, &res)
	if err != nil {
		t.Errorf("unmarshal error for value %v: %v. data: %d [%x]", v, err, len(data), data)
	}

	if !bytes.Equal(v, res) {
		t.Errorf("marshal/unmarshal error %x become %x. data: %d [%x]", v, res, len(data), data)
	}
}

func BenchmarkPackerBytesEncode(t *testing.B) {
	var buf bytes.Buffer
	w := NewEncoder(&buf)
	r := NewDecoder(&buf)

	var v = []byte("QWEQEQqwdw001nfcudbaqwlkenli31lni")
	var res []byte

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

	if !bytes.Equal(v, res) {
		t.Errorf("marshal/unmarshal error %x become %x. data: %d [%x]", v, res, buf.Len(), buf.Bytes())
	}

	for i := 0; i < t.N; i++ {
		_ = w.Encode(v)
		_ = r.Decode(&res)
		buf.Reset()
	}
}

func BenchmarkPackerBytesLongEncode(t *testing.B) {
	var buf bytes.Buffer
	w := NewEncoder(&buf)
	r := NewDecoder(&buf)

	var v = make([]byte, 1024)
	var res []byte

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

	if !bytes.Equal(v, res) {
		t.Errorf("marshal/unmarshal error %x become %x. data: %d [%x]", v, res, buf.Len(), buf.Bytes())
	}

	for i := 0; i < t.N; i++ {
		_ = w.Encode(v)
		_ = r.Decode(&res)
		buf.Reset()
	}
}

func BenchmarkGob___BytesEncode(t *testing.B) {
	var buf bytes.Buffer
	w := gob.NewEncoder(&buf)
	r := gob.NewDecoder(&buf)

	var v = []byte("QWEQEQqwdw001nfcudbaqwlkenli31lni")
	var res []byte

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

	if !bytes.Equal(v, res) {
		t.Errorf("marshal/unmarshal error %x become %x. data: %d [%x]", v, res, buf.Len(), buf.Bytes())
	}

	for i := 0; i < t.N; i++ {
		_ = w.Encode(v)
		_ = r.Decode(&res)
		buf.Reset()
	}
}

func BenchmarkGob___BytesLongEncode(t *testing.B) {
	var buf bytes.Buffer
	w := gob.NewEncoder(&buf)
	r := gob.NewDecoder(&buf)

	var v = make([]byte, 1024)
	var res []byte

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

	if !bytes.Equal(v, res) {
		t.Errorf("marshal/unmarshal error %x become %x. data: %d [%x]", v, res, buf.Len(), buf.Bytes())
	}

	for i := 0; i < t.N; i++ {
		_ = w.Encode(v)
		_ = r.Decode(&res)
		buf.Reset()
	}
}

type TestableStruct1 struct {
	N int
	n int
	S []byte
}

func TestStruct1(t *testing.T) {
	v := TestableStruct1{
		N: 10000,
		S: []byte("qqwe;qcqlnoibiobowibeql"),
	}

	p := CreatePacker(v)

	size := p.Size(reflect.ValueOf(v))

	data, err := Marshal(v)
	if err != nil {
		t.Errorf("marshal error for value %v: %v", v, err)
	}

	if size != len(data) {
		t.Errorf("Size (%v) != len(data) (%v). data: [%x]", size, len(data), data)
	}

	var res TestableStruct1
	err = Unmarshal(data, &res)
	if err != nil {
		t.Errorf("unmarshal error for value %v: %v. data: %d [%x]", v, err, len(data), data)
	}

	if !reflect.DeepEqual(v, res) {
		t.Errorf("marshal/unmarshal error %v become %v. data: %d [%x]", v, res, len(data), data)
	}
}

func BenchmarkPackerStruct1Encode(t *testing.B) {
	var buf bytes.Buffer
	w := NewEncoder(&buf)
	r := NewDecoder(&buf)

	v := TestableStruct1{
		N: 10000,
		S: []byte("qqwe;qcqlnoibiobowibeql"),
	}
	var res TestableStruct1

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

	if !reflect.DeepEqual(v, res) {
		t.Errorf("marshal/unmarshal error %v become %v. data: %d [%x]", v, res, buf.Len(), buf.Bytes())
	}

	for i := 0; i < t.N; i++ {
		_ = w.Encode(v)
		_ = r.Decode(&res)
		buf.Reset()
	}
}

func BenchmarkGob___Struct1Encode(t *testing.B) {
	var buf bytes.Buffer
	w := gob.NewEncoder(&buf)
	r := gob.NewDecoder(&buf)

	v := TestableStruct1{
		N: 10000,
		S: []byte("qqwe;qcqlnoibiobowibeql"),
	}
	var res TestableStruct1

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

	if !reflect.DeepEqual(v, res) {
		t.Errorf("marshal/unmarshal error %v become %v. data: %d [%x]", v, res, buf.Len(), buf.Bytes())
	}

	for i := 0; i < t.N; i++ {
		_ = w.Encode(v)
		_ = r.Decode(&res)
		buf.Reset()
	}
}
