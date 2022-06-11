package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func RunRedirectServer() (code string) {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.HandleFunc("/auth_redirect", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("uri=%s\n", request.RequestURI)
		code = request.FormValue("code")
		log.Printf("code=%s\n", code)
		writer.WriteHeader(200)
		writer.Write([]byte(fmt.Sprintf("you done did it code=%s", code)))
		go func() {
			time.Sleep(time.Second * 5)
			_ = server.Shutdown(context.Background())
		}()
	})

	_ = server.ListenAndServe()
	return code
}
