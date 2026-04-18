package handler

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/internal/parser"
	"github.com/codecrafters-io/http-server-starter-go/internal/config"
	"github.com/codecrafters-io/http-server-starter-go/internal/response"
	"net"
	"strings"
	"os"
)

func HandlePost(request *parser.Request, conn net.Conn, config *config.Config) {
	directory := config.DirName
	res := &response.Response{
		Headers: make(map[string]string),
	}
	if strings.HasPrefix(request.Path, "/files/"){
		file_name := request.Path[len("/files/"):]
		file_path := directory + "/" + file_name
		err := os.WriteFile(file_path, []byte(request.Body), 0644)
		if err != nil {
			response.HandleServerError(res)
		}else{
			response.HandleFileCreate(res)
		}
		res.Write(conn, false)
	}
}
