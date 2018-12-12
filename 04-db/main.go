package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const CREATE_TEST_TABLE = `
CREATE TABLE test (
  id BIGINT NOT NULL AUTO_INCREMENT,
  firstname varchar(255),
  lastname varchar(255),
PRIMARY KEY (id)
)`

type test struct {
	id        int64
	firstname string
	lastname  string
}

func main() {
	mysql, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s", "contact", "contact", "tcp(localhost:3306)/contact"))
	if err != nil {
		log.Println("Error opening database tcp(localhost:3306)/contact")
		panic(fmt.Sprintf("%v", err))
	}
	if err = mysql.Ping(); err != nil {
		log.Fatalf("Error accessing database: %v\n", err)
		panic(fmt.Sprintf("%v", err))
	}
	mysql.Exec("DROP TABLE test")

	_, err = mysql.Exec(CREATE_TEST_TABLE)
	if err != nil {
		log.Fatalln("Error creating contact table")
		panic(fmt.Sprintf("%v", err))
	} else {
		log.Println("Contact Table successfully created....")
	}

	ersterTest := test{firstname: "Kunibert", lastname: "Knäuel"}
	result, err := mysql.Exec("INSERT INTO test (firstname, lastname) VALUES(?,?)", ersterTest.firstname, ersterTest.lastname)
	if err != nil {
		log.Fatalf("Error inserting contact: %v", err)
	}
	ersterTest.id, err = result.LastInsertId()
	if err != nil {
		log.Fatalf("Error contact id: %v", err)
	}

	zweiterTest := test{firstname: "Kunigunde", lastname: "Käfer"}
	result, err = mysql.Exec("INSERT INTO test (firstname, lastname) VALUES(?,?)", zweiterTest.firstname, zweiterTest.lastname)
	if err != nil {
		log.Fatalf("Error inserting contact: %v", err)
	}
	zweiterTest.id, err = result.LastInsertId()
	if err != nil {
		log.Fatalf("Error contact id: %v", err)
	}

	contact := test{}
	err = mysql.QueryRow("SELECT id,firstname, lastname FROM test WHERE id=?", ersterTest.id).Scan(
		&contact.id,
		&contact.firstname,
		&contact.lastname,
	)
	if err != nil {
		log.Fatalf("Error reading one row: %v\n", err)
	}
	log.Printf("%v\n", contact)

	rows, err := mysql.Query("SELECT id, firstname, lastname FROM test ORDER BY lastname, firstname")
	if err != nil {
		log.Fatalf("Error fetching contact table: %v", err)
	}
	defer rows.Close()
	contacts := make([]*test, 0)
	for rows.Next() {
		contact = test{}
		err = rows.Scan(
			&contact.id,
			&contact.firstname,
			&contact.lastname,
		)
		if err != nil {
			log.Fatalf("Error reading one row: %v\n", err)
		}
		log.Println(contact)
		contacts = append(contacts, &contact)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	log.Println(contacts)
}
