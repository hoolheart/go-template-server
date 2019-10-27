package main

import (
	"fmt"
	"net/http"
	"git.dev.tencent.com/petit_kayak/go-template-server/services"
)

// salute 处理问好请求
func salute(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	names := request.Form["name"]
	if len(names)>0 {
		fmt.Fprintf(writer, "Hello %s", names[0])
	} else {
		fmt.Fprintf(writer, "Hello world")
	}
}

// main 主程序
func main() {
	fmt.Println("Hello web world.");

	// prepare mux
	mux := http.NewServeMux()
	mux.HandleFunc("/salute/",salute)

	//setup services
	services.Setup(mux)

	// create and start server
	server := &http.Server {
		Addr: "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
