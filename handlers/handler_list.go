package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)


func entriesFromDir (dirName string) ([]string, error) {
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
	Files []string
	UrlToFile []string
	FullPathToDownload []string
}

func ListFiles(w http.ResponseWriter, r *http.Request) {
	
	page := &Page{}
	var err error	

	page.Files, err = entriesFromDir(r.URL.Path[len("/files"):])
	if err != nil {
		log.Println("there is error while reading the dir:", r.URL.Path[len("/files"):])
		respondWithError(w, 401, "Error while reading directory")
		return
	}
	
	// page.UrlToFile = page.Files

	for i := 0; i < len(page.Files); i ++ {
		if r.URL.Path[len(r.URL.Path)-1:] == "/" {
			page.UrlToFile = append(page.UrlToFile, r.URL.Path + page.Files[i])
		} else {
			page.UrlToFile = append(page.UrlToFile, r.URL.Path + "/" + page.Files[i])
		}

		page.FullPathToDownload = append(page.FullPathToDownload, page.UrlToFile[i][len("/disk"):])
	}


	// t, err := template.ParseFiles("tmpl/list_files/list_of_files.html")
	t, err := template.New("list_of_files.html").Funcs(template.FuncMap{"isDirectory": isDirectory}).ParseFiles("tmpl/list_files/list_of_files.html")
	if err != nil {
		log.Println("Error while parsing the template of list", err)
		respondWithError(w, 401, "Error while reading directory")
		return
	}
	t.Execute(w, page)
}

// func ListFiles(w http.ResponseWriter, r *http.Request) {
// 	page := &Page{}
//
// 	var err error
// 	page.Files, err = entriesFromDir(r.URL.Path[len("/files"):])
// 	page.UrlToFile = page.Files
// 	if err != nil {
// 		log.Println("there is an error while reading the dir:", r.URL.Path[len("/files"):])
// 		respondWithError(w, 401, "Error while reading directory")
// 		return
// 	}
//
// 	t, err := template.ParseFiles("tmpl/list_files/list_of_files.html")
// 	if err != nil {
// 		log.Println("Error while parsing the template of list", err)
// 		respondWithError(w, 401, "Error while reading directory")
// 		return
// 	}
//
// 	err = t.Execute(w, page)
// 	if err != nil {
// 		log.Println("Error while executing the template:", err)
// 		respondWithError(w, 500, "Internal Server Error")
// 	}
// }
