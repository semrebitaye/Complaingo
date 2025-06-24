package utility

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	appErrors "crud_api/internal/errors"
)

func SaveUploadFile(file multipart.File, handler multipart.FileHeader) (string, error) {
	// create directory for uploading file if not exist
	os.Mkdir("uploads", os.ModePerm)
	// generate unique filename with timestamp
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
	filePath := "./uploads/" + fileName

	// create file at the path to save the uploaded data
	dest, err := os.Create(filePath)
	if err != nil {
		return "", appErrors.ErrInvalidPayload.Wrap(err, "Failed to save uploaded file")
	}
	defer dest.Close()

	// copy all the uploded file to local file dest
	_, err = io.Copy(dest, file)
	if err != nil {
		return "", appErrors.ErrInvalidPayload.Wrap(err, "Failed to write file to disk")
	}

	return "/static/" + fileName, nil
}
