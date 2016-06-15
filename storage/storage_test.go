package storage

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/andgrit/bridge/bridgeTypes"
	"gopkg.in/mgo.v2/bson"
	"github.com/kr/pretty"
	"github.com/andgrit/bridge/configuration"
)

func xTestPlayerStorage(t *testing.T) {
	storage, err := NewStorage(TEST_DATABSE)
	assert.Nil(t, err)

	players, err := storage.GetPlayers()
	assert.Nil(t, err)
	assert.NotNil(t, players)
	assert.Equal(t, 0, len(players))

	player := &bridgeTypes.Player{Username:"Sam"}
	err = storage.CreatePlayer(player)
	assert.Nil(t, err)

	players, err = storage.GetPlayers()
	assert.Nil(t, err)
	assert.NotNil(t, players)
	assert.Equal(t, 1, len(players))
	storedPlayer := players[0]
	assert.Equal(t, "Sam", storedPlayer.Username)
}

func TestTableStorage(t *testing.T) {
	SetConfiguration(&configuration.AppConfiguration{MongoIp:"localhost"})
	storage, err := NewStorage(TEST_DATABSE)
	assert.Nil(t, err)

	table  := &bridgeTypes.Table{}
	err = storage.CreateTable(table)
	assert.NoError(t, err)
	assert.Equal(t, 0, table.Version)

	n := &bridgeTypes.Player{Username:"North"}
	err = storage.CreatePlayer(n)
	assert.NoError(t, err)

	table.Players = []bson.ObjectId{n.OId}
	err = storage.UpdateTable(table)
	assert.NoError(t, err)
	assert.Equal(t, 1, table.Version)

	e := &bridgeTypes.Player{Username:"East"}
	err = storage.CreatePlayer(e)
	assert.NoError(t, err)

	table.Players = append(table.Players, e.OId)
	err = storage.UpdateTable(table)
	assert.NoError(t, err)
	assert.Equal(t, 2, table.Version)

	table.Version = 1	// different thread wanted to update, it should fail
	err = storage.UpdateTable(table)
	pretty.Println(err)
	assert.Error(t, err)

	table, err = storage.GetTable(table.OId)
	assert.NoError(t, err)
	assert.Equal(t, 2, table.Version)

}
