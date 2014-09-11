package client

import "github.com/headmade/backuper/backuper"

type Client struct {
	Url   string
	Token string
}

func InitServer(token string) error {
	config := &Config{token}
	err := error(nil)
	if err == nil {
		err = WriteConfig(config, configPath())
	}
	return err
}

func Get(url string) (*Client, error) {
	config, err := LoadConfig(configPath())
	return &Client{url, config.Token}, err
}

func (client *Client) Backup(backupResult *backuper.BackupResult) error {
	return nil
}

func (client *Client) GetConfig() ([]byte, error) {
	response := []byte("{}")
	return response, nil
}
