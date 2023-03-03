package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Type byte

// Various RESP kinds
const (
	SimpleString = '+'
	BulkString   = '$'
	Array        = '*'
)

type RESP struct {
	typ   Type
	bytes []byte
	array []*RESP
}

func (resp RESP) String() string {
	if resp.typ == SimpleString || resp.typ == BulkString {
		return string(resp.bytes)
	}
	return ""
}

func (resp RESP) Array() []*RESP {
	if resp.typ == Array {
		return resp.array
	}
	return nil
}

func DecodeRESP(byteStream *bufio.Reader) (*RESP, error) {
	dataTypeByte, err := byteStream.ReadByte()
	if err != nil {
		return nil, err
	}

	switch rune(dataTypeByte) {
	case SimpleString:
		return decodeRESPSimpleString(byteStream)
	case BulkString:
		return decodeRESPBulkString(byteStream)
	case Array:
		return decodeRESPArray(byteStream)
	}

	return nil, fmt.Errorf("unsupported data type: %s", string(dataTypeByte))
}

func decodeRESPSimpleString(byteStream *bufio.Reader) (*RESP, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return nil, fmt.Errorf("failed to read until CRLF: %s", err.Error())
	}

	return &RESP{
		typ:   SimpleString,
		bytes: readBytes,
	}, nil
}

func decodeRESPBulkString(byteStream *bufio.Reader) (*RESP, error) {
	stringCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return nil, fmt.Errorf("failed to read until CRLF: %s", err.Error())
	}

	count, err := strconv.Atoi(string(stringCount))
	if err != nil {
		return nil, fmt.Errorf("failed to convert '%s' to int: %s", string(stringCount), err.Error())
	}

	readBytes := make([]byte, count+2)

	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return nil, fmt.Errorf("failed to read bulk string contents: %s", err.Error())
	}

	return &RESP{
		typ:   BulkString,
		bytes: readBytes[:count],
	}, nil
}

func decodeRESPArray(byteStream *bufio.Reader) (*RESP, error) {
	arrayCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return nil, fmt.Errorf("failed to read until CRLF: %s", err.Error())
	}

	count, err := strconv.Atoi(string(arrayCount))
	if err != nil {
		return nil, fmt.Errorf("failed to convert '%s' to int: %s", string(arrayCount), err.Error())
	}

	array := []*RESP{}

	for i := 0; i < count; i++ {
		resp, err := DecodeRESP(byteStream)
		if err != nil {
			return nil, err
		}

		array = append(array, resp)
	}

	return &RESP{
		typ:   Array,
		array: array,
	}, nil
}

func readUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
	readBytes := []byte{}

	for {
		b, err := byteStream.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		readBytes = append(readBytes, b...)

		if len(readBytes) >= 2 && readBytes[len(readBytes)-2] == '\r' {
			break
		}
	}

	return readBytes[:len(readBytes)-2], nil
}
