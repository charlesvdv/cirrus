package rest

import (
	"encoding/json"
	"io"
)

type marshaller interface {
	marshal(w io.Writer, v interface{}) error
	unmarshal(r io.Reader, v interface{}) error
	mediaType() string
}

var _ marshaller = jsonMarshaller{}

type jsonMarshaller struct {
}

func (m jsonMarshaller) marshal(w io.Writer, v interface{}) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(v)
}

func (m jsonMarshaller) unmarshal(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(v)
}

func (m jsonMarshaller) mediaType() string {
	return "application/json"
}
