package main

import (
	"context"
	"fmt"
	_ "fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"io"
	"net/http"
	_ "net/http"

	"os"
	"path/filepath"
	"time"
)

var accountName = "your bucket name"
var accountKey = "your-bucket-access-key"

func main() {
	http.HandleFunc("/upload", Upload)
	http.ListenAndServe("localhost:8082", nil)

}

func getClient() (*azblob.Client, error) {
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	return azblob.NewClientWithSharedKeyCredential(url, cred, nil)
}

func Upload(res http.ResponseWriter, req *http.Request) {
	//fileName := "image"
	containerName := "your container name"
	client, err := getClient()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileMultiPart, fileHeader, err := req.FormFile("image")
	err = req.ParseMultipartForm(10 << 20)
	if err != nil {
		//log.Error(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	imgName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))

	file, err := os.Create(imgName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()
	defer os.Remove(imgName)
	_, err = io.Copy(file, fileMultiPart)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = client.UploadFile(context.Background(), containerName, imgName, file, nil)
	if err != nil {

		res.WriteHeader(http.StatusBadRequest)
	}
	return

}
