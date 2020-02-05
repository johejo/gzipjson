package gzipjson

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func assertJ(t *testing.T, a, b []byte) bool {
	t.Helper()
	var ja interface{}
	if err := json.Unmarshal(a, &ja); err != nil {
		t.Fatal(err)
	}
	var jb interface{}
	if err := json.Unmarshal(b, &jb); err != nil {
		t.Fatal(err)
	}
	return reflect.DeepEqual(ja, jb)
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want []byte
	}{
		{
			name: "map",
			v:    map[string]interface{}{"foo": "bar"},
			want: []byte(`{"foo": "bar"}`),
		},
		{
			name: "slice",
			v:    []interface{}{"foo", "bar"},
			want: []byte(`["foo", "bar"]`),
		},
		{
			name: "struct",
			v: struct {
				ID   int
				Name string
			}{ID: 99, Name: "gopher"},
			want: []byte(`{"ID": 99, "Name": "gopher"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.v)
			if err != nil {
				t.Fatal(err)
			}
			gr, err := gzip.NewReader(bytes.NewReader(got))
			if err != nil {
				t.Fatal(err)
			}
			gotb, err := ioutil.ReadAll(gr)
			if err != nil {
				t.Fatal(err)
			}
			if !assertJ(t, tt.want, gotb) {
				t.Fatalf("failed to compare json: want=%v, got=%v", tt.want, got)
			}
		})
	}
}

func TestEncoder(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want []byte
	}{
		{
			name: "map",
			v:    map[string]interface{}{"foo": "bar"},
			want: []byte(`{"foo": "bar"}`),
		},
		{
			name: "slice",
			v:    []interface{}{"foo", "bar"},
			want: []byte(`["foo", "bar"]`),
		},
		{
			name: "struct",
			v: struct {
				ID   int
				Name string
			}{ID: 99, Name: "gopher"},
			want: []byte(`{"ID": 99, "Name": "gopher"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			if err := NewEncoder(buf).Encode(tt.v); err != nil {
				t.Fatal(err)
			}
			gr, err := gzip.NewReader(buf)
			if err != nil {
				t.Fatal(err)
			}
			got, err := ioutil.ReadAll(gr)
			if err != nil {
				t.Fatal(err)
			}
			if !assertJ(t, tt.want, got) {
				t.Fatalf("failed to compare json: want=%v, got=%v", tt.want, got)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want interface{}
	}{
		{
			name: "map",
			b:    []byte(`{"foo": "bar"}`),
			want: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "slice",
			b:    []byte(`["foo", "bar"]`),
			want: []interface{}{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			gw := gzip.NewWriter(buf)
			if _, err := gw.Write(tt.b); err != nil {
				t.Fatal(err)
			}
			if err := gw.Close(); err != nil {
				t.Fatal(err)
			}
			var got interface{}
			if err := Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(tt.want, got) {
				t.Fatal(cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestDecoder(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want interface{}
	}{
		{
			name: "map",
			b:    []byte(`{"foo": "bar"}`),
			want: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "slice",
			b:    []byte(`["foo", "bar"]`),
			want: []interface{}{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			gw := gzip.NewWriter(buf)
			if _, err := io.Copy(gw, bytes.NewReader(tt.b)); err != nil {
				t.Fatal(err)
			}
			if err := gw.Close(); err != nil {
				t.Fatal(err)
			}
			var got interface{}
			if err := NewDecoder(buf).Decode(&got); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(tt.want, got) {
				t.Fatal(cmp.Diff(tt.want, got))
			}
		})
	}

	type testStruct struct {
		ID   int
		Name string
	}

	b := []byte(`{"ID": 99, "Name": "gopher"}`)
	want := testStruct{ID: 99, Name: "gopher"}

	buf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(buf)
	if _, err := io.Copy(gw, bytes.NewReader(b)); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	var got testStruct
	if err := NewDecoder(buf).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}
