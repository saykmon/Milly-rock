package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/files.html"))

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/files", filesHandler)
	http.Handle("/storage/", http.StripPrefix("/storage/", http.FileServer(http.Dir("storage"))))

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "index.html", nil)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		text := r.FormValue("text")
		if text == "" {
			http.Error(w, "Text cannot be empty", http.StatusBadRequest)
			return
		}

		filename := fmt.Sprintf("storage/file_%d.txt", len(getFileList())+1)
		err := ioutil.WriteFile(filename, []byte(text), 0644)
		if err != nil {
			http.Error(w, "Could not save file", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/files", http.StatusSeeOther)
	}
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		files := getFileList()
		templates.ExecuteTemplate(w, "files.html", files)
	}
}

func getFileList() []string {
	files := []string{}
	fileInfos, err := ioutil.ReadDir("storage")
	if err != nil {
		log.Println("Could not read storage directory:", err)
		return files
	}

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			files = append(files, fileInfo.Name())
		}
	}

	return files
}
