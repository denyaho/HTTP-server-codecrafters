package handler

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/internal/parser"
	"github.com/codecrafters-io/http-server-starter-go/internal/config"
	"net"
	"strings"
	"os"
	"path/filepath"
	"bytes"
	"compress/gzip"
)

var SupportedEncoding = map[string]bool{
	"gzip": true,
	"deflate": true,
}

func HandleGet(request *parser.Request, conn net.Conn, config *config.Config) {
	if request.Path == "/" {
		conn.Write([]byte(fmt.Sprintf("%s 200 OK\r\n\r\n", request.Version)))
	}else if strings.HasPrefix(request.Path, "/echo/") {
		words := request.Path[len("/echo/"):]
		encoding, ok := request.Headers["Accept-Encoding"]
		if ok{
			encodings := strings.Split(encoding, ",")
			encode := "invalid-encoding"
			for _, e := range encodings {
				if SupportedEncoding[strings.TrimSpace(e)]{
					encode = strings.TrimSpace(e)
					break
				}
			}
			switch encode {
				case "gzip":
					var buf bytes.Buffer
					zw := gzip.NewWriter(&buf)
					if _, err := zw.Write([]byte(words)); err != nil {
						conn.Write([]byte(fmt.Sprintf("%s 500 Internal Server Error\r\n\r\n", request.Version)))
					}
					zw.Close()
					compressed := buf.Bytes()
					conn.Write([]byte(fmt.Sprintf("%s %d OK\r\nContent-Type: %s\r\nContent-Encoding: %s\r\nContent-Length: %d\r\n\r\n%s", request.Version, 200, "text/plain", "gzip", len(compressed), compressed)))
				default:
					conn.Write([]byte(fmt.Sprintf("%s %d OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", request.Version, 200, "text/plain", len(words), words)))
			}
		}else{
			conn.Write([]byte(fmt.Sprintf("%s %d OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", request.Version, 200, "text/plain", len(words), words)))
		}

	}else if request.Path == "/user-agent" {
		user_agent := request.Headers["User-Agent"]
		conn.Write([]byte(fmt.Sprintf("%s %d OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", request.Version, 200, "text/plain", len(user_agent), user_agent)))
	} else if strings.HasPrefix(request.Path, "/files"){
		file_path := request.Path[len("/files/"):]
		data, err := os.ReadFile(filepath.Join(config.DirName, file_path))
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("%s 404 Not Found\r\n\r\n", request.Version)))
			return
		}
		conn.Write([]byte(fmt.Sprintf("%s %d OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", request.Version, 200, "application/octet-stream", len(data), data)))
	}else{
		conn.Write([]byte(fmt.Sprintf("%s 404 Not Found\r\n\r\n", request.Version)))
	}
}