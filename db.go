package main

import (
	"github.com/asdine/storm/v3"
	"go.etcd.io/bbolt"
	"time"
)

type AlexandriaDB struct {
	db *storm.DB
}

func (adb *AlexandriaDB) OpenDB() (err error) {
	adb.db, err = storm.Open("alex.db", storm.BoltOptions(0600, &bbolt.Options{Timeout: 3 * time.Second}))
	if err != nil {
		return err
	}
	return
}

func (adb *AlexandriaDB) CloseDB() (err error) {
	err = adb.db.Close()
	if err != nil {
		return err
	}
	return nil
}

type BookRecord struct {
	ID   string `storm:"id"`
	Data Book
}

func (adb *AlexandriaDB) StoreRecord(record *Book) (err error) {
	err = adb.db.Save(&BookRecord{ID: record.ID, Data: *record})
	if err != nil {
		return err
	}
	return nil
}

func (adb *AlexandriaDB) GetRecord(id string) (record *Book, err error) {
	var bookRecord BookRecord
	err = adb.db.One("ID", id, &bookRecord)
	if err != nil {
		return nil, err
	}
	return &bookRecord.Data, nil
}
