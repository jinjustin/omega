package storage

import (

    "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/http"
	"fmt"
	"os"
	"io"
	"crypto/rand"
	"context"
	"encoding/json"
	"strings"
)

//StorageInfo is struct that use to store information of storage.
type StorageInfo struct{
	Endpoint string
	AccessKeyID string
	SecretAccessKey string
	UseSSL bool
	BucketName string
	URL string
}

//Storage is store storage information
var Storage = StorageInfo{
	Endpoint: "142.93.177.152:9000",
	AccessKeyID: "minioadmin",
	SecretAccessKey: "minioadmin",
	UseSSL: false,
	BucketName: "omega",
	URL: "http://142.93.177.152:9000/omega/",
}

func deleteFile(name string){
	err := os.Remove(name)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func generateFilename() string{
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}

//UploadPic is a function that use to upload picture to storage.
var UploadPic = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		w.Write([]byte("Unsuccessfully uploaded\n"))
		return
	}

	defer file.Close()

	filetype := handler.Header.Get("Content-type")
	filename := handler.Filename

	fileSplit := strings.Split(filename,".")

	var newFileName string

	if len(fileSplit) == 2{
		newFileName = generateFilename() + "." + fileSplit[1]
	}else{
		newFileName = filename
	}

	defer deleteFile(newFileName)
	defer file.Close()

	dst, err := os.Create(newFileName)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	// Initialize minio client object.
	minioClient, err := minio.New(Storage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(Storage.AccessKeyID, Storage.SecretAccessKey, ""),
		Secure: Storage.UseSSL,
	})
	if err != nil {
		panic(err)
	}

	filePath := "./" + newFileName

    // Upload the zip file with FPutObject
    _, err = minioClient.FPutObject(ctx, Storage.BucketName, newFileName, filePath, minio.PutObjectOptions{ContentType: filetype})
    if err != nil {
    	panic(err)
    }

	type ImageLink struct{
		URL string
	}

	var imageLink ImageLink

	imageLink.URL = Storage.URL + newFileName

	json.NewEncoder(w).Encode(imageLink)
})
