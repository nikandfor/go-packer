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

	N2 *int
	//	S2 string

	Sub struct {
		V int
		S []byte
	}
}

func TestStruct1(t *testing.T) {
	v := TestableStruct1{
		N: 10000,
		S: []byte("qqwe;qcqlnoibiobowibeql"),
	}
	v.Sub.V = 400
	v.Sub.S = []byte("qwsss3wwqqqqf24330")

	p := CreatePacker(v)

	size := p.Size(reflect.ValueOf(&v))

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
	data := buf.Bytes()
	//	t.Logf("packed size: %v", s)

	t.ReportAllocs()
	t.SetBytes(int64(s))

	err = r.Decode(&res)
	if err != nil {
		t.Errorf("encode %v error: %v", v, err)
	}

	if !reflect.DeepEqual(v, res) {
		t.Errorf("marshal/unmarshal error %v become %v. data: %d [%x]", v, res, s, data)
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

type Context struct {
	buf bytes.Buffer
	e   enci
	d   deci
	val interface{}
	res interface{}
}

type enci interface {
	Encode(v interface{}) error
}
type deci interface {
	Decode(v interface{}) error
}

func runGlobal(t testing.TB) {
	ctx := Context{}
	for _, ed := range []struct {
		name string
		e    enci
		d    deci
	}{
		{
			name: "packer",
			e:    NewEncoder(&ctx.buf),
			d:    NewDecoder(&ctx.buf),
		},
		{
			name: "gob",
			e:    gob.NewEncoder(&ctx.buf),
			d:    gob.NewDecoder(&ctx.buf),
		},
	} {
		switch t := t.(type) {
		case *testing.T:
			t.Run(ed.name, func(t *testing.T) {
				ctx.e = ed.e
				ctx.d = ed.d
				runForEncoderDecoder(t, &ctx)
			})
		case *testing.B:
			t.Run(ed.name, func(t *testing.B) {
				ctx.e = ed.e
				ctx.d = ed.d
				runForEncoderDecoder(t, &ctx)
			})
		}
	}
}

func runForEncoderDecoder(t testing.TB, ctx *Context) {
	for _, suit := range []struct {
		name string
		N    int
		val  func() interface{}
		res  interface{}
	}{
		{
			name: "int",
			N:    1,
			val:  func() interface{} { return rand.Int() },
			res:  new(int),
		},
	} {
		switch t := t.(type) {
		case *testing.T:
			t.Run(suit.name, func(t *testing.T) {
				for i := 0; i < suit.N; i++ {
					ctx.val = suit.val()
					ctx.res = suit.res
					runForSuit(t, ctx)
				}
			})
		case *testing.B:
			t.Run(suit.name, func(t *testing.B) {
				ctx.val = suit.val()
				ctx.res = suit.res
				runForSuit(t, ctx)
			})
		}
	}
}

func runForSuit(t testing.TB, ctx *Context) {
	err := ctx.e.Encode(ctx.val)
	if err != nil {
		t.Errorf("encode %v error: %v", ctx.val, err)
	}

	s := ctx.buf.Len()
	//	t.Logf("packed size: %v", s)

	err = ctx.d.Decode(ctx.res)
	if err != nil {
		t.Errorf("encode %v error: %v", ctx.val, err)
	}

	res := reflect.ValueOf(ctx.res).Elem().Interface()
	if !reflect.DeepEqual(ctx.val, res) {
		t.Errorf("marshal/unmarshal error: %T %v become %T %v. data: %d [%x]", ctx.val, ctx.val, ctx.res, ctx.res, ctx.buf.Len(), ctx.buf.Bytes())
	}

	if t, ok := t.(*testing.B); ok {
		t.ReportAllocs()
		t.SetBytes(int64(s))

		for i := 0; i < t.N; i++ {
			_ = ctx.e.Encode(ctx.val)
			_ = ctx.d.Decode(ctx.res)
			ctx.buf.Reset()
		}
	}

	ctx.buf.Reset()
}

func TestGlobal(t *testing.T) {
	runGlobal(t)
}

func BenchmarkGlobal(t *testing.B) {
	runGlobal(t)
}

func TestDump(t *testing.T) {
}
