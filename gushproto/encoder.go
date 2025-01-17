package gushproto

import (
	"errors"
	"io"

	"google.golang.org/protobuf/proto"
)

func WriteProto(w io.Writer, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if err := encodeVarint(w, uint64(len(data))); err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func ReadProto(r io.Reader, msg proto.Message) error {
	length, err := decodeVarint(r)
	if err != nil {
		return err
	}

	if length == 0 {
		return errors.New("zero length message")
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return err
	}

	return proto.Unmarshal(buf, msg)
}

func encodeVarint(w io.Writer, x uint64) error {
	for x >= 1<<7 {
		b := byte(x&0x7f | 0x80)
		if _, err := w.Write([]byte{b}); err != nil {
			return err
		}

		x >>= 7
	}

	_, err := w.Write([]byte{byte(x)})
	return err
}

func decodeVarint(r io.Reader) (uint64, error) {
	var x uint64
	var s uint

	for i := 0; i < 10; i++ {
		var b [1]byte
		if _, err := r.Read(b[:]); err != nil {
			return 0, err
		}

		if b[0] < 0x80 {
			if i == 9 && b[0] > 1 {
				return 0, errors.New("varint too large")
			}
			return x | uint64(b[0])<<s, nil
		}
		x |= uint64(b[0]&0x7f) << s
		s += 7
	}
	return 0, errors.New("varint too large, more than 10 bytes")
}
