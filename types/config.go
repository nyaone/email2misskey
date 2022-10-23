package types

type Config struct {
	Misskey struct {
		Instance string  `json:"instance"`
		Token    string  `json:"token"`
		FolderID *string `json:"folderId"`
	} `json:"misskey"`
	EMail struct {
		Host []string `json:"host"`
	} `json:"email"`
	System struct {
		Redis      string `json:"redis"`
		Production bool   `json:"production"`
	} `json:"system"`
}
