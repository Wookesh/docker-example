package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/wookesh/example/app/database"
)

var (
	port = flag.Int("port", 8000, "serving port")
	key  = flag.String("key", "default", "db key")
)

type handler struct {
	db  *database.Database
	key string
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		return
	}

	err := h.db.Incr(request.Context(), h.key)
	if err != nil {
		writer.Write([]byte(err.Error()))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write([]byte("OK"))
}

func main() {
	flag.Parse()

	db := database.New()

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", *port), &handler{db: db, key: *key}); err != nil {
		logrus.Fatal(err)
	}
}
