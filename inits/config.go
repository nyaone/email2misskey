package inits

import (
	"email2misskey/global"
	"encoding/json"
	"log"
	"os"
)

func Config() error {
	// Read config.json
	f, err := os.Open("config.json")
	if err != nil {
		log.Printf("Failed to open config.json file with error: %v", err)
		return err
	}

	defer f.Close()
	err = json.NewDecoder(f).Decode(&global.Config)
	if err != nil {
		log.Printf("Failed to decode config.json file contents with error: %v", err)
		return err
	}

	log.Printf("Configurations initialized")
	return nil
}
