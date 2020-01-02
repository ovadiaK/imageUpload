package router

import (
	"fmt"
	"github.com/ovadiaK/imageUpload/img"
	"github.com/ovadiaK/imageUpload/router/handler"
	"net/http"
)

func Router() {
	http.Handle("/perm-images/", http.StripPrefix("/perm-images/", http.FileServer(http.Dir(img.IMAGE_FOLDER_PERM))))
	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/upload", handler.UploadFileHandler)
	http.HandleFunc("/select", handler.Selector)
	http.HandleFunc("/manipulate/", handler.Manipulator)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
