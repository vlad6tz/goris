package protocol

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type Command string

const (
	GET Command = "GET"
	SET Command = "SET"
)

var (
	ErrUnknownCommand = errors.New("unknown command")
	ErrTooLong        = errors.New("line too long")
	ErrArity          = errors.New("wrong number of arguments")
)

const maxLine = 64 << 10

type Request struct {
	Command Command
	Key     string
	Value   []byte
}

func ReadRequest(r *bufio.Reader) (Request, error) {
	line, err := r.ReadString('\n')
	if errors.Is(err, io.EOF) && line == "" {
		return Request{}, io.EOF // clean disconnect, nothing pending
	}
	if err != nil {
		return Request{}, err
	}
	if len(line) > maxLine {
		return Request{}, ErrTooLong
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return Request{}, ErrArity
	}

	fields := strings.SplitN(line, " ", 3) // splitn preserves the space so edge cases like key = hello world dont get stored as "hello"
	switch Command(strings.ToUpper(fields[0])) {
	case GET:
		if len(fields) != 2 { // for get requests the len should be 2 because the only param accepted is the key
			return Request{}, ErrArity
		}
		return Request{
			Command: GET,
			Key:     fields[1],
		}, nil
	case SET:
		if len(fields) != 3 { // same logic as above, setting a key needs the meth + key + value, anything else is invalid
			return Request{}, ErrArity
		}
		return Request{
			Command: SET,
			Key:     fields[1],
			Value:   []byte(fields[2]),
		}, nil
	default:
		return Request{}, ErrUnknownCommand
	}
}
