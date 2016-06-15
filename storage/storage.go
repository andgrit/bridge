// package storage implements persistence
package storage

import "github.com/andgrit/bridge/bridgeTypes"
import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/andgrit/bridge/configuration"
)

const TEST_DATABSE = "test"
const TABLE_COLLECTION = "table"
const PLAYER_COLLECTION = "player"

var session *mgo.Session

func SetConfiguration(appConfiguration *configuration.AppConfiguration) error {
	var err error
	session, err = mgo.Dial(appConfiguration.MongoIp)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	return err
}

type Storage struct {
	tableCollection *mgo.Collection
	playerCollection *mgo.Collection
	database *mgo.Database
}

// Drop the database
func NewStorage(database string) (*Storage, error) {
	if database != TEST_DATABSE{
		panic("do not screw up the production database")
	}
	storage := Storage{database: session.DB(database)}
	err := storage.database.DropDatabase() // new
	if err != nil {
		return &storage, err
	}
	err = storage.existingStorage()
	return &storage, err
}

func ExistingStorage(database string) (*Storage, error) {
	storage := Storage{database: session.DB(database)}
	err := storage.existingStorage()
	return &storage, err
}

// TODO add code to verify the collections exist in the database?
func (storage *Storage)existingStorage() error {
	storage.createCollections()
	return nil
}

func (storage *Storage)createCollections() {
	storage.tableCollection = storage.database.C(TABLE_COLLECTION)
	storage.playerCollection = storage.database.C(PLAYER_COLLECTION)
}

// NewMatch returns a mongo document for a new match
func NewMatch() (*bridgeTypes.Match) {
	return &bridgeTypes.Match{}
}

// store a new player, the OId is not required and will be ignored if provided
func (storage *Storage)CreatePlayer(player *bridgeTypes.Player) error {
	player.OId = bson.NewObjectId()
	err := storage.playerCollection.Insert(player)
	return err
}

//
func (storage *Storage)GetPlayers() ([]*bridgeTypes.Player, error){
	players := []*bridgeTypes.Player{}
	err := storage.playerCollection.Find(nil).All(&players)
	return players, err
}

// create a new table, the OId and Version is not required and will be ignored if provided
func (storage *Storage)CreateTable(table *bridgeTypes.Table) error {
	table.OId = bson.NewObjectId()
	table.Version = 0
	err := storage.tableCollection.Insert(table)
	return err
}

// return all of the tables
func (storage *Storage)GetTables() ([]*bridgeTypes.Table, error){
	tables := []*bridgeTypes.Table{}
	err := storage.tableCollection.Find(nil).All(&tables)
	return tables, err
}

// replace the table identified by table.{OId, Version} with a new value
// this is expected to fail if another writer beat you to the punch
func (storage *Storage)UpdateTable(table *bridgeTypes.Table) error {
	selector := bson.M{"_id":table.OId, "version":table.Version}
	table.Version++
	err := storage.tableCollection.Update(selector, table)
	return err
}

// return a table
func (storage *Storage)GetTable(oId bson.ObjectId) (*bridgeTypes.Table, error) {
	table := &bridgeTypes.Table{}
	err := storage.tableCollection.FindId(oId).One(table)
	return table, err
}


