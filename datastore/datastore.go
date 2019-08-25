package datastore

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// User contains a single user.
type User struct {
	Name     string
	Password string
}

// UserDB emulates a user datastore
type UserDB []User

// NewUserDB initializes an empty UserDB
func NewUserDB() UserDB {
	return UserDB{}
}

// LoadDB loads data into UserDB struct from a csv blob
func (u *UserDB) LoadDB(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not load datastore: %v", err)
	}
	defer f.Close()
	r := csv.NewReader(f)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break // exit for loop
		}
		if err != nil {
			return fmt.Errorf("could not read record from datastore: %v", err)
		}
		*u = append(*u, User{Name: record[0], Password: record[1]})
	}
	return nil
}

// WriteDB writes contents of UserDB to csv blob
func (u *UserDB) WriteDB(filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could open datastore to write to: %v", err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	for _, record := range *u {
		fmt.Println(record)
		if err := w.Write([]string{record.Name, record.Password}); err != nil {
			return fmt.Errorf("could not write record %+v to datastore: %v", record, err)
		}
	}
	// !Important: Write any remaining buffered data to underlying writer
	w.Flush()
	// check if writer has thrown an error while flushing buffer
	if err := w.Error(); err != nil {
		return fmt.Errorf("error flushing buffer: %v", err)
	}
	fmt.Println("Successfully written to datastore.")
	return nil
}

// AddUser adds a user object to UserDB
func (u *UserDB) AddUser(thisUser User) error {
	*u = append(*u, thisUser)
	return nil
}

// IsNameUnique loops through all records in UserDB and checks if a given Name exists
func (u *UserDB) IsNameUnique(thisName string) bool {
	for _, record := range *u {
		if record.Name == thisName {
			fmt.Printf("User name '%s' exists; please type in a different user name.", thisName)
			return false
		}
	}
	return true
}
