package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)


func HandlerDownload (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get the file name to download from url
    name := r.URL.Query().Get("name")

    // join to get the full file path
    directory := filepath.Join("files", name)

    // open file (check if exists)
    _, err := os.Open(directory)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode("Unable to open file ")
        return
    }

    // force a download with the content- disposition field
    w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(directory))

    // serve file out.
    http.ServeFile(w, r, directory)
}
