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
	"strconv"
	"sync"
	"time"
)

const (
	defaultSize = 1000
	maxSize     = 3000
	minSize     = 50
)

var allowedKinds = []string{"image/png", "image/jpeg", "image/tiff", "image/gif"}

type messageUpload struct {
	Success []string
	Failed  []string
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {

	//Parse our multipart form, 10 << 20 specifies a maximum
	//upload of 10 MB files.
	start := time.Now()
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	format := r.FormValue("format")
	size := sizeValue(r.FormValue("size"))
	makeRectangle := boolValue(r.FormValue("makeRectangle"))
	//get a ref to the parsed multipart form
	m := r.MultipartForm
	//get the *fileheaders
	files := m.File["myFile"]
	success := make(chan string, len(files))
	failed := make(chan string, len(files))

	var wg sync.WaitGroup
	for _, h := range files {
		header := h
		wg.Add(1)
		go upload(header, format, size, makeRectangle, success, failed, &wg)
	}
	wg.Wait()
	successStrings := make([]string, 0, len(success))
	resPos := len(success)
	for i := 0; i < resPos; i++ {
		successStrings = append(successStrings, <-success)
	}
	failedStrings := make([]string, 0, len(failed))
	resNeg := len(failed)
	for i := 0; i < resNeg; i++ {
		failedStrings = append(failedStrings, <-failed)
	}
	mess := messageUpload{Success: successStrings, Failed: failedStrings}
	fmt.Println(time.Since(start))
	renderTemplate(w, "upload", mess)
}

func sizeValue(str string) int {
	n, err := strconv.Atoi(str)
	if err != nil {
		return defaultSize
	}
	if n > maxSize {
		return maxSize
	}
	if n < minSize {
		return minSize
	}
	return n
}
func boolValue(str string) bool {
	fmt.Println("boolValue:", str)
	if str == "true" {
		return true
	}
	return false
}

func upload(header *multipart.FileHeader, format string, size int, makeRectangle bool, success, failed chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	if !imageCorrect(header) {
		fmt.Println("image not correct")
		fail(failed, header.Filename)
		return
	}
	//for each fileheader, get a handle to the actual file
	file, err := header.Open()
	if err != nil {
		fmt.Println("open header:", err)
		fail(failed, header.Filename)
		return
	}
	fileName := header.Filename
	tempPath := filepath.Join(img.IMAGE_FOLDER_TEMP, fileName)
	f, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file:", err)
		fail(failed, fileName, tempPath)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("copy file:", err)
		fail(failed, fileName, tempPath)
		return
	}
	if err := file.Close(); err != nil {
		fmt.Println("close file:", err)
		fail(failed, header.Filename, tempPath)
		return
	}
	if !correctFormat(fileName, format) {
		if fileName, err = img.Format(fileName, format); err != nil {
			fail(failed, fileName, filepath.Join(img.IMAGE_FOLDER_TEMP, fileName))
			return
		}
	}
	tempPath = filepath.Join(img.IMAGE_FOLDER_TEMP, fileName)
	sizedImage, err := img.Resize(fileName, size, makeRectangle)
	if err != nil {
		fmt.Println("resize:", err)
		fail(failed, fileName, tempPath)
		return
	}
	//create destination file making sure the path is writeable.
	//copy the uploaded file to the destination file
	err = imaging.Save(sizedImage, filepath.Join(img.IMAGE_FOLDER_PERM, fileName))
	if err != nil {
		fmt.Println("save image:", err)
		fail(failed, fileName, tempPath)
		return
	}
	if err := os.Remove(tempPath); err != nil {
		log.Println(err)
	}
	success <- fileName
}

func fail(failed chan string, filename string, toDelete ...string) {
	failed <- filename
	for _, dir := range toDelete {
		if err := os.Remove(dir); err != nil {
			log.Println(err)
		}
	}
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
func correctFormat(fileName, requestedFormat string) bool {
	currentFormat := filepath.Ext(fileName)
	fmt.Println(currentFormat, requestedFormat)
	if currentFormat[1:] == requestedFormat {
		fmt.Println(currentFormat[1:], requestedFormat)
		return true
	}
	return false
}
