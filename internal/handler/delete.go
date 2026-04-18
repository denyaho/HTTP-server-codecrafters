package handler

import (
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/internal/config"
	"github.com/codecrafters-io/http-server-starter-go/internal/parser"
	"github.com/codecrafters-io/http-server-starter-go/internal/response"
)

func HandleDelete(request *parser.Request, conn net.Conn, config *config.Config) {
	res := &response.Response{
		Headers: make(map[string]string),
	}
	filename := request.Path[len("/files/"):]
	file_path := config.DirName + "/" + filename
	err := os.Remove(file_path)	
	if err != nil {
		response.HandleNotFound(res)
	} else {
		response.HandleFileDelete(res)
	}
	res.Write(conn, false, false)
}
