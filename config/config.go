package config

type config struct {
	Misskey struct {
		Instance string `yaml:"instance"`
		Token    string `yaml:"token"`
		FolderID string `yaml:"folderId"`
	} `yaml:"misskey"`
	EMail struct {
		Host       []string `yaml:"host"`
		ReaderUrl  string   `yaml:"readerUrl"`
		VerifyDKIM bool     `yaml:"verifyDKIM"`
		SizeLimit  int      `yaml:"sizeLimit"`
		TLS        struct {
			Enable bool   `yaml:"enable"`
			Cert   string `yaml:"cert"`
			Key    string `yaml:"key"`
		} `yaml:"tls"`
	} `yaml:"email"`
	System struct {
		Redis      string `yaml:"redis"`
		Listen     string `yaml:"listen"`
		Production bool   `yaml:"production"`
	} `yaml:"system"`
}

var Config config
