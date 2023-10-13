package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Alphonnse/file_server/internal/database"
)


func HandlerView (w http.ResponseWriter, r *http.Request, user database.User) {

	if strings.Split(r.URL.Path,"/")[1] != user.Name {
		RespondWithError(w, 401, "Bad url")
		return
	}

	if r.URL.Query().Get("action") == "view" {
		fmt.Println("thats the view handler")
	}



	path := r.URL.Path[len("/disk/view/"):]

	// Check if the file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		RespondWithError(w, 401, "File not found")
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Serve the file directly from disk
	http.ServeFile(w, r, path)
}
