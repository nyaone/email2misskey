package misskey

import (
	"context"
	"email2misskey/config"
	"email2misskey/consts"
	"email2misskey/global"
	"fmt"
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
	cacheKey := fmt.Sprintf(consts.CacheUsernameTemplate, config.Config.Misskey.Instance, username)
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
	cacheKey := fmt.Sprintf(consts.CacheUsernameTemplate, config.Config.Misskey.Instance, username)
	global.Redis.Set(context.Background(), cacheKey, userID, 0)
}

type UserShow_Request struct {
	Username string  `json:"username"`
	Host     *string `json:"host"` // Null
}

func getUserIDFromMisskey(username string) (string, bool, error) {
	// Get data from Misskey

	res, err := PostAPIRequest[I_Response]("users/show", &UserShow_Request{
		Username: username,
		Host:     nil,
	})
	if err != nil {
		global.Logger.Errorf("Failed to get user id for @%s with error: %v", username, err)
		return "", false, err
	}

	// Exist
	global.Logger.Debugf("Found user @%s (%s)", res.Username, res.ID)

	saveUserIDToCache(username, res.ID)

	return res.ID, true, nil
}
