package packer

import (
	"bytes"
	"io"
	"reflect"
)

type Packer interface {
	Size(v interface{}) int
	WriteTo(w io.Writer, v interface{}) (int, error)
	ReadFrom(r io.Reader, v interface{}) (int, error)
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer

	e := NewEncoder(&buf)
	err := e.Encode(v)

	return buf.Bytes(), err
}

func Unmarshal(data []byte, v interface{}) error {
	r := bytes.NewReader(data)
	d := NewDecoder(r)

	return d.Decode(v)
}

type packerStorage struct {
	m map[reflect.Type]Packer
}

type Encoder struct {
	packerStorage
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	e := &Encoder{w: w}
	e.init()
	return e
}

func (e *Encoder) Encode(v interface{}) error {
	p := e.getPacker(v)
	_, err := p.WriteTo(e.w, v)
	return err
}

type Decoder struct {
	packerStorage
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{r: r}
	d.init()
	return d
}

func (d *Decoder) Decode(v interface{}) error {
	p := d.getPacker(v)
	_, err := p.ReadFrom(d.r, v)
	return err
}

func (s *packerStorage) init() {
	s.m = make(map[reflect.Type]Packer)
}

func (s *packerStorage) getPacker(v interface{}) Packer {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	p, ok := s.m[t]
	if !ok {
		p = makePacker(t, v)
		s.m[t] = p
	}

	if p == nil {
		panic("cant create packer")
	}

	return p
}

func makePacker(t reflect.Type, v interface{}) Packer {
	switch t.Kind() {
	case reflect.Int:
		return &Int64Packer{}
	}

	return nil
}

type Int64Packer struct {
	buf [10]byte
}

func (p *Int64Packer) Size(vi interface{}) int {
	v := uint64(toInt64(vi))
	return sizeOfVarint(v)
}

func (p *Int64Packer) WriteTo(w io.Writer, vi interface{}) (int, error) {
	v := uint64(toInt64(vi))
	offset := 0
	for v >= 1<<7 {
		p.buf[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	p.buf[offset] = uint8(v)
	offset++

	return w.Write(p.buf[:offset])
}

func (p *Int64Packer) ReadFrom(r io.Reader, vi interface{}) (int, error) {
	val := uint64(0)
	off := 0
	for shift := uint(0); ; shift += 7 {
		if shift >= 64 {
			panic("read error")
		}
		n, err := r.Read(p.buf[:1])
		if err != nil {
			return 0, err
		}
		if n != 1 {
			panic("n != 1")
		}
		b := p.buf[0]
		val |= (uint64(b) & 0x7F) << shift
		off++
		if b < 0x80 {
			break
		}
	}

	setInt64Ptr(vi, val)
	return off, nil
}

func toInt64(vi interface{}) int64 {
	var v int64
	switch t := vi.(type) {
	case int:
		v = int64(t)
	case int64:
		v = t
	default:
		panic("type assertion error")
	}
	return v
}

func setInt64Ptr(vi interface{}, val uint64) {
	switch t := vi.(type) {
	case *int:
		*t = int(val)
	case *int64:
		*t = int64(val)
	default:
		panic("type assertion error")
	}
}

func sizeOfVarint(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
