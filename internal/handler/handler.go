package handler

import (
	"fmt"
	"net"
	"bufio"
	"github.com/codecrafters-io/http-server-starter-go/internal/parser"
	"github.com/codecrafters-io/http-server-starter-go/internal/config"
)

func HandleConnection(conn net.Conn, config *config.Config) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		request, err := parser.ParseRequest(reader)
		if err != nil {
			fmt.Println("Error parsing request: ", err.Error())
			return
		}
		
		fmt.Println("Received request: ", request)
		switch request.Method {
			case "GET":
				HandleGet(request, conn, config)
			case "POST":
				HandlePost(request, conn, config)
			case "DELETE":
				HandleDelete(request, conn, config)
			case "HEAD":
				HandleHead(request, conn, config)		
			case "PUT":
				
		}
		if request.Headers["Connection"] == "close" {
			break
		}
	}
}
