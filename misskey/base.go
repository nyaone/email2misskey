package misskey

import (
	"bytes"
	"email2misskey/config"
	"email2misskey/global"
	"encoding/json"
	"fmt"
	"net/http"
)

type Error_Response struct {
	Error struct {
		Message string `json:"message"`
		Code    string `json:"code"`
		ID      string `json:"id"`
		Kind    string `json:"kind"`
	} `json:"error"`
}

func PostAPIRequest[T I_Response | CreatePrivateNote_Response](
	apiEndpointPath string, reqBody any,
) (*T, error) {
	// Prepare request
	apiEndpoint := fmt.Sprintf("https://%s/api/%s", config.Config.Misskey.Instance, apiEndpointPath)

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		global.Logger.Errorf("Failed to marshall request body with error: %v", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewReader(reqBodyBytes))
	if err != nil {
		global.Logger.Errorf("Failed to prepare request with error: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Do request
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		global.Logger.Errorf("Failed to finish request with error: %v", err)
		return nil, err
	}

	// Parse response
	if res.StatusCode == http.StatusOK {
		var resBody T
		err = json.NewDecoder(res.Body).Decode(&resBody)
		if err != nil {
			global.Logger.Errorf("Failed to decode response body with error: %v", err)
			return nil, err
		}

		return &resBody, nil
	} else {
		global.Logger.Errorf("Request failed with code: %d.", res.StatusCode)
		var errBody Error_Response
		err = json.NewDecoder(res.Body).Decode(&errBody)
		if err != nil {
			global.Logger.Errorf("Failed to decode error body with error: %v", err)
			return nil, err
		}

		global.Logger.Errorf("Failed details: %v", errBody)
		return nil, fmt.Errorf(errBody.Error.Message)
	}

}
