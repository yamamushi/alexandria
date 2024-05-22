package main

import (
	"github.com/asdine/storm/v3"
	"github.com/ciehanski/libgen-cli/libgen"
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
	Data libgen.Book
}

func (adb *AlexandriaDB) StoreRecord(record *libgen.Book) (err error) {
	err = adb.db.Save(&BookRecord{ID: record.Md5, Data: *record})
	if err != nil {
		return err
	}
	return nil
}

func (adb *AlexandriaDB) GetRecord(id string) (record *libgen.Book, err error) {
	var bookRecord BookRecord
	err = adb.db.One("ID", id, &bookRecord)
	if err != nil {
		return nil, err
	}
	return &bookRecord.Data, nil
}
