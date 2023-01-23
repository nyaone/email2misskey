package jobs

import (
	"bytes"
	"email2misskey/consts"
	"email2misskey/global"
	"email2misskey/types"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendMessage(userID string, emailFileID string, emailSubject string, emailSender string) error {
	messageText := fmt.Sprintf(consts.MessageTemplate, emailSubject, emailSender)
	reqBody := types.MessageCreate_Request{
		I:      global.Config.Misskey.Token,
		UserID: userID,
		Text:   messageText,
		FileID: emailFileID,
	}
	reqBodyBytes, err := json.Marshal(&reqBody)
	if err != nil {
		global.Logger.Errorf("Failed to marshall ShowUser request with error: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/messaging/messages/create", global.Config.Misskey.Instance), bytes.NewReader(reqBodyBytes))
	if err != nil {
		global.Logger.Errorf("Failed to initialize request with error: %v", err)
		return err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		global.Logger.Errorf("Failed to finish request with error: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if res.StatusCode != http.StatusOK {
		// Failed
		global.Logger.Errorf("Request failed.")
		var apiErrorRes types.ErrorResponse
		err := json.NewDecoder(res.Body).Decode(&apiErrorRes)
		if err != nil {
			global.Logger.Errorf("Failed to decode response body with error: %v", err)
			return err
		}
		global.Logger.Errorf("Request failed with error: %v", apiErrorRes)
		return fmt.Errorf(apiErrorRes.Error.Message)
	}

	// Succeeded
	// Ignore response

	global.Logger.Debugf("Message sent successfully")
	return nil
}
