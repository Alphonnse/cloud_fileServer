package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Alphonnse/file_server/internal/database"
)

func entriesFromDir(dirName string) ([]string, error) {
	entries, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	list := make([]string, len(entries))

	i := 0
	for _, file := range entries {
		if file.IsDir() {
			list[i] = file.Name() + "/"
		} else {
			list[i] = file.Name()
		}
		i += 1
	}
	return list[:i], err
}

func isDirectory(entry string) bool {
	return strings.HasSuffix(entry, "/")
}

type Page struct {
	FilesInDir         []string
	UrlToFile          []string
	// FullPathToDownload []string
	ViewPath           []string
}

func FS(w http.ResponseWriter, r *http.Request, user database.User) {
	// There might be a trouble with cookie, when using links into site
	if strings.Split(r.URL.Path, "/")[1] != user.Name {
		RespondWithError(w, 401, "Bad url")
		return
	}

	page := &Page{}
	var err error
	
	if r.URL.Query().Get("action") == "view" {
		path := r.URL.Path[len("/arsen2/disk/"):]

		// Check if the file exists
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			RespondWithError(w, 401, "File not found")
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Serve the file directly from disk
		http.ServeFile(w, r, path)
			
	} else {
		page.FilesInDir, err = entriesFromDir(r.URL.Path[len(fmt.Sprintf("/"+user.Name+"/disk/")):])

		if err != nil {
			log.Println("there is error while reading the dir:", r.URL.Path[len(fmt.Sprintf("/"+user.Name+"/disk/")):])
			RespondWithError(w, 401, "Error while reading directory")
			return
		}

		for i := 0; i < len(page.FilesInDir); i++ {
			// If there is not slash in the end of our URL than adding it 
			// to the beginning of the next hop
			if r.URL.Path[len(r.URL.Path)-1:] == "/" {
				page.UrlToFile = append(page.UrlToFile, r.URL.Path+page.FilesInDir[i])
			} else {
				page.UrlToFile = append(page.UrlToFile, r.URL.Path+"/"+page.FilesInDir[i])
			}

			// page.FullPathToDownload = append(page.FullPathToDownload, page.UrlToFile[i][len("/disk"):]) // тут нужно поправить но пока не знаю как

			// if isDirectory(page.UrlToFile[i]) == false {
			// 	// its not a dir
			// 	page.ViewPath = append(page.ViewPath, fmt.Sprintf(page.UrlToFile[i]+"?action=view"))
			// } else {
			// 	page.ViewPath = append(page.ViewPath, "none")
			// }

			if isDirectory(page.UrlToFile[i]) == false {
				page.UrlToFile[i] = fmt.Sprintf(page.UrlToFile[i]+"?action=view")
			}
		}

		t, err := template.New("list_of_files.html").Funcs(template.FuncMap{"isDirectory": isDirectory}).ParseFiles("tmpl/list_files/list_of_files.html")
		if err != nil {
			log.Println("Error while parsing the template of list", err)
			RespondWithError(w, 401, "Error while reading directory")
			return
		}
		t.Execute(w, page)
	}

}

// func ListFilesOld(w http.ResponseWriter, r *http.Request) {
//
// 	page := &Page{}
// 	var err error
//
// 	page.FilesInDir, err = entriesFromDir(r.URL.Path[len("/files"):])
//
// 	if err != nil {
// 		log.Println("there is error while reading the dir:", r.URL.Path[len("/files"):])
// 		RespondWithError(w, 401, "Error while reading directory")
// 		return
// 	}
//
// 	for i := 0; i < len(page.FilesInDir); i++ {
// 		if r.URL.Path[len(r.URL.Path)-1:] == "/" {
// 			page.UrlToFile = append(page.UrlToFile, r.URL.Path+page.FilesInDir[i])
// 		} else {
// 			page.UrlToFile = append(page.UrlToFile, r.URL.Path+"/"+page.FilesInDir[i])
// 		}
//
// 		page.FullPathToDownload = append(page.FullPathToDownload, page.UrlToFile[i][len("/disk"):])
// 	}
//
// 	t, err := template.New("list_of_files.html").Funcs(template.FuncMap{"isDirectory": isDirectory}).ParseFiles("tmpl/list_files/list_of_files.html")
// 	if err != nil {
// 		log.Println("Error while parsing the template of list", err)
// 		RespondWithError(w, 401, "Error while reading directory")
// 		return
// 	}
// 	t.Execute(w, page)
// }
