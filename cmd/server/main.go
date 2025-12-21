package main

import (
	"cmp"
	"log"
	"net/http"
	"os"

	"github.com/typelate/muxt-template-module-htmx/internal/hypertext"
)

func main() {
	srv := &hypertext.Server{}
	mux := http.NewServeMux()
	hypertext.TemplateRoutes(mux, srv)
	log.Fatal(http.ListenAndServe(":"+cmp.Or(os.Getenv("PORT"), "8080"), mux))
}
