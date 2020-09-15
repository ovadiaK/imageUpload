package main

import (
	"fmt"
	"github.com/ovadiaK/imageUpload/img"
	"github.com/ovadiaK/imageUpload/router"
	"os"
)

func main() {
	initFolders()
	fmt.Println("running on :8080")
	router.Router()
}
func initFolders() {
	if _, err := os.Stat(img.IMAGE_FOLDER_TEMP); err != nil {
		if err := os.Mkdir(img.IMAGE_FOLDER_TEMP, 0777); err != nil {
			panic(fmt.Errorf("create temporary image folder at %s failed: %s", img.IMAGE_FOLDER_TEMP, err))
		}
	}
	if _, err := os.Stat(img.IMAGE_FOLDER_PERM); err != nil {
		if err := os.Mkdir(img.IMAGE_FOLDER_PERM, 0777); err != nil {
			panic(fmt.Errorf("create permanent image folder at %s failed: %s", img.IMAGE_FOLDER_PERM, err))
		}
	}
}
