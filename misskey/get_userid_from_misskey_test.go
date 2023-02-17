package misskey

import (
	"email2misskey/config"
	"testing"
)

func TestGetUserIDFromMisskey(t *testing.T) {
	// Init settings
	config.Config.Misskey.Instance = "nya.one"

	t.Log(getUserIDFromMisskey("candinya"))

}
