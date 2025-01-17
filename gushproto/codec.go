package gushproto

import (
	"io"

	"google.golang.org/protobuf/proto"
)

type bufCodec struct {
	w io.Writer
	r io.Reader
}

type Encoder interface {
	ProtoBufEncode
	ProtoBufDecode
}

type ProtoBufEncode interface {
	Encode(msg proto.Message) error
}

type ProtoBufDecode interface {
	Decode(msg proto.Message) error
}

func NewProtoRW(r io.Reader, w io.Writer) Encoder {
	return &bufCodec{w: w, r: r}
}

func (rw *bufCodec) Encode(msg proto.Message) error {
	return WriteProto(rw.w, msg)
}

func (rw *bufCodec) Decode(msg proto.Message) error {
	return ReadProto(rw.r, msg)
}
