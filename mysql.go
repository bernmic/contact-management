package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

const CREATE_CONTACT_TABLE = `
CREATE TABLE contact (
  id BIGINT NOT NULL AUTO_INCREMENT,
  firstname varchar(255),
  lastname varchar(255),
  company varchar(255),
  address1 varchar(255),
  address2 varchar(255),
  zipcode varchar(30),
  city varchar(255),
  country varchar(255),
  tag varchar(1024),
  email varchar(255),
  web varchar(255),
  birthday date,
PRIMARY KEY (id)
)`

const CREATE_PHONE_TABLE = `
CREATE TABLE phone (
  id BIGINT NOT NULL AUTO_INCREMENT,
  name varchar(255),
  number varchar(255),
  contact_id bigint,
PRIMARY KEY (id),
FOREIGN KEY (contact_id) REFERENCES contact(id)
)`

type DB struct {
	sql.DB
}

type TX struct {
	sql.Tx
}

func New() (*DB, error) {
	mysql, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s", "contact", "contact", "tcp(localhost:3306)/contact"))
	if err != nil {
		log.Errorln("Error opening database tcp(localhost:3306)/contact")
		panic(fmt.Sprintf("%v", err))
	}
	if err = mysql.Ping(); err != nil {
		log.Errorf("Error accessing database: %v\n", err)
		panic(fmt.Sprintf("%v", err))
	}
	db := &DB{*mysql}
	db.Exec("DROP TABLE phone")
	db.Exec("DROP TABLE contact")
	created := db.initializeContact()
	db.initializePhone()

	if created {
		db.insertTestData()
	}
	return db, nil
}

/*----------------------------------------------------------------------------------------*/

func (db *DB) initializeContact() bool {
	_, err := db.Query("SELECT 1 FROM contact LIMIT 1")
	if err != nil {
		log.Info("Table contact does not exists. Creating now.")
		_, err := db.Exec(CREATE_CONTACT_TABLE)
		if err != nil {
			log.Error("Error creating contact table")
			panic(fmt.Sprintf("%v", err))
		} else {
			log.Info("Contact Table successfully created....")
		}
		return true
	}
	return false
}

func (db *DB) FindAllContacts() ([]*Contact, error) {
	rows, err := db.Query("SELECT id, firstname, lastname, company, address1, address2, zipcode, city, country, tag, email, web, birthday FROM contact ORDER BY lastname, firstname")
	if err != nil {
		log.Errorf("Error fetching contact table: %v", err)
		return nil, err
	}
	defer rows.Close()
	contacts := make([]*Contact, 0)
	for rows.Next() {
		contact := new(Contact)
		birthday := mysql.NullTime{}
		err := rows.Scan(
			&contact.Id,
			&contact.Firstname,
			&contact.Lastname,
			&contact.Company,
			&contact.Address1,
			&contact.Address2,
			&contact.Zipcode,
			&contact.City,
			&contact.Country,
			&contact.Tag,
			&contact.Email,
			&contact.Web,
			&birthday,
		)
		if err != nil {
			log.Error(err)
		}
		if birthday.Valid {
			contact.Birthday = &birthday.Time
		}
		contacts = append(contacts, contact)
	}
	if err = rows.Err(); err != nil {
		log.Error(err)
	}
	return contacts, err
}

func (db *DB) FindContactById(id int64) (*Contact, error) {
	contact := Contact{}
	birthday := mysql.NullTime{}
	err := db.QueryRow("SELECT id,firstname, lastname, company, address1, address2, zipcode, city, country, tag, email, web, birthday FROM contact WHERE id=?", id).Scan(
		&contact.Id,
		&contact.Firstname,
		&contact.Lastname,
		&contact.Company,
		&contact.Address1,
		&contact.Address2,
		&contact.Zipcode,
		&contact.City,
		&contact.Country,
		&contact.Tag,
		&contact.Email,
		&contact.Web,
		&birthday,
	)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if birthday.Valid {
		contact.Birthday = &birthday.Time
	}
	contact.Phones, err = db.findPhoneByContact(contact.Id)
	if err != nil {
		return nil, err
	}

	return &contact, err
}

func (db *DB) InsertContact(contact *Contact) (*Contact, error) {
	tx, err := db.Begin()
	result, err := tx.Exec("INSERT INTO contact (firstname, lastname, company, address1, address2, zipcode, city, country, tag, email, web, birthday) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)", contact.Firstname, contact.Lastname, contact.Company, contact.Address1, contact.Address2, contact.Zipcode, contact.City, contact.Country, contact.Tag, contact.Email, contact.Web, contact.Birthday)
	if err != nil {
		log.Errorf("Error inserting contact: %v", err)
		tx.Rollback()
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Errorf("Error contact id: %v", err)
		tx.Rollback()
		return nil, err
	}
	contact.Id = id
	if contact.Phones != nil {
		newTx := TX{*tx}
		phones := make([]*Phone, 0)
		for _, phone := range contact.Phones {
			phone.ContactId = contact.Id
			p, err := newTx.insertPhone(phone)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			phones = append(phones, p)
		}
		contact.Phones = phones
	}
	err = tx.Commit()
	return contact, err
}

