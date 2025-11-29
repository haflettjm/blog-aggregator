package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	DB_URL     string `json:"db_url"`
	CURR_UNAME string `json:"current_user_name"`
}

func (c *Config) Read(filepath string) {
	home, _ := os.UserHomeDir()
	path := home + "/" + filepath

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error Reading file: \n %v", err)
	}

	err = json.Unmarshal(file, &c)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: \n %v", err)
	}

}

func (c *Config) SetUser(user string) {
	c.CURR_UNAME = user
}
