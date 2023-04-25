package cool

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func StartHttpServer(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open(os.Args[0])
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\"cool-transmission\"")
		io.Copy(w, file)
	})
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
