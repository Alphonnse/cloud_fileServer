package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	FilesInDir   []string
	NextHop      []string
	ViewPath     []string
	DownloadPath []string
	UploadPath   string
	DeletePath   string
	CreatePath   string
}

func getFilesInDir(w http.ResponseWriter, r *http.Request, user database.User) []string{
	filesInDir, err := entriesFromDir(r.URL.Path[len(fmt.Sprintf("/"+user.Name+"/disk/")):])
	if err != nil {
		log.Println("there is error while reading the dir:", r.URL.Path[len(fmt.Sprintf("/"+user.Name+"/disk/")):])
		RespondWithError(w, 401, "Error while reading directory")
		return nil
	}
	return filesInDir
}


func FS(w http.ResponseWriter, r *http.Request, user database.User) {
	// There might be a trouble with cookie, when using links into site
	if strings.Split(r.URL.Path, "/")[1] != user.Name {
		RespondWithError(w, 401, "Bad url")
		return
	}

	page := &Page{}

	if r.URL.Query().Get("action") == "view" {
		// view handler
		path := r.URL.Path[len(fmt.Sprintf("/"+user.Name+"/disk/")):]

		// Check if the file exists
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			RespondWithError(w, 401, "File not found")
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Serve the file directly from disk
		http.ServeFile(w, r, path)

	} else if r.URL.Query().Get("action") == "download" {
		w.Header().Set("Content-Type", "application/json")

		// get the file name to download from url
		name := r.URL.Query().Get("name") // Error is here whern i am finding out the name of download file
		fmt.Println(name)

		// join to get the full file path
		directory := filepath.Join("files", name)

		// open file (check if exists)
		_, err := os.Open(directory)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "trouble with file:")
			return
		}

		// force a download with the content- disposition field
		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(directory))

		// serve file out.
		http.ServeFile(w, r, directory)
	} else if r.URL.Query().Get("action") == "upload" && r.Method == "POST" {

		reader, err := r.MultipartReader()
		if err != nil {
			log.Println("error while multiparting")
			return
		}

		for {
			part, err := reader.NextPart()
			if err != nil {
				break // No more parts to read
			}
			defer part.Close()

			if part.FormName() == "file" {

				// err := os.MkdirAll("./files", os.ModePerm)
				// if err != nil {
				// 	log.Println("error while reading dir")
				// }

				// dst, err := os.Create(strings.Replace("./"+r.URL.Path[len("arsen5/disk/"):], substring, "", -1) + "/" + part.FileName()) // Change this

				page.FilesInDir = getFilesInDir(w, r, user)

				for _, file := range page.FilesInDir {
					if part.FileName() == file {
						RespondWithError(w, 401, "file with this name already exists")
						return
					}
				}

				dst, err := os.Create("./"+r.URL.Path[len("arsen5/disk/"):] + "/" + part.FileName()) // Change this
				if err != nil {
					log.Println("err while creating the file", err)
					RespondWithError(w, 401, "error while creating the file")
					return
				}
				defer dst.Close()

				if _, err := io.Copy(dst, part); err != nil {
					log.Println("error while copying the file")
					RespondWithError(w, 401, "error while creating the file")
					return
				}

			} else {
				log.Println("no needed part")
				RespondWithError(w, 401, "error while creating the file")
				return
			}
		}
		http.Redirect(w, r, r.URL.Path, 301)

	} else if r.URL.Query().Get("action") == "delete" && r.Method == "POST" {
		type deleteFiles struct {
			Files []string
		}
		files := deleteFiles{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&files)
		if err != nil {
			RespondWithError(w, 500, fmt.Sprintf("error while decoding JSON"))
		}

		rightURL := r.URL.Path[len("/arsen5/disk/"):] // change

		for _, item := range files.Files {
			err := os.RemoveAll(filepath.Join(".", rightURL, item))
			if err != nil {
				RespondWithError(w, 500, fmt.Sprintf("error while removing file"))
			}
		}

	} else if r.URL.Query().Get("action") == "create" && r.Method == "POST"{
		rightURL := r.URL.Path[len("/arsen5/disk/"):] // change

		type CreateDir struct {
			DirectoryName string 
		}
		var dirName CreateDir
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&dirName)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error while decoding JSON")
			fmt.Println(err)
			return
		}


		page.FilesInDir = getFilesInDir(w, r, user)
		for _, file := range page.FilesInDir {
			if !strings.HasSuffix(dirName.DirectoryName, "/") {
				dirName.DirectoryName += "/"
			}

			if file == dirName.DirectoryName {
				RespondWithError(w, http.StatusBadRequest, "directory with that name already exists")
				return
			}
		}

		err = os.Mkdir((filepath.Join(".",rightURL,dirName.DirectoryName)), os.FileMode(0777))
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "error while creating the dir")
			return
		}
	} else {
		// Hop handler
		page.FilesInDir = getFilesInDir(w, r, user)

		// Upload link
		if strings.HasSuffix(r.URL.Path, "/") {
			page.UploadPath = fmt.Sprintf(r.URL.Path[:len(r.URL.Path)-len("/")] + "?action=upload")
			page.DeletePath = fmt.Sprintf(r.URL.Path[:len(r.URL.Path)-len("/")] + "?action=delete")
			page.CreatePath = fmt.Sprintf(r.URL.Path[:len(r.URL.Path)-len("/")] + "?action=create")
		} else {
			page.UploadPath = fmt.Sprintf(r.URL.Path + "?action=upload")
			page.DeletePath = fmt.Sprintf(r.URL.Path + "?action=delete")
			page.CreatePath = fmt.Sprintf(r.URL.Path + "?action=create")
		}

		for i := 0; i < len(page.FilesInDir); i++ {
			// If there is not slash in the end of our URL than adding it
			// to the beginning of the next hop
			if r.URL.Path[len(r.URL.Path)-1:] == "/" {
				page.NextHop = append(page.NextHop, r.URL.Path+page.FilesInDir[i])
			} else {
				page.NextHop = append(page.NextHop, r.URL.Path+"/"+page.FilesInDir[i])
			}

			// Creating the paths
			if isDirectory(page.NextHop[i]) == false {
				page.DownloadPath = append(page.DownloadPath, fmt.Sprintf(page.NextHop[i]+"?action=download"))
				page.ViewPath = append(page.ViewPath, fmt.Sprintf(page.NextHop[i]+"?action=view"))
			} else {
				page.ViewPath = append(page.ViewPath, "nil")
				page.DownloadPath = append(page.DownloadPath, "nil")
			}
		}

		t, err := template.New("list_of_files.html").Funcs(template.FuncMap{"isDirectory": isDirectory}).ParseFiles("tmpl/list_files/list_of_files.html")
		if err != nil {
			log.Println("Error while parsing the template of list", err)
			RespondWithError(w, 401, "Error while reading directory")
			return
		}
		t.Execute(w, page)
		fmt.Println(page.DownloadPath)
	}

}
