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
	"os"
	"path/filepath"
)

var allowedKinds = []string{"image/png", "image/jpeg"}

type messageUpload struct {
	Success []string
	Failed  []string
}

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
	success := make([]string, 0, len(files))
	failed := make([]string, 0, len(files))
	for i, header := range files {
		if !imageCorrect(header) {
			failed = append(failed, header.Filename)
			continue
		}
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		if err != nil {
			failed = append(failed, header.Filename)
			continue
		}
		fileName := header.Filename
		f, err := os.OpenFile(filepath.Join(img.IMAGE_FOLDER_TEMP, fileName), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			failed = append(failed, header.Filename)
			if err := os.Remove(filepath.Join(img.IMAGE_FOLDER_TEMP, fileName)); err != nil {
				log.Println(err)
			}
			continue
		}
		defer f.Close()
		n, err := io.Copy(f, file)
		if err != nil {
			failed = append(failed, header.Filename)
			if err := os.Remove(filepath.Join(img.IMAGE_FOLDER_TEMP, fileName)); err != nil {
				log.Println(err)
			}
			continue
		}
		log.Printf("filename: %v %v of %v bytes written/", fileName, n, header.Size)
		if err := file.Close(); err != nil {
			failed = append(failed, header.Filename)
			if err := os.Remove(filepath.Join(img.IMAGE_FOLDER_TEMP, fileName)); err != nil {
				log.Println(err)
			}
			continue
		}
		if header.Header.Get("Content-Type") != "image/jpeg" {
			if fileName, err = img.Format(fileName); err != nil {
				failed = append(failed, header.Filename)
				if err := os.Remove(filepath.Join(img.IMAGE_FOLDER_TEMP, fileName)); err != nil {
					log.Println(err)
				}
				if err := os.Remove(filepath.Join(img.IMAGE_FOLDER_TEMP, header.Filename)); err != nil {
					log.Println(err)
				}
				continue
			}
		}
		fmt.Println(fileName)
		sizedImage, err := img.Resize(fileName)
		if err != nil {
			failed = append(failed, header.Filename)
			if err := os.Remove(filepath.Join(img.IMAGE_FOLDER_TEMP, fileName)); err != nil {
				log.Println(err)
			}
			if err := os.Remove(filepath.Join(img.IMAGE_FOLDER_TEMP, header.Filename)); err != nil {
				log.Println(err)
			}
			continue
		}
		//create destination file making sure the path is writeable.
		//copy the uploaded file to the destination file
		err = imaging.Save(sizedImage, filepath.Join(img.IMAGE_FOLDER_PERM, fileName))
		if err != nil {
			failed = append(failed, header.Filename)
			if err := os.Remove(filepath.Join(img.IMAGE_FOLDER_TEMP, fileName)); err != nil {
				log.Println(err)
			}
			if err := os.Remove(filepath.Join(img.IMAGE_FOLDER_TEMP, header.Filename)); err != nil {
				log.Println(err)
			}
			continue
		}
		if err := os.Remove(filepath.Join(img.IMAGE_FOLDER_TEMP, header.Filename)); err != nil {
			log.Println(err)
		}
		success = append(success, header.Filename)
	}
	mess := messageUpload{Success: success, Failed: failed}

	renderTemplate(w, "upload", mess)
}
func imageCorrect(head *multipart.FileHeader) bool {
	f, err := head.Open()
	if err != nil {
		log.Println(err)
		return false
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		log.Println(err)
		return false
	}
	buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
	if _, err := buf.Read(buff); err != nil {
		log.Println(err)
		return false
	}
	format := http.DetectContentType(buff)
	if !stringInSlice(format, allowedKinds) {
		log.Println(format, fmt.Errorf("Got ya, ivan!\n"))
		return false
	}
	if head.Header.Get("Content-Type") != format {
		log.Printf("formats mismatch: %v %v", head.Header.Get("Content-Type"), format)
	}
	return true
}
