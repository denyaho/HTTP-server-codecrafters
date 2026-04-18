package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"strconv"
)

type Request struct {
	Method string
	Path string
	Version string
	Headers map[string]string
	Body []byte
}

func ParseRequest(reader *bufio.Reader) (*Request, error){
	req := Request{
		Headers: make(map[string]string),
	}

	request_line, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Error reading request: %v", err)
	}
	request_slice := strings.Fields(request_line)
	req.Method = request_slice[0]
	req.Path = request_slice[1]
	req.Version = request_slice[2]

	for {
		line, _ := reader.ReadString('\n')
		if line == "\r\n" {
			break
		}
		header := strings.Split(line, ": ")
		req.Headers[header[0]] = strings.TrimSpace(header[1])
	}
	content_length, _ := strconv.Atoi(req.Headers["Content-Length"])	
	buf := make([]byte, content_length)
	io.ReadFull(reader, buf)

	req.Body = buf
	return &req, nil
}

