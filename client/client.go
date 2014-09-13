package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/headmade/backuper/backuper"
)

type Client struct {
	Url   string
	Token string
}

func InitServer(backendAddr, token string) error {
	resp, err := http.Get("http://" + backendAddr + "/backend/InitServer?token=" + token)
	defer resp.Body.Close()
	if err == nil {
		if resp.StatusCode != 200 {
			return errors.New(resp.Status)
		}
		var config Config
		body, err := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &config); err == nil {
			err = WriteConfig(&config, configPath())
		}
	}
	return err
}

func Get(backendAddr string) (*Client, error) {
	config, err := LoadConfig(configPath())
	return &Client{backendAddr, config.Token}, err
}

func (client *Client) Backup(backupResult *backuper.BackupResult) error {
	return nil
}

func (client *Client) GetConfig() ([]byte, error) {
	response := []byte{}
	resp, err := http.Get("http://" + client.Url + "/backend/GetConfig?token=" + client.Token)
	defer resp.Body.Close()
	if err == nil {
		if resp.StatusCode != 200 {
			return response, errors.New(resp.Status)
		}
		response, err = ioutil.ReadAll(resp.Body)
	}
	return response, nil
}
