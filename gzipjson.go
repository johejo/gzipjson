package gzipjson

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
)

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Encode(v interface{}) error {
	zw := gzip.NewWriter(e.w)
	if err := json.NewEncoder(zw).Encode(&v); err != nil {
		return err
	}
	return zw.Close()
}

func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(data)).Decode(&v)
}

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Decode(v interface{}) error {
	gr, err := gzip.NewReader(d.r)
	if err != nil {
		return err
	}
	return json.NewDecoder(gr).Decode(&v)
}
