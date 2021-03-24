package main

import (
	"fmt"
	"log"
	"net/http"
	"git.dev.tencent.com/petit_kayak/go-template-server/services"
)

type writerPacker struct {
	writer http.ResponseWriter
	statusCode int
}

func (packer *writerPacker) Header() (http.Header) {
	return packer.writer.Header()
}

func (packer *writerPacker) Write(bytes []byte) (int, error) {
	return packer.writer.Write(bytes)
}

func (packer *writerPacker) WriteHeader(statusCode int) {
	packer.statusCode = statusCode
	packer.writer.WriteHeader(statusCode)
}

type logHandler struct {
	handler http.Handler
}

func (h *logHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	packer := writerPacker{writer,http.StatusOK}//pack response
	log.Printf("--> %s %s", request.Method, request.URL.RequestURI())//log request
	h.handler.ServeHTTP(&packer,request)// call inner handler
	log.Printf("<-- %d %s", packer.statusCode, http.StatusText(packer.statusCode))//log response
}

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
		Addr: "0.0.0.0:4000",
		Handler: &logHandler{mux},
	}
	server.ListenAndServe()
}
