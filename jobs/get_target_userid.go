package jobs

import (
	"bytes"
	"context"
	"email2misskey/consts"
	"email2misskey/global"
	"email2misskey/types"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetTargetUserID(username string) (string, bool, error) {
	id, err := getUserIDFromCache(username)
	if err == nil {
		return id, true, nil
	}
	// Else
	return getUserIDFromMisskey(username)
}

func getUserIDFromCache(username string) (string, error) {
	// Get data from cache
	cacheKey := fmt.Sprintf(consts.CacheUsernameTemplate, global.Config.Misskey.Instance, username)
	// Check if in cache
	exist, err := global.Redis.Exists(context.Background(), cacheKey).Result()
	if err != nil {
		global.Logger.Errorf("Failed to check user existence in cache")
		return "", err
	}
	if exist == 0 {
		global.Logger.Debugf("User not exist in cache")
		return "", fmt.Errorf("user not in cache")
	}

	userID, err := global.Redis.Get(context.Background(), cacheKey).Result()
	if err != nil {
		global.Logger.Errorf("Failed to get user ID from cache")
		return "", err
	}
	return userID, nil

}

func saveUserIDToCache(username string, userID string) {
	cacheKey := fmt.Sprintf(consts.CacheUsernameTemplate, global.Config.Misskey.Instance, username)
	global.Redis.Set(context.Background(), cacheKey, userID, 0)
}

func getUserIDFromMisskey(username string) (string, bool, error) {
	// Get data from Misskey
	reqBody := types.UserShow_Request{
		Username: username,
		Host:     nil,
	}
	reqBodyBytes, err := json.Marshal(&reqBody)
	if err != nil {
		global.Logger.Errorf("Failed to marshall ShowUser request with error: %v", err)
		return "", false, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/users/show", global.Config.Misskey.Instance), bytes.NewReader(reqBodyBytes))
	if err != nil {
		global.Logger.Errorf("Failed to initialize request with error: %v", err)
		return "", false, err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		global.Logger.Errorf("Failed to finish request with error: %v", err)
		return "", false, err
	}

	if res.StatusCode != http.StatusOK {
		// User doesn't exist
		return "", false, nil
	} else {
		// Succeeded
		var iRes types.I_Response
		err = json.NewDecoder(res.Body).Decode(&iRes)
		if err != nil {
			global.Logger.Errorf("Failed to decode response body with error: %v", err)
			return "", false, err
		}

		// Exist
		global.Logger.Debugf("Found user @%s (%s)", iRes.Username, iRes.ID)

		saveUserIDToCache(username, iRes.ID)

		return iRes.ID, true, nil
	}
}
