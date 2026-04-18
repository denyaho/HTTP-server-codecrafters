package handler

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/internal/parser"
	"github.com/codecrafters-io/http-server-starter-go/internal/config"
	"net"
	"strings"
	"os"
)

func HandlePost(request *parser.Request, conn net.Conn, config *config.Config) {
	directory := config.DirName
	if strings.HasPrefix(request.Path, "/files/"){
		file_name := request.Path[len("/files/"):]
		file_path := directory + "/" + file_name
		err := os.WriteFile(file_path, []byte(request.Body), 0644)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
			return
		}
		conn.Write([]byte(fmt.Sprintf("%s 201 Created\r\n\r\n", request.Version)))
	}
}
