package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)


func HandlerDownload (w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/disk/download/"):]
	fmt.Println(path)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("there is error while reding the file", err)
		respondWithError(w, 400, "shit")
	}

	_, err = io.Copy(w, bytes.NewReader(file))
	if err != nil {
		log.Println("error while copying the file")
		respondWithError(w, 401, "error while creating the file")
		return
	}
}
