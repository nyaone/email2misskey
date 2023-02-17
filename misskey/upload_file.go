package misskey

import (
	"bytes"
	"email2misskey/config"
	"email2misskey/global"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

type DriveFilesCreate_Response struct {
	ID  string `json:"id"`
	URL string `yaml:"url"`
	// Ignore other fields
}

type MessageCreate_Request struct {
	I      string `json:"i"` // Token
	UserID string `json:"userId"`
	Text   string `json:"text"`
	FileID string `json:"fileId"`
}

func UploadFile(filename string, fileBuffer *bytes.Buffer) (string, string, error) {
	// Prepare
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	err := bodyWriter.WriteField("i", config.Config.Misskey.Token)
	if err != nil {
		global.Logger.Errorf("Failed to write field data i with error: %v", err)
		return "", "", err
	}
	if config.Config.Misskey.FolderID != "" {
		err = bodyWriter.WriteField("folderId", config.Config.Misskey.FolderID)
		if err != nil {
			global.Logger.Errorf("Failed to write field data folderId with error: %v", err)
			return "", "", err
		}
	}
	err = bodyWriter.WriteField("name", filename)
	if err != nil {
		global.Logger.Errorf("Failed to write field data name with error: %v", err)
		return "", "", err
	}

	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		global.Logger.Errorf("Failed to create file field with error: %v", err)
		return "", "", err
	}

	_, err = fileWriter.Write(fileBuffer.Bytes())
	if err != nil {
		global.Logger.Errorf("Failed to copy buffer data with error: %v", err)
		return "", "", err
	}

	// Upload to Misskey
	formContentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close() // Ignore error
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/drive/files/create", config.Config.Misskey.Instance), bodyBuffer)
	if err != nil {
		global.Logger.Errorf("Failed to initialize request with error: %v", err)
		return "", "", err
	}
	req.Header.Set("Content-Type", formContentType)
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		global.Logger.Errorf("Failed to finish request with error: %v", err)
		return "", "", err
	}

	if res.StatusCode != http.StatusOK {
		// Failed
		global.Logger.Errorf("Upload failed.")
		var apiErrorRes Error_Response
		err = json.NewDecoder(res.Body).Decode(&apiErrorRes)
		if err != nil {
			global.Logger.Errorf("Failed to decode response body with error: %v", err)
			return "", "", err
		}
		global.Logger.Errorf("Request failed with error: %v", apiErrorRes)
		return "", "", fmt.Errorf(apiErrorRes.Error.Message)
	}

	// Succeeded
	var dfcRes DriveFilesCreate_Response
	err = json.NewDecoder(res.Body).Decode(&dfcRes)
	if err != nil {
		global.Logger.Errorf("Failed to decode response body with error: %v", err)
		return "", "", err
	}

	return dfcRes.ID, dfcRes.URL, nil

}
