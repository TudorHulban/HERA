package hera

import (
	"database/sql"
	"errors"
)

type DBConnection struct {
	Active    bool
	DBHandler *sql.DB
	// unix timestamp of when connection was lost, for clean up after some time, maybe other processes still refer it
	LastStatusChange int64
}

type DBConnManager struct {
	Connections map[string]*DBConnection
}

// newManager - constructor for db manager
func newManager() *DBConnManager {
	instance := new(DBConnManager)
	instance.Connections = make(map[string]*DBConnection)
	return instance
}

// addConnection - persists connection till inactive
func (m *DBConnManager) addConnection(pID string, pRDBMS RDBMS) error {
	if _, exists := m.Connections[pID]; exists {
		return errors.New("handler exists")
	}
	handler, errNewHandler := pRDBMS.NewConnection()
	if errNewHandler != nil {
		return errNewHandler
	}
	m.Connections[pID].DBHandler = handler
	m.Connections[pID].Active = true
	return nil
}

// deleteConnection - delete connction when issues
func (m *DBConnManager) deleteConnection(pID string) error {
	if _, exists := m.Connections[pID]; !exists {
		return errors.New("connection does not exist")
	}
	delete(m.Connections, pID)
	return nil
}

// requestConnection - provides connection handler based on connection ID
func (m *DBConnManager) requestConnection(pID string) (*sql.DB, error) {
	if _, exists := m.Connections[pID]; !exists {
		return nil, errors.New("connection does not exist")
	}
	return m.Connections[pID].DBHandler, nil
}
