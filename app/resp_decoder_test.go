package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestDecodeSimpleString(t *testing.T) {
	resp, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString("+foo\r\n")))

	if err != nil {
		t.Errorf("error decoding simple string: %s", err)
	}

	if resp.typ != SimpleString {
		t.Errorf("expected SimpleString, got %v", resp.typ)
	}

	if resp.String() != "foo" {
		t.Errorf("expected 'foo', got '%s'", resp.String())
	}
}

func TestDecodeBulkString(t *testing.T) {
	resp, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString("$4\r\nabcd\r\n")))

	if err != nil {
		t.Errorf("error decoding bulk string: %s", err)
	}

	if resp.typ != BulkString {
		t.Errorf("expected BulkString, got %v", resp.typ)
	}

	if resp.String() != "abcd" {
		t.Errorf("expected 'abcd', got '%s'", resp.String())
	}
}

func TestDecodeArray(t *testing.T) {
	resp, err := DecodeRESP(bufio.NewReader(bytes.NewBufferString("*2\r\n$3\r\nGET\r\n$4\r\nthis\r\n")))

	if err != nil {
		t.Errorf("error decoding array: %s", err)
	}

	if resp.typ != Array {
		t.Errorf("expected Array, got %v", resp.typ)
	}

	if len(resp.Array()) != 2 {
		t.Errorf("expected array of length 2, got %v", len(resp.Array()))
	}

	if resp.Array()[0].String() != "GET" {
		t.Errorf("expected 'GET', got '%s'", resp.Array()[0].String())
	}

	if resp.Array()[1].String() != "this" {
		t.Errorf("expected 'this', got '%s'", resp.Array()[1].String())
	}
}
