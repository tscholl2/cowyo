package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var db *bolt.DB
var open bool

// Open to create the database and open
func Open(filename string) error {
	var err error
	config := &bolt.Options{Timeout: 30 * time.Second}
	db, err = bolt.Open(filename, 0600, config)
	if err != nil {
		fmt.Println("Opening BoltDB timed out")
		log.Fatal(err)
	}
	open = true
	return nil
}

// Close database
func Close() {
	open = false
	db.Close()
}

// WikiData is data for storing in DB
type WikiData struct {
	Title       string
	CurrentText string
	Diffs       []string
	Timestamps  []string
	Encrypted   bool
	Locked      string
}

func getCurrentText(title string, version int) (string, []versionsInfo, bool, time.Duration, bool, string, int) {
	Open(RuntimeArgs.DatabaseLocation)
	defer Close()
	title = strings.ToLower(title)
	var vi []versionsInfo
	totalTime := time.Now().Sub(time.Now())
	isCurrent := true
	currentText := ""
	encrypted := false
	locked := ""
	currentVersionNum := -1
	if !open {
		return currentText, vi, isCurrent, totalTime, encrypted, locked, currentVersionNum
	}
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("datas"))
		if b == nil {
			return fmt.Errorf("db must be opened before loading")
		}
		k := []byte(title)
		val := b.Get(k)
		if val == nil {
			return nil
		}
		var p WikiData
		err = p.decode(val)
		if err != nil {
			return err
		}
		currentText = p.CurrentText
		encrypted = p.Encrypted
		locked = p.Locked
		currentVersionNum = len(p.Diffs) - 1
		if version > -1 && version < len(p.Diffs) {
			// get that version of text instead
			currentText = rebuildTextsToDiffN(p, version)
			isCurrent = false
		}
		vi, totalTime = getImportantVersions(p)
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get WikiData: %s", err)
	}
	return currentText, vi, isCurrent, totalTime, encrypted, locked, currentVersionNum
}

func (p *WikiData) load(title string) error {
	title = strings.ToLower(title)
	if !open {
		Open(RuntimeArgs.DatabaseLocation)
		defer Close()
	}
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("datas"))
		if b == nil {
			return nil
		}
		k := []byte(title)
		val := b.Get(k)
		if val == nil {
			// make new one
			p.Title = title
			p.CurrentText = ""
			p.Diffs = []string{}
			p.Timestamps = []string{}
			return nil
		}
		err = p.decode(val)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get WikiData: %s", err)
		return err
	}
	return nil
}

func (p *WikiData) save(newText string) error {
	if !open {
		Open(RuntimeArgs.DatabaseLocation)
		defer Close()
	}
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("datas"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		// find diffs
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(p.CurrentText, newText, true)
		delta := dmp.DiffToDelta(diffs)
		p.CurrentText = newText
		p.Timestamps = append(p.Timestamps, time.Now().Format(time.ANSIC))
		p.Diffs = append(p.Diffs, delta)
		enc, err := p.encode()
		if err != nil {
			return fmt.Errorf("could not encode WikiData: %s", err)
		}
		p.Title = strings.ToLower(p.Title)
		err = bucket.Put([]byte(p.Title), enc)
		if err != nil {
			return fmt.Errorf("could add to bucket: %s", err)
		}
		return err
	})
	// // Add the new name to the programdata so its not randomly generated
	// if err == nil && len(p.Timestamps) > 0 && len(p.CurrentText) > 0 {
	// 	err2 := db.Update(func(tx *bolt.Tx) error {
	// 		b := tx.Bucket([]byte("programdata"))
	// 		id, _ := b.NextSequence()
	// 		idInt := int(id)
	// 		return b.Put(itob(idInt), []byte(p.Title))
	// 	})
	// 	if err2 != nil {
	// 		return fmt.Errorf("could not add to programdata: %s", err)
	// 	}
	// }
	return err
}

func (p *WikiData) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (p *WikiData) decode(data []byte) error {
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	return nil
}
