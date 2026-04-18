package handler

import (
	"github.com/codecrafters-io/http-server-starter-go/internal/parser"
	"github.com/codecrafters-io/http-server-starter-go/internal/config"
	"net"
	"os"
)

func HandleDelete(request *parser.Request, conn net.Conn, config *config.Config) {
	filename := request.Path[len("/files/"):]
	file_path := config.DirName + "/" + filename
	err := os.Remove(file_path)	
	if err != nil {
		
	}
}