func (db *DB) UpdateContact(contact *Contact) (*Contact, error) {
	tx, err := db.Begin()
	_, err = tx.Exec("UPDATE contact SET firstname=?, lastname=?, company=?, address1=?, address2=?, zipcode=?, city=?, country=?, tag=?, email=?, web=?, birthday=? WHERE id=?", contact.Firstname, contact.Lastname, contact.Company, contact.Address1, contact.Address2, contact.Zipcode, contact.City, contact.Country, contact.Tag, contact.Email, contact.Web, contact.Birthday, contact.Id)
	if err != nil {
		log.Errorf("Error updating contact: %v", err)
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec("DELETE FROM phone WHERE contact_id=?", contact.Id)
	if err != nil {
		log.Errorf("Error deleting phones of contact: %v", err)
		tx.Rollback()
		return nil, err
	}
	if contact.Phones != nil {
		newTx := TX{*tx}
		phones := make([]*Phone, 0)
		for _, phone := range contact.Phones {
			phone.ContactId = contact.Id
			p, err := newTx.insertPhone(phone)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			phones = append(phones, p)
		}
		contact.Phones = phones
	}
	err = tx.Commit()
	return contact, err
}

func (db *DB) DeleteContact(id int64) error {
	tx, err := db.Begin()
	_, err = tx.Exec("DELETE FROM phone WHERE contact_id=?", id)
	if err != nil {
		log.Errorf("Error deleting phones of contact: %v", err)
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM contact WHERE id=?", id)
	if err != nil {
		log.Errorf("Error deleting contact: %v", err)
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

/*----------------------------------------------------------------------------------------*/

func (db *DB) initializePhone() {
	_, err := db.Query("SELECT 1 FROM phone LIMIT 1")
	if err != nil {
		log.Info("Table phone does not exists. Creating now.")
		_, err := db.Exec(CREATE_PHONE_TABLE)
		if err != nil {
			log.Error("Error creating phone table")
			panic(fmt.Sprintf("%v", err))
		} else {
			log.Info("Address Table successfully created....")
		}
	}
}

func (db *DB) findPhoneByContact(id int64) ([]*Phone, error) {
	rows, err := db.Query("SELECT id, name, number, contact_id FROM phone WHERE contact_id=?", id)
	if err != nil {
		log.Errorf("Error fetching phone table: %v", err)
		return nil, err
	}
	defer rows.Close()
	phones := make([]*Phone, 0)
	for rows.Next() {
		phone := new(Phone)
		err := rows.Scan(&phone.Id, &phone.Name, &phone.Number, &phone.ContactId)
		if err != nil {
			log.Error(err)
		}
		phones = append(phones, phone)
	}
	if err = rows.Err(); err != nil {
		log.Error(err)
	}
	return phones, err
}

func (txi *TX) insertPhone(phone *Phone) (*Phone, error) {
	result, err := txi.Exec("INSERT INTO phone (name, number, contact_id) VALUES(?,?,?)", phone.Name, phone.Number, phone.ContactId)
	if err != nil {
		log.Errorf("Error inserting phone: %v", err)
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Errorf("Error phone id: %v", err)
		return nil, err
	}
	phone.Id = id
	return phone, nil
}

func (db *DB) insertTestData() {
	f, err := os.Open("us-500.csv")
	if err != nil {
		log.Errorf("Testdata not found (us-500.csv)")
		return
	}
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record[0] != "first_name" {
			contact := Contact{
				Firstname: record[0],
				Lastname:  record[1],
				Company:   record[2],
				Address1:  record[3],
				City:      record[4],
				Country:   record[6],
				Zipcode:   record[7],
				Email:     record[10],
				Web:       record[11],
			}
			phones := make([]*Phone, 0)
			if record[8] != "" {
				phone := Phone{
					Name:   "Telefon 1",
					Number: record[8],
				}
				phones = append(phones, &phone)
			}
			if record[9] != "" {
				phone := Phone{
					Name:   "Telefon 2",
					Number: record[9],
				}
				phones = append(phones, &phone)
			}
			if len(phones) > 0 {
				contact.Phones = phones
			}
			db.InsertContact(&contact)
		}
	}
}
