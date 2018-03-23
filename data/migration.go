package main

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fastcep/src/address"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getChunkedData(filename string) [][]string {
	fi, err := os.Open(filename)
	checkErr(err)
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	defer fz.Close()

	r := csv.NewReader(fz)

	records, err := r.ReadAll()
	checkErr(err)

	return records
}

func toAddress(item []string) address.Address {
	return address.Address{item[0], item[1], item[2], item[3], item[4], item[5]}
}

func main() {
	dbfile := "data.db"
	db, err := bolt.Open(dbfile, 0600, nil)
	checkErr(err)
	defer db.Close()

	postalCodes := getChunkedData("data/data.csv.gz")

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("postal_codes"))
		if err != nil {
			return err
		}

		for _, postalCode := range postalCodes {
			add := toAddress(postalCode)
			encoded, err := json.Marshal(add)
			if err != nil {
				return err
			}
			err = b.Put([]byte(add.CEP), encoded)
			if err != nil {
				return err
			}
		}

		return nil
	})
	checkErr(err)
}
