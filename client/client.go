package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/headmade/backuper/backuper"
	"github.com/headmade/backuper/config"
)

type Client struct {
	URL   string
	Token string
}

func (client *Client) request(method string, action string, o interface{}, params ...interface{}) error {
	var bodyReader io.Reader
	var req *http.Request
	var httpClient = &http.Client{}
	uri := fmt.Sprintf("%s/backend/%s?token=%s", client.URL, action, client.Token)
	if len(params) > 0 {
		json, err := json.Marshal(params[0])
		if err != nil {
			return err
		}
		req, _ = http.NewRequest(method, uri, bytes.NewBuffer(json))
		req.Header.Add("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, uri, bodyReader)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	if o != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(body, o); err != nil {
			return err
		}
	}
	return nil
}

func InitServer(backendAddr, token string) (err error) {
	var clientConfig backuper.ClientConfig
	client := Client{backendAddr, token}
	if err = client.request("GET", "InitServer", &clientConfig); err == nil {
		conf := &config.Config{}
		err = conf.Write(&clientConfig)
	}
	return err
}

func NewClient(backendAddr string) (*Client, error) {
	conf, err := config.New()
	return &Client{backendAddr, conf.Client.Token}, err
}

func (client *Client) Backup(backupResult *backuper.BackupResult) error {
	err := client.request("POST", "Backup", nil, backupResult)
	return err
}

func (client *Client) GetConfig(o interface{}) error {
	err := client.request("GET", "GetConfig", o)
	return err
}

func (client *Client) GetBackupsList(brs *[]backuper.BackupResult) error {
	return client.request("GET", "GetBackups", brs)
}

func (client *Client) GetBackup(backup *backuper.BackupResult, id int) error {
	return client.request("GET", fmt.Sprintf("GetBackup/%d.json", id), backup)
}
