package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var folder string

func init (){
	flag.StringVar(&folder, "folder", "", "folder to store static files")
}

func main() {
	flag.Parse()
	if folder == "" {
		log.Fatal("Usage: httpServer -folder /web")
	}
	fs := http.FileServer(http.Dir(folder))
	http.Handle("/", fs)
	http.HandleFunc("/upload", upload)

	log.Println("Listening on :3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func upload(w http.ResponseWriter ,  r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>File Upload</title>
  </head>
  <body>
    <form
      enctype="multipart/form-data"
      action="/upload"
      method="post"
    >
      <input type="file" name="myFile" />
      <input type="submit" value="upload" />
    </form>
  </body>
</html>
`))
		return
	}
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		log.Fatal(err)
		return
	}
	defer file.Close()
	destFile, err := os.Create(fmt.Sprintf("%s/%s", folder, handler.Filename))
	if err != nil {
		log.Fatal(err)
	}
	defer destFile.Close()

	buffer := make([]byte, 4096)
	for {
		_, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			destFile.Write(buffer)
			break
		}
		destFile.Write(buffer)
	}
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
