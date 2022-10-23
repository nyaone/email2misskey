package jobs

import (
	"email2misskey/global"
	"testing"
)

func TestGetUserIDFromMisskey(t *testing.T) {
	// Init settings
	global.Config.Misskey.Instance = "nya.one"

	t.Log(getUserIDFromMisskey("candinya"))

}
