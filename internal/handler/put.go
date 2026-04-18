package handler

import (
	"net"

	"github.com/codecrafters-io/http-server-starter-go/internal/config"
	"github.com/codecrafters-io/http-server-starter-go/internal/parser"
)

func HandlePut(request *parser.Request, conn net.Conn, config *config.Config) {
	res := &respone.Response{
		Header: make(map[string]string),
	}
	filename := request.Path[len("/files/"):]
	file_path := config.FilePath + "/" + filename

}
