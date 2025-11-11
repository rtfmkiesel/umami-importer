package db

import (
	"fmt"

	logger "github.com/rtfmkiesel/kisslog"
	"go.etcd.io/bbolt"
)

const bucketName = "hashes"

var (
	currentDB *bbolt.DB
	log       = logger.New("db")
)

func Open(dbPath string) error {
	if currentDB != nil {
		return log.NewError("database already open")
	}

	log.Debug("Opening database '%s'", dbPath)

	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return log.NewError("failed to open database: %s", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		if closeErr := db.Close(); closeErr != nil {
			log.Warning("failed to close database after bucket creation error: %v", closeErr)
		}
		return log.NewError("failed to initialize database: %w", err)
	}

	currentDB = db
	log.Info("Database '%s' opened successfully", dbPath)
	return nil
}
func Close() error {
	log.Debug("Closing database")

	if currentDB != nil {
		return currentDB.Close()
	}

	return nil
}

func IsKnown(checksum string) (bool, error) {
	if currentDB == nil {
		return false, log.NewError("database is not open")
	}

	var exists bool
	err := currentDB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return log.NewError("bucket not found")
		}

		v := b.Get([]byte(checksum))
		exists = v != nil
		return nil
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}

func Add(checksum string) error {
	if currentDB == nil {
		return log.NewError("database is not open")
	}

	err := currentDB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return log.NewError("bucket not found")
		}

		return b.Put([]byte(checksum), []byte{1})
	})

	if err != nil {
		return log.NewError(err)
	}

	return nil
}
