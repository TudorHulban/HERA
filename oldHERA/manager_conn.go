/*
Package hera assists with easy RDBMS ops.
Entry point to a db connection is the connection manager.
Connection is *sql.DB.
*/
package hera

import (
	"database/sql"
	"errors"
)

// DBConnection is part of db connection manager pool.
type DBConnection struct {
	Active    bool
	DBHandler *sql.DB
	// unix timestamp of when connection was lost, for clean up after some time, maybe other processes still refer it
	LastStatusChange int64
}

// DBConnManager is pool of db connections.
type DBConnManager struct {
	Connections map[string]*DBConnection
}

// NewDBConnManager - constructor for db manager
func NewDBConnManager() *DBConnManager {
	instance := new(DBConnManager)
	instance.Connections = make(map[string]*DBConnection)
	return instance
}

// AddConnection - persists connection till inactive
func (m *DBConnManager) AddConnection(pCODE string, pRDBMS RDBMS) error {
	_, exists := m.Connections[pCODE]
	if exists {
		return errors.New("handler exists")
	}
	handler, errNewHandler := pRDBMS.NewConnection()
	if errNewHandler != nil {
		return errNewHandler
	}
	m.Connections[pCODE].DBHandler = handler
	m.Connections[pCODE].Active = true
	return nil
}

// DeleteConnection - delete connction when issues, connection inactive for some time
func (m *DBConnManager) DeleteConnection(pCODE string) error {
	_, exists := m.Connections[pCODE]
	if !exists {
		return errors.New("connection does not exist")
	}
	delete(m.Connections, pCODE)
	return nil
}

// RequestConnection - provides connection handler based on connection CODE
func (m *DBConnManager) RequestConnection(pCODE string) (*sql.DB, error) {
	_, exists := m.Connections[pCODE]
	if !exists {
		return nil, errors.New("connection does not exist")
	}
	if !m.Connections[pCODE].Active {

		m.DeleteConnection(pCODE)
		return nil, errors.New("connection does not exist")
	}
	return m.Connections[pCODE].DBHandler, nil
}
