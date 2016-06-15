// package controllers holds the route and all of the controllers
package controllers

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/andgrit/bridge/misc"
	"github.com/andgrit/bridge/configuration"
	"github.com/andgrit/bridge/bridgeTypes"
	"github.com/andgrit/bridge/storage"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"github.com/kr/pretty"
)

var globalStorage *storage.Storage

const OID = "oid"
var store = sessions.NewCookieStore([]byte("something-very-secret"))

func SetConfiguration(appConfiguration *configuration.AppConfiguration) error {
	var err error
	if appConfiguration.DropDatabase {
		globalStorage, err = storage.NewStorage(appConfiguration.DatabaseName)
	} else {
		globalStorage, err = storage.ExistingStorage(appConfiguration.DatabaseName)
	}
	return err
}

func MuxRouter(apiPrefix string) *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix(apiPrefix).Subrouter()
	api.HandleFunc("/version", VesionInfo) // if nothing else shows up return the version info
	api.HandleFunc("/players", Players) // if nothing else shows up return the version info
	api.HandleFunc("/tables", Tables) // if nothing else shows up return the version info

	api.Path("/player").Methods("POST").HandlerFunc(PostPlayer)
	api.Path("/table").Methods("POST").HandlerFunc(PostTable)

	api.Path("/table/{" + OID + "}").Methods("PUT").HandlerFunc(PutTable)
	api.Path("/table/{" + OID + "}").Methods("GET").HandlerFunc(Table)

	http.Handle("/", r)
	return r
}

type VersionInfo struct {
	Application string
	Version     string
	Author      string
}

func VesionInfo(w http.ResponseWriter, r *http.Request) {
	vi := VersionInfo{Application: "bridge", Version:"0.0.1", Author:"Powell Quiring"}
	//////////////////

	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pretty.Println("session:", session)

	// Set some session values.
	session.Values["foo"] = "bar"
	session.Values[42] = 43
	// Save it before we write to the response/return from the handler.
	session.Save(r, w)

	///////////////////
	ret, _ := json.Marshal(&vi)
	fmt.Fprintf(w, string(ret))
}

// create a new player
func PostPlayer(responseWriter http.ResponseWriter, request *http.Request) {
	player := &bridgeTypes.Player{}
	err := misc.RequestBodyUnmarshal(request, player)
	if err == nil {
		err = globalStorage.CreatePlayer(player)
	}
	misc.WriteResponse(responseWriter, player, err)
}

// list the players
func Players(responseWriter http.ResponseWriter, request *http.Request) {
	players, err := globalStorage.GetPlayers()
	misc.WriteResponse(responseWriter, players, err)
}

// create a new table
func PostTable(responseWriter http.ResponseWriter, request *http.Request) {
	table := &bridgeTypes.Table{}
	err := misc.RequestBodyUnmarshal(request, table)
	if err == nil {
		err = globalStorage.CreateTable(table)
	}
	misc.WriteResponse(responseWriter, table, err)
}

type InvalidOidString string
func (oidString InvalidOidString)Error() string {
	return fmt.Sprintf("Invalid OidString used to identify a record.  Expecting a hex string, got:%s", string(oidString))
}

// update a table, replace the whole thing
func PutTable(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	oidString := vars[OID]
	var err error
	newTable := &bridgeTypes.Table{}
	if !bson.IsObjectIdHex(oidString) {
		err = InvalidOidString(oidString)
	} else {
		err := misc.RequestBodyUnmarshal(request, newTable)
		if err == nil {
			oid := bson.ObjectIdHex(oidString)
			oldTable, err := globalStorage.GetTable(oid)
			if err == nil {
				if oldTable.OId != oid {
					err = errors.New("storage table mismatch expecting oid:" + oid.String() + "got:" + oldTable.OId.String())
				} else {
					newTable.OId = oid
					newTable.Version = oldTable.Version
					err = globalStorage.UpdateTable(newTable)
				}
			}
		}
	}
	misc.WriteResponse(responseWriter, newTable, err)
}

// list the tables
func Tables(responseWriter http.ResponseWriter, request *http.Request) {
	tables, err := globalStorage.GetTables()
	misc.WriteResponse(responseWriter, tables, err)
}

// return a table
func Table(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	oidString := vars[OID]
	var err error
	table := &bridgeTypes.Table{}
	if !bson.IsObjectIdHex(oidString) {
		err = InvalidOidString(oidString)
	} else {
		oid := bson.ObjectIdHex(oidString)
		table, err = globalStorage.GetTable(oid)
	}
	misc.WriteResponse(responseWriter, table, err)
}
