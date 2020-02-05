package gzipjson

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
)

func Marshal(v interface{}) ([]byte, error) {
	pr, pw := io.Pipe()
	defer pr.Close()
	go func() {
		defer pw.Close()
		_ = NewEncoder(pw).Encode(v)
	}()
	return ioutil.ReadAll(pr)
}

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Encode(v interface{}) error {
	pr, pw := io.Pipe()
	defer pr.Close()
	go func() {
		defer pw.Close()
		_ = json.NewEncoder(pw).Encode(v)
	}()

	gw := gzip.NewWriter(e.w)
	defer gw.Close()
	if _, err := io.Copy(gw, pr); err != nil {
		return err
	}
	return nil
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
