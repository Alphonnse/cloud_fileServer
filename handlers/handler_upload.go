package handlers

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func UploadGetHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tmpl/choose_and_upload/choose_upload.html")
		if err != nil {
			log.Println("Error while parsing the template of upload")
		  	return
		}
	t.Execute(w, nil)
}


func UploadPostHandler(w http.ResponseWriter, r *http.Request) {

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
			
			err := os.MkdirAll("./files", os.ModePerm)
			if err != nil {
				log.Println("error while reating dir")
			}

			dst, err := os.Create("./files/" + part.FileName())
			if err != nil {
				log.Println("err while creating the file")
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
}
