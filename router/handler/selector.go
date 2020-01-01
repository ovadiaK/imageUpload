package handler

import (
	"fmt"
	"github.com/ovadiaK/imageUpload/img"
	"io/ioutil"
	"net/http"
)

type selectorMessage struct {
	Images []string
}

func Selector(w http.ResponseWriter, r *http.Request) {
	mess := selectorMessage{}
	fis, err := ioutil.ReadDir(img.IMAGE_FOLDER)
	if err != nil {
		fmt.Println(err)
	}
	mess.Images = make([]string, len(fis))
	for i, fi := range fis {
		mess.Images[i] = fi.Name()
	}
	renderTemplate(w, "select", mess)
}
