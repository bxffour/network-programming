package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	// Declaring types of with a size of 1 byte
	BinaryType uint8 = iota + 1
	StingType

	// The 4-byte integer used to designate the Maximum payload size has a
	// maximum value of 4,294,967,295 indicating a payload of over 4GB. It would
	// be easy for a malicious actor to perform a Denial-of-Service attack that
	// exhausts all available RAM on my computer. Keeping the maximum payload size
	// at a reasonable size makes memory exhaustion attacks harder to execute
	MaxPayloadSize uint32 = 10 << 20 // 10MB
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}

type Binary []byte

func (m Binary) Bytes() []byte { return m }

func (m Binary) String() string { return string(m) }

func (m Binary) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, BinaryType) // Writing 1-byte type
	if err != nil {
		return 0, err
	}

	n = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m))) // Writing 4-Byte size value
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write(m) // Writing the payload

	return n + int64(o), err
}

func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	// var typ uint8
	// err := binary.Read(r, binary.BigEndian, &typ) // reading 1-byte type
	// if err != nil {
	// 	return 0, err
	// }

	var n int64 = 1

	// Checking if the type read from the payload is a valid binary
	// if typ != BinaryType {
	// 	return n, errors.New("invalid binary")
	// }

	var size uint32
	err := binary.Read(r, binary.BigEndian, &size) // reading 4-byte size
	if err != nil {
		return n, err
	}

	n += 4

	// Checking if the size of the payload is within the 10MB limit
	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	*m = make([]byte, size) // Creating a byte slice the size of the payload
	o, err := r.Read(*m)    // Reading the actual payload

	return n + int64(o), err
}

type String string

func (m String) Bytes() []byte { return []byte(m) }

func (m String) String() string { return string(m) }

func (m String) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, StingType)
	if err != nil {
		return 0, err
	}

	n += 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m)))
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write([]byte(m))

	return n + int64(o), err
}

func (m *String) ReadFrom(r io.Reader) (n int64, err error) {
	// var typ uint8
	// err = binary.Read(r, binary.BigEndian, &typ)
	// if err != nil {
	// 	return 0, err
	// }

	n = 1

	// if typ != StingType {
	// 	return n, errors.New("invalid String")
	// }

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return n, err
	}

	n += 4

	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	buf := make([]byte, size)
	o, err := r.Read(buf)
	if err != nil {
		return n, err
	}

	*m = String(buf)

	return n + int64(o), nil
}

func decode(r io.Reader) (Payload, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, err
	}

	var payload Payload

	switch typ {
	case BinaryType:
		payload = new(Binary)
	case StingType:
		payload = new(String)
	default:
		return nil, errors.New("unknown Type")
	}

	// MultiReader is used to inject the byte we've already read i.e the type
	// back into the reader. This is to account for the fact thar the ReadFrom
	// reads the type field, thus it requires it to be present. An optimal solution would
	// be to eliminate the need for ReadFrom to read the type, thus eliminating
	// the need for io.MuliReader.

	// _, err = payload.ReadFrom(io.MultiReader(bytes.NewReader([]byte{typ}), r))

	// Since the need for ReaderFrom to read the type field has been eliminated, we
	// don't have to use io.MultiReader to inject the already read type field back in.
	_, err = payload.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
