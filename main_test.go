package main

import (
	"testing"
	"net/http/httptest"
	"github.com/andgrit/bridge/controllers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"encoding/json"
	"github.com/andgrit/bridge/bridgeTypes"
	"github.com/andgrit/bridge/configuration"
	"github.com/andgrit/bridge/misc"
	"strings"
	"github.com/andgrit/bridge/storage"
	"gopkg.in/mgo.v2/bson"
)

func TestRealServer(t *testing.T) {
	appConfiguration, _ := configuration.DefaultConfiguration()
	appConfiguration.DatabaseName = storage.TEST_DATABSE
	appConfiguration.DropDatabase = true
	storage.SetConfiguration(appConfiguration)
	controllers.SetConfiguration(appConfiguration)

	bridgeServer := httptest.NewServer(controllers.MuxRouter("/"))
	defer bridgeServer.Close()

	var versionInfo controllers.VersionInfo
	err := getType(bridgeServer.URL, &versionInfo)
	assert.NoError(t, err)
	assert.Equal(t, "bridge", versionInfo.Application)

	var players []bridgeTypes.Player
	err = getType(bridgeServer.URL + "/players", &players)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(players))

	const USER = "sam q"
	player := &bridgeTypes.Player{Username:USER}
	assert.Equal(t, bson.ObjectId(""), player.OId)
	err = postType(bridgeServer.URL + "/player", player, player)
	assert.Equal(t, USER, player.Username)
	assert.NotEqual(t, bson.ObjectId(""), player.OId)

	err = getType(bridgeServer.URL + "/players", &players)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(players))
	player = &players[0]
	assert.Equal(t, USER, player.Username)
	assert.NotEqual(t, "", player.OId)

	var tables []*bridgeTypes.Table
	err = getType(bridgeServer.URL + "/tables", &tables)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tables))

	table := &bridgeTypes.Table{}
	assert.Equal(t, bson.ObjectId(""), table.OId)
	err = postType(bridgeServer.URL + "/table", table, table)
	assert.NotEqual(t, bson.ObjectId(""), table.OId)
	oid := table.OId

	err = getType(bridgeServer.URL + "/tables", &tables)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tables))
	table = tables[0]
	assert.Equal(t, oid, table.OId)
}

func getType(url string, result interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	return misc.ResponseBodyUnmarshal(res, result)
}

func postType(url string, requestBody interface{}, returnBody interface{}) error {
	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	res, err := http.Post(url, "application/json", strings.NewReader(string(body)))
	return misc.ResponseBodyUnmarshal(res, returnBody)

}
