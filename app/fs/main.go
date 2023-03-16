package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 8000, "serving port")
	root = flag.String("root", "/data2123", "fs root")
)

type handler struct {
	root string
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		return
	}
	p := request.URL.Path
	p = path.Join(h.root, p)
	p = path.Clean(p)
	if !strings.HasPrefix(p, h.root) {
		writer.Write([]byte(fmt.Sprintf("path is above allowed root")))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	info, err := os.Stat(p)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf("os.State(%v) failed: %v", p, err)))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		entities, err := os.ReadDir(p)
		if err != nil {
			writer.Write([]byte(fmt.Sprintf("os.ReadDir(%v) failed: %v", p, err)))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, e := range entities {
			writer.Write([]byte(fmt.Sprintf("%v\n", e.Name())))
		}
	} else {
		data, err := os.ReadFile(p)
		if err != nil {
			writer.Write([]byte(fmt.Sprintf("os.ReadFile(%v) failed: %v", p, err)))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Write(data)
	}
}

func main() {
	flag.Parse()

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", *port), &handler{root: *root}); err != nil {
		logrus.Fatal(err)
	}
}
