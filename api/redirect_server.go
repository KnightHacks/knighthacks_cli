package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func RunRedirectServer(ctx context.Context) (code string) {
	ctx, cancel := context.WithCancel(ctx)
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
		cancel()
	})
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("Error occured while shutting down: %s\n", err)
		}
	}()
	<-ctx.Done()

	return code
}
