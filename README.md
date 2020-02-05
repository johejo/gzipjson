# gzipjson

[![go-test](https://github.com/johejo/gzipjson/workflows/go-test/badge.svg)](https://github.com/johejo/gzipjson/actions?query=workflow%3Ago-test)
[![GitHub license](https://img.shields.io/github/license/johejo/gzipjson)](https://github.com/johejo/gzipjson/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/johejo/gzipjson?status.svg)](https://godoc.org/github.com/johejo/gzipjson)

gzipped json Encoder and Decoder

## Install

```
go get github.com/johejo/gzipjson
```

## Example

```go
import "github.com/johejo/gzipjson"

// json.Marshal style
j := map[string]interface{}{"foo": "bar"}
b, _ := gzipjson.Marshal(j)

// json.Encoder style
j := map[string]interface{}{"foo": "bar"}
var w io.Writer
gzipjson.NewEncoder(w).Encode(j)

// json.Unmarshal style
b := []byte(`gzipped json byte slice`)
var j interface{}
gzipjson.Unmarshal(b, &j)

// json.Decoder style
f, _ := os.Open("test.json.gz")
var j interface{}
gzipjson.NewDecoder(f).Decode(&j)
```

## Licence

MIT

## Author

Mitsuo Heijo (@johejo)
