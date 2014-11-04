package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/headmade/backuper/backuper"
	"github.com/headmade/backuper/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigPath(t *testing.T) {
	conf := config.Config{Client: &backuper.ClientConfig{Token: ""}}
	config.WriteConfig(&conf)
	client, _ := NewClient("test.ru")
	assert.Equal(t, client, &Client{"test.ru", ""})
}

func TestBackup(t *testing.T) {
	backupResult := backuper.BackupResult{Status: "complete"}
	conf := config.Config{Client: &backuper.ClientConfig{Token: ""}}
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var br backuper.BackupResult
		json.Unmarshal(body, &br)
		assert.Equal(t, backupResult.Status, br.Status)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	client := Client{server.URL, conf.Client.Token}
	err := client.Backup(&backupResult)
	assert.Nil(t, err)
}
