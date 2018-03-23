package main

import (
	"compress/gzip"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

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

func createInsert(table string, fields []string) string {
	insertStatement := fmt.Sprintf("INSERT INTO `%s` (", table)

	for i, field := range fields {
		insertStatement += fmt.Sprintf("`%s`", field)
		if i != len(fields)-1 {
			insertStatement += ","
		}
	}
	insertStatement += ") VALUES ("
	for i := range fields {
		insertStatement += "?"

		if i != len(fields)-1 {
			insertStatement += ","
		}
	}

	insertStatement += ")"
	return insertStatement
}

func toSliceInterface(a []string) []interface{} {
	s := make([]interface{}, len(a))

	for i, v := range a {
		s[i] = v
	}
	return s
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbfile := "data.db"
	if _, err := os.Stat(dbfile); err == nil {
		err = os.Remove(dbfile)
		checkErr(err)
	}

	db, err := sql.Open("sqlite3", dbfile)
	checkErr(err)

	transaction, err := db.Begin()
	checkErr(err)

	statement, err := transaction.Prepare("create table `postal_codes` (`id` integer not null primary key autoincrement, `cep` varchar(8) not null, `neighborhood` varchar(255) not null, `street` varchar(255) not null, `state` varchar(255) not null, `uf` varchar(2) not null, `city` varchar(255) not null)")
	checkErr(err)
	statement.Exec()

	statement, err = transaction.Prepare("create unique index `postal_codes_cep_unique` on `postal_codes` (`cep`)")
	checkErr(err)
	statement.Exec()

	postalCodesFields := []string{"cep", "street", "neighborhood", "city", "state", "uf"}
	postalCodesStatement, err := transaction.Prepare(createInsert("postal_codes", postalCodesFields))
	checkErr(err)

	postalCodes := getChunkedData("data/data.csv.gz")
	checkErr(err)

	for _, postalCode := range postalCodes {
		postalCodesStatement.Exec(toSliceInterface(postalCode)...)
	}

	err = transaction.Commit()
	checkErr(err)
}
