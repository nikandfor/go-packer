package packer

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

type Packer interface {
	Size(v reflect.Value) int
	WriteTo(w io.Writer, v reflect.Value) (int, error)
	ReadFrom(r io.Reader, v reflect.Value) (int, error)
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
	_, err := p.WriteTo(e.w, reflect.ValueOf(v))
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
	_, err := p.ReadFrom(d.r, reflect.ValueOf(v))
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
		p = CreatePacker(v)
		s.m[t] = p
	}

	if p == nil {
		panic("cant create packer")
	}

	return p
}

func CreatePacker(v interface{}) Packer {
	t := reflect.TypeOf(v)
	return createPacker(t)
}

func createPacker(t0 reflect.Type) Packer {
	t := t0
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Int:
		return &Int64Packer{}
	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Uint8:
			return &BytesPacker{}
		default:
			panic(fmt.Errorf("can't create packer for slice of %v, kind: %v", t.Elem(), t.Elem().Kind()))
		}
	case reflect.Struct:
		return NewStructPacker(t)
	default:
		panic(fmt.Errorf("can't create packer for type %v, kind: %v\n%v", t, t.Kind(), DeepTypeDump(t0)))
	}
}

type Int64Packer struct {
	buf [10]byte
}

func (p *Int64Packer) Size(rv reflect.Value) int {
	rv = inderectValueConst(rv)
	switch rv.Type().Kind() {
	case reflect.Int:
	default:
		panic(fmt.Errorf("wrong type: %v.  kind %v", rv.Type(), rv.Kind()))
	}
	v := uint64(rv.Int())
	return sizeOfVarint(v)
}

func (p *Int64Packer) WriteTo(w io.Writer, rv reflect.Value) (int, error) {
	rv = inderectValueConst(rv)
	v := uint64(rv.Int())
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

func (p *Int64Packer) ReadFrom(r io.Reader, rv reflect.Value) (int, error) {
	rv = inderectValue(rv)
	val := int64(0)
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
		val |= (int64(b) & 0x7F) << shift
		off++
		if b < 0x80 {
			break
		}
	}

	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rv.SetInt(val)
	return off, nil
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

type BytesPacker struct {
	buf [128]byte
}

func (p *BytesPacker) Size(vi reflect.Value) int {
	l := vi.Len()
	return sizeOfVarint(uint64(l)) + l
}

func (p *BytesPacker) WriteTo(w io.Writer, rv reflect.Value) (int, error) {
	data := rv.Bytes()
	l := len(data)
	lsz := sizeOfVarint(uint64(l))

	v := l
	offset := 0
	for v >= 1<<7 {
		p.buf[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	p.buf[offset] = uint8(v)
	offset++

	if lsz+l < len(p.buf) {
		offset += copy(p.buf[offset:], data)
		return w.Write(p.buf[:offset])
	}

	n1, err := w.Write(p.buf[:offset])
	if err != nil {
		return n1, err
	}
	n2, err := w.Write(data)
	return n1 + n2, err
}

func (p *BytesPacker) ReadFrom(r io.Reader, rv reflect.Value) (int, error) {
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

	l := int(val)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	data := rv.Bytes()
	if len(data) < l {
		data = make([]byte, l)
		rv.SetBytes(data)
	}

	return r.Read(data[:l])
}

type StructPacker struct {
	typ    reflect.Type
	fields []structField
}

type structField struct {
	f reflect.StructField
	p Packer
}

func NewStructPacker(t reflect.Type) *StructPacker {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("type assertion error")
	}

	p := &StructPacker{typ: t}
	p.addFields(t)
	return p
}

func (s *StructPacker) addFields(t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		//	println(fmt.Sprintf("field: %v", ft))
		if ft.PkgPath != "" {
			continue
		}
		f := structField{f: ft, p: createPacker(ft.Type)}
		s.fields = append(s.fields, f)
	}
}

func (s *StructPacker) Size(rv reflect.Value) int {
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	size := 0

	for _, f := range s.fields {
		fv := rv.FieldByIndex(f.f.Index)
		size += f.p.Size(fv)
	}

	return size
}

func (s *StructPacker) WriteTo(w io.Writer, rv reflect.Value) (int, error) {
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	size := 0

	for _, f := range s.fields {
		fv := rv.FieldByIndex(f.f.Index)
		si, err := f.p.WriteTo(w, fv)
		size += si
		if err != nil {
			return size, err
		}
	}

	return size, nil
}

func (s *StructPacker) ReadFrom(r io.Reader, v reflect.Value) (int, error) {
	var rv reflect.Value
	if v.Type() == s.typ && v.CanSet() {
		rv = v
	} else {
		rv = inderectValue(v)
	}
	if rv.Type() != s.typ {
		panic("type assertion error")
	}

	size := 0

	for _, f := range s.fields {
		fv := rv.FieldByIndex(f.f.Index)
		si, err := f.p.ReadFrom(r, fv)
		size += si
		if err != nil {
			return size, err
		}
	}

	return size, nil
}

func inderectValue(rv0 reflect.Value) reflect.Value {
	rv := rv0
	//	println(fmt.Sprintf("indirect %v %v %v", rv.Kind(), rv.Type(), rv))
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
		//	println(fmt.Sprintf("    step %v %v %v", rv.Kind(), rv.Type(), rv))
	}
	return rv
}
func inderectValueConst(rv0 reflect.Value) reflect.Value {
	rv := rv0
	for rv.Kind() == reflect.Ptr {
		if !rv.IsNil() {
			rv = rv.Elem()
			continue
		}
		t := rv.Type()
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		return reflect.Zero(t)
	}
	return rv
}

func DeepValueDump(rv reflect.Value) string {
	var buf bytes.Buffer
	deepValueDump(rv, 0, &buf)
	return buf.String()
}
func deepValueDump(rv reflect.Value, sh int, w io.Writer) {
	for i := 0; i < sh; i++ {
		w.Write([]byte("    "))
	}
	fmt.Fprintf(w, "%v %v %v\n", rv.Type().Kind(), rv.Type(), rv)
}

func DeepTypeDump(t reflect.Type) string {
	var buf bytes.Buffer
	deepTypeDump(t, 0, &buf)
	return buf.String()
}
func deepTypeDump(t reflect.Type, sh int, w io.Writer) {
	for i := 0; i < sh; i++ {
		w.Write([]byte("    "))
	}
	fmt.Fprintf(w, "%v %v\n", t.Kind(), t)
}
