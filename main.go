package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	db "github.com/zeddee/password-hashing-demo/datastore"
)

func main() {
	// Initializes a new instance of the user database
	u := db.NewUserDB()
	const datastore = "users.csv"
	if err := u.LoadDB(datastore); err != nil {
		log.Fatal(err)
	}
	for { // Runs forever until ctrl+c or we 'return' from the loop
		fmt.Printf("Contents of datastore:\n%+v\n", u)

		fmt.Printf(
			"Hello! What do you want to do?:\n" + "" +
				"1. Register.\n" +
				"2. Register and hash my password.\n" +
				"3. Show user database.\n" +
				"4. Sign in.\n" +
				"5. Quit.\n")
		// Reads a single UTF-8 character (rune)
		// from STDIN and switches to case.
		switch ask() {
		case "1":
			fmt.Println("Registering new account!")
			fmt.Println("Enter a new user name:")
			newUserName := ask()
			// if name does not already exist, continue with the rest of user registration
			isUnique := u.IsNameUnique(newUserName)
			if isUnique == true {
				fmt.Println("Enter a new password:")
				newPassword := ask()
				thisUser := db.User{Name: newUserName, Password: newPassword}
				fmt.Printf("Your new user is: %s\n Your password is: %s\n", thisUser.Name, thisUser.Password)
				if err := u.AddUser(thisUser); err != nil {
					log.Fatal(err)
				}
				if err := u.WriteDB(datastore); err != nil {
					log.Fatal(err)
				} // Save new account information; don't need to reload database because we're still persisting the UserDB struct
			}
			break // 'breaks' out of the switch, but we're still in the 'for' loop.
		case "2":
			fmt.Println("Registering new account!")
			fmt.Println("Enter a new user name:")
			newUserName := ask()
			fmt.Println("Enter a new user password:")
			newPassword := sha256.Sum256([]byte(ask()))
			// store password as utf-8 encoded hex because we're working with a text storage format.
			thisUser := db.User{Name: newUserName, Password: hex.EncodeToString(newPassword[:])}
			fmt.Printf("Your new user is: %s\n Your password is: %s\n", thisUser.Name, thisUser.Password)
			if err := u.AddUser(thisUser); err != nil {
				log.Fatal(err)
			}
			if err := u.WriteDB(datastore); err != nil {
				log.Fatal(err)
			} // Save new account information; don't need to reload database because we're still persisting the UserDB struct
			break // 'breaks' out of the switch, but we're still in the 'for' loop.
		case "3":
			fmt.Println("Showing contents of user database:")
			for _, record := range u {
				fmt.Printf("User: %s\tPassword: %s\n", record.Name, record.Password)
			}
			break
		case "4":
			fmt.Println("Signing in!")
			break
		case "5":
			fmt.Println("Bye!")
			return // quits the loop
		default:
			fmt.Println("Invalid option. Please try again.")
			break
		}
	}
}

func ask() string {
	reader := bufio.NewReader(os.Stdin)
	inputVal, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("invalid option: %v\n", err)
		return ""
	}

	output := strings.TrimSuffix(inputVal, "\n") // Important!
	return output
}
