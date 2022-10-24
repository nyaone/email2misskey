package jobs

import (
	"archive/zip"
	"bytes"
	"email2misskey/global"
	"email2misskey/types"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

func CompressAndUploadFile(filename string, fileBuffer *bytes.Buffer) (string, error) {
	// Zip file data
	zippedFileName, zippedBuffer, err := compressFile(filename, fileBuffer)
	if err != nil {
		global.Logger.Errorf("Failed to compress file with error: %v", err)
		return "", err
	}

	// Prepare
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	err = bodyWriter.WriteField("i", global.Config.Misskey.Token)
	if err != nil {
		global.Logger.Errorf("Failed to write field data i with error: %v", err)
		return "", err
	}
	if global.Config.Misskey.FolderID != "" {
		err = bodyWriter.WriteField("folderId", global.Config.Misskey.FolderID)
		if err != nil {
			global.Logger.Errorf("Failed to write field data folderId with error: %v", err)
			return "", err
		}
	}
	err = bodyWriter.WriteField("name", zippedFileName)
	if err != nil {
		global.Logger.Errorf("Failed to write field data name with error: %v", err)
		return "", err
	}

	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		global.Logger.Errorf("Failed to create file field with error: %v", err)
		return "", err
	}

	_, err = fileWriter.Write(zippedBuffer.Bytes())
	if err != nil {
		global.Logger.Errorf("Failed to copy buffer data with error: %v", err)
		return "", err
	}

	// Upload to Misskey
	formContentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close() // Ignore error
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/drive/files/create", global.Config.Misskey.Instance), bodyBuffer)
	if err != nil {
		global.Logger.Errorf("Failed to initialize request with error: %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", formContentType)
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		global.Logger.Errorf("Failed to finish request with error: %v", err)
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		// Failed
		global.Logger.Errorf("Upload failed.")
		var apiErrorRes types.ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&apiErrorRes)
		if err != nil {
			global.Logger.Errorf("Failed to decode response body with error: %v", err)
			return "", err
		}
		global.Logger.Errorf("Request failed with error: %v", apiErrorRes)
		return "", fmt.Errorf(apiErrorRes.Error.Message)
	}

	// Succeeded
	var dfcRes types.DriveFilesCreate_Response
	err = json.NewDecoder(res.Body).Decode(&dfcRes)
	if err != nil {
		global.Logger.Errorf("Failed to decode response body with error: %v", err)
		return "", err
	}

	return dfcRes.ID, nil

}

func compressFile(filename string, fileBuffer *bytes.Buffer) (string, *bytes.Buffer, error) {

	zippedFileName := filename + ".zip"
	zippedBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zippedBuffer)
	zipFile, err := zipWriter.Create(filename)
	if err != nil {
		global.Logger.Errorf("Failed to create zip file with error: %v", err)
		return "", nil, err
	}
	_, err = zipFile.Write(fileBuffer.Bytes())
	if err != nil {
		global.Logger.Errorf("Failed to write zip file with error: %v", err)
		return "", nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		global.Logger.Errorf("Failed to close zip writer with error: %v", err)
		return "", nil, err
	}

	return zippedFileName, zippedBuffer, nil
}
