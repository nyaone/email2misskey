package inits

import (
	"bytes"
	"email2misskey/global"
	"email2misskey/types"
	"encoding/json"
	"fmt"
	"net/http"
)

func Misskey() error {
	// Try to access endpoint `i` and get self info
	reqBody := types.I_Request{
		I: global.Config.Misskey.Token,
	}
	reqBodyBytes, err := json.Marshal(&reqBody)
	if err != nil {
		global.Logger.Errorf("Failed to marshall I request with error: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/i", global.Config.Misskey.Instance), bytes.NewReader(reqBodyBytes))
	if err != nil {
		global.Logger.Errorf("Failed to initialize request with error: %v", err)
		return err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		global.Logger.Errorf("Failed to finish request with error: %v", err)
		return err
	}

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
	var iRes types.I_Response
	err = json.NewDecoder(res.Body).Decode(&iRes)
	if err != nil {
		global.Logger.Errorf("Failed to decode response body with error: %v", err)
		return err
	}

	global.Logger.Infof("System activated, I'm %s", iRes.Name)
	return nil
}
