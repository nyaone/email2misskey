package inits

import (
	"email2misskey/config"
	"email2misskey/global"
	"email2misskey/misskey"
)

func Misskey() error {

	res, err := misskey.I(config.Config.Misskey.Token)
	if err != nil {
		global.Logger.Errorf("Failed to check token with error: %v", err)
		return err
	}

	global.Logger.Infof("System activated, I'm %s", res.Name)
	return nil
}
