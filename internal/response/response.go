package response

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"

	"github.com/codecrafters-io/http-server-starter-go/internal/parser"
	"github.com/codecrafters-io/http-server-starter-go/internal/config"
	"strings"
	"os"
)

var SupportedEncodings = map[string]bool{
	"gzip": true,
	"deflate": true,
}

type Response struct {
	Version string
	StatusCode int
	StatusText string
	Headers map[string]string
	Body []byte
}


func HandleRoot(response *Response) {
	response.Version = "HTTP/1.1"
	response.StatusCode = 200
	response.StatusText = "OK"
	response.Headers["Content-Type"] = "text/plain"
}

func HandleEcho(response *Response, request *parser.Request) {
	response.Version = "HTTP/1.1"
	words := request.Path[len("/echo/"):]
	encoding, ok := request.Headers["Accept-Encoding"]
	if ok {
		encodings := strings.Split(encoding, ",")
		encode := "invalid-encoding"
		for _, e := range encodings {
			if SupportedEncodings[strings.TrimSpace(e)] {
				encode = strings.TrimSpace(e)
				break
			}
		}
		switch encode {
			case "gzip":
				var buf bytes.Buffer
				zw := gzip.NewWriter(&buf)
				if _, err := zw.Write([]byte(words)); err != nil {
					response.StatusCode = 500
					response.StatusText = "Internal Server Error"
					return
				}
				zw.Close()
				compressed := buf.Bytes()
				response.StatusCode = 200
				response.StatusText = "OK"
				response.Headers["Content-Type"] = "Text/plain"
				response.Headers["Content-Encoding"] = "gzip"
				response.Headers["Content-Length"] = fmt.Sprintf("%d", len(compressed))
				response.Body = compressed
			default:
				response.StatusCode = 200
				response.StatusText = "OK"
				response.Headers["Content-Type"] = "text/plain"
				response.Headers["Content-Length"] = fmt.Sprintf("%d", len(words))
				response.Body = []byte(words)
		}
	}else{
		response.StatusCode = 200
		response.StatusText = "OK"
		response.Headers["Content-Type"] = "text/plain"
		response.Headers["Content-Length"] = fmt.Sprintf("%d", len(words))
		response.Body = []byte(words)
	}
}

func HandleUserAgent(response *Response, request *parser.Request) {
	response.Version = "HTTP/1.1"
	response.StatusCode = 200
	response.StatusText = "OK"
	response.Headers["Content-Type"] = "text/plain"
	response.Headers["User-Agent"] = request.Headers["User-Agent"]
	response.Headers["Content-Length"] = fmt.Sprintf("%d", len(request.Headers["User-Agent"]))
	response.Body = []byte(request.Headers["User-Agent"])
}

func HandleFiles(response *Response, request *parser.Request, config *config.Config) {
	response.Version = "HTTP/1.1"
	file_name := request.Path[len("/files/"):]
	full_path := config.DirName + "/" + file_name
	data, err := os.ReadFile(full_path)
	if err != nil {
		response.StatusCode = 404
		response.StatusText = "Not Found"
		return
	}
	response.StatusCode = 200
	response.StatusText = "OK"
	response.Headers["Content-Type"] = "application/octet-stream"
	response.Headers["Content-Length"] = fmt.Sprintf("%d", len(data))
	response.Body = data
}

func HandleNotFound(response *Response) {
	response.Version = "HTTP/1.1"
	response.StatusCode = 404
	response.StatusText = "Not Found"
}

func HandleServerError(response *Response){
	response.Version = "HTTP/1.1"
	response.StatusCode = 500
	response.StatusText = "Internal Server Error"
}

func HandleFileCreate(response *Response){
	response.Version = "HTTP/1.1"
	response.StatusCode = 201
	response.StatusText = "Created"
}

func (r *Response) Write(conn net.Conn, keepAlive bool) {
	if keepAlive {
		r.Headers["Connection"] = "keep-alive"
	}else {
		r.Headers["Connection"] = "close"
	}
	status_line := fmt.Sprintf("%s %d %s\r\n", r.Version, r.StatusCode, r.StatusText)
	header_lines := fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d	\r\nConnection: %s\r\n", r.Headers["Content-Type"], len(r.Body), r.Headers["Connection"])
	switch r.StatusCode {
		case 200:
			if r.Body == nil && r.Headers == nil {
				conn.Write([]byte(status_line + "\r\n"))
				return
			} else if r.Headers["Content-Encoding"] != "" {
				header_lines = header_lines + fmt.Sprintf("Content-Encoding: %s\r\n", r.Headers["Content-Encoding"])
			}
			body_line := fmt.Sprintf("\r\n%s", r.Body)
			conn.Write([]byte(status_line + header_lines + body_line))
		case 404:
			conn.Write([]byte(status_line + "\r\n"))
		case 500:
			conn.Write([]byte(status_line + "\r\n"))
	}
}