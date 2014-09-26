package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/headmade/backuper/backuper"
	"github.com/headmade/backuper/config"
)

type Client struct {
	Url   string
	Token string
}

func InitServer(backendAddr, token string) error {
	resp, err := http.Get("http://" + backendAddr + "/backend/InitServer?token=" + token)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return errors.New(resp.Status)
		}
		var clientConfig backuper.ClientConfig
		body, err := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &clientConfig); err == nil {
			conf := &config.Config{}
			err = conf.Write(&clientConfig)
		}
	}
	return err
}

func Get(backendAddr string) (*Client, error) {
	conf, err := config.New()
	return &Client{backendAddr, conf.Client.Token}, err
}

func (client *Client) Backup(backupResult *backuper.BackupResult) error {
	json, err := json.Marshal(backupResult)
	response, err := http.Post("http://"+client.Url+"/backend/Backup?token="+client.Token, "application/json", bytes.NewBuffer(json))
	if err == nil && response.StatusCode != 200 {
		err = errors.New(response.Status)
	}
	return err
}

func (client *Client) GetConfig() ([]byte, error) {
	response := []byte{}
	resp, err := http.Get("http://" + client.Url + "/backend/GetConfig?token=" + client.Token)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return response, errors.New(resp.Status)
		}
		response, err = ioutil.ReadAll(resp.Body)
	}
	return response, err
}
