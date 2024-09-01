package service

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveUploadedImage(fileHeaders []*multipart.FileHeader) []string {
	images := []string{}

	//create upload directory if it doesn't exist
	uploadsDirectory := "uploads"
	err := os.MkdirAll(uploadsDirectory, 0755)

	if err != nil {
		panic(err)
	}

	for _, fileHeader := range fileHeaders {
		// open the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// generate unique file name
		fileName := "image-" + fileHeader.Filename
		filePath := filepath.Join(uploadsDirectory, fileName)

		//create a new file to the uploads directory
		newFile, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		defer newFile.Close()

		// copy the uploaded file to the new file
		if _, err := io.Copy(newFile, file); err != nil {
			panic(err)
		}

		images = append(images, fileName)
	}

	return images
}
