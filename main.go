//id geven van wachtwoord om op te halen

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	username string
	password string
	hostname string
	dbname   string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env file:", err)
	}

	username = os.Getenv("NAME")
	password = os.Getenv("PASSWORD")
	hostname = os.Getenv("HOSTNAME")
	dbname = os.Getenv("DBNAME")
}

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

func checkPasswordExists(db *sql.DB, password string) bool {
	var exists bool
	// Query to check if the password already exists in the database
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM passwords WHERE password = ? LIMIT 1)", password).Scan(&exists)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return exists
}

func insertPassword(db *sql.DB, password string) {
	// Insert the new password into the database
	_, err := db.Exec("INSERT INTO passwords (password) VALUES (?)", password)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Password inserted successfully!")
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numbers = []rune("0123456789")
var characters = []rune("?><}{}[]!@#$%&*()")

func randSeq(n int, includeNumbers bool, includeSpecialChars bool) string {
	var charSet []rune
	charSet = append(charSet, letters...)

	if includeNumbers {
		charSet = append(charSet, numbers...)
	}

	if includeSpecialChars {
		charSet = append(charSet, characters...)
	}

	p := make([]rune, n)
	for i := range p {
		p[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(p)
}

func main() {
	// Open connection to MySQL database
	db, err := sql.Open("mysql", dsn())
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}
	defer db.Close()

	// Set up connection timeout and ping the database to check connection
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	err = db.Ping()
	if err != nil {
		log.Printf("Error %s pinging DB", err)
		return
	}

	// Log success if connection is successful
	log.Println("Connected to DB successfully")

	var use_number bool
	var use_character bool
	fmt.Print("Use numbers in the pasword? (awnser true or false): ")
	fmt.Scan(&use_number)
	fmt.Print("Use characters in the pasword? (awnser true or false): ")
	fmt.Scan(&use_character)

	password := randSeq(20, use_number, use_character) // Generate password with letters, numbers, and special characters

	// Check if the password already exists in the database
	for checkPasswordExists(db, password) {
		log.Println("Password already exists, generating a new one...")
		password = randSeq(20, use_number, use_character)
	}

	// Once we have a unique password, insert it into the database
	insertPassword(db, password)

	// Output the generated password
	fmt.Println("Generated Password:", password)
}
