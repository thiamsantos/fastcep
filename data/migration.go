package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func getChunkedData(filename string) [][]string {
	dat, err := ioutil.ReadFile(filename)
	checkErr(err)

	r := csv.NewReader(strings.NewReader(string(dat)))

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
	dbfile := "database.sqlite"
	if _, err := os.Stat(dbfile); err == nil {
		err = os.Remove(dbfile)
		checkErr(err)
	}

	db, err := sql.Open("sqlite3", dbfile)
	checkErr(err)

	transaction, err := db.Begin()
	checkErr(err)

	statement, err := transaction.Prepare("create table `states` (`id` integer not null primary key autoincrement, `name` varchar(255) not null, `abbreviation` varchar(2) not null)")
	checkErr(err)

	statement.Exec()
	statement, err = transaction.Prepare("create table `cities` (`id` integer not null primary key autoincrement, `state_id` integer not null, `name` varchar(255) not null, foreign key(`state_id`) references `states`(`id`))")
	checkErr(err)
	statement.Exec()

	statement, err = transaction.Prepare("create table `postal_codes` (`id` integer not null primary key autoincrement, `cep` varchar(8) not null, `neighborhood` varchar(255) not null, `street` varchar(255) not null, `state_id` integer not null, `city_id` integer not null, foreign key(`state_id`) references `states`(`id`), foreign key(`city_id`) references `cities`(`id`))")
	checkErr(err)
	statement.Exec()

	statement, err = transaction.Prepare("create unique index `postal_codes_cep_unique` on `postal_codes` (`cep`)")
	checkErr(err)
	statement.Exec()

	statement, err = transaction.Prepare("create index `postal_codes_state_id_index` on `postal_codes` (`state_id`)")
	checkErr(err)
	statement.Exec()

	statement, err = transaction.Prepare("create index `postal_codes_city_id_index` on `postal_codes` (`city_id`)")
	checkErr(err)
	statement.Exec()

	statesFields := []string{"id", "name", "abbreviation"}
	statesStatement, err := transaction.Prepare(createInsert("states", statesFields))
	checkErr(err)

	states := getChunkedData("data/states.csv")
	for _, state := range states {
		statesStatement.Exec(toSliceInterface(state)...)
	}

	citiesFields := []string{"id", "name", "state_id"}
	citiesStatement, err := transaction.Prepare(createInsert("cities", citiesFields))
	checkErr(err)

	cities := getChunkedData("data/cities.csv")
	for _, city := range cities {
		citiesStatement.Exec(toSliceInterface(city)...)
	}

	postalCodesFields := []string{"cep", "street", "neighborhood", "city_id", "state_id"}
	postalCodesStatement, err := transaction.Prepare(createInsert("postal_codes", postalCodesFields))
	checkErr(err)

	postalCodes := getChunkedData("data/ceps.csv")
	for _, postalCode := range postalCodes {
		postalCodesStatement.Exec(toSliceInterface(postalCode)...)
	}

	err = transaction.Commit()
	checkErr(err)
}
