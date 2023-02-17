package misskey

import "email2misskey/config"

type CreatePrivateNote_Request struct {
	I              string   `json:"i"`
	Visibility     string   `json:"visibility"` // "specified"
	VisibleUserIDs []string `json:"visibleUserIds"`
	CW             string   `json:"cw"`        // Summary
	Text           string   `json:"text"`      // Message
	LocalOnly      bool     `json:"localOnly"` // true
	FileIDs        []string `json:"fileIds"`
}

type CreatePrivateNote_Response struct {
	// Ignore response
}

func CreatePrivateNote(userIDs []string, summary string, message string, attachmentId string) error {
	_, err := PostAPIRequest[CreatePrivateNote_Response]("notes/create", &CreatePrivateNote_Request{
		I:              config.Config.Misskey.Token,
		Visibility:     "specified",
		VisibleUserIDs: userIDs,
		CW:             summary,
		Text:           message,
		LocalOnly:      true,
		FileIDs:        []string{attachmentId},
	})
	return err
}
