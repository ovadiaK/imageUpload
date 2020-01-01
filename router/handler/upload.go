package handler

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/ovadiaK/imageUpload/img"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

var allowedKinds = []string{"image/png", "image/jpeg"}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {

	//Parse our multipart form, 10 << 20 specifies a maximum
	//upload of 10 MB files.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get a ref to the parsed multipart form
	m := r.MultipartForm

	//get the *fileheaders
	files := m.File["myFile"]
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		format, err := checkImage(&file)
		fmt.Println(format)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sizedImage, err := img.Resize(&file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//create destination file making sure the path is writeable.
		//copy the uploaded file to the destination file
		err = imaging.Save(sizedImage, filepath.Join(img.IMAGE_FOLDER, files[i].Filename))
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}
		if err := file.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/select", 202)
}
func checkImage(file *multipart.File) (string, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, *file); err != nil {
		fmt.Println(err)
	}
	buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
	if _, err := buf.Read(buff); err != nil {
		return "", err
	}
	kind := http.DetectContentType(buff)
	if !stringInSlice(kind, allowedKinds) {
		return kind, fmt.Errorf("Got ya, ivan!\n")
	}
	return kind, nil
}
