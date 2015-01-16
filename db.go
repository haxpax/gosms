package gosms

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var db *sql.DB

func InitDB(driver, dbname string) (*sql.DB, error) {
	var err error
	createDb := false
	if _, err := os.Stat(dbname); os.IsNotExist(err) {
		log.Printf("database does not exist %s", dbname)
		createDb = true
	}
	db, err = sql.Open(driver, dbname)
	if createDb {
		if err = syncDB(); err != nil {
			return nil, errors.New("Error creating database")
		}
	}
	return db, nil
}

func syncDB() error {
	log.Println("--- syncDB")
	//create messages table
	createMessages := `CREATE TABLE messages (
                id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
                uuid char(32) UNIQUE NOT NULL,
                message char(160)   NOT NULL,
                mobile   char(15)    NOT NULL,
                status  INTEGER DEFAULT 0,
                retries INTEGER DEFAULT 0
            );`
	_, err := db.Exec(createMessages, nil)
	return err
}

func insertMessage(sms *SMS) error {
	log.Println("--- insertMessage ", sms)
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO messages(uuid, message, mobile) VALUES(?, ?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(sms.UUID, sms.Body, sms.Mobile)
	if err != nil {
		log.Println(err)
		return err
	}
	tx.Commit()
	return nil
}

func updateMessageStatus(sms SMS) error {
	log.Println("--- updateMessageStatus ", sms)
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	stmt, err := tx.Prepare("UPDATE messages SET status=? WHERE uuid=?")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(sms.Status, sms.UUID)
	if err != nil {
		log.Println(err)
		return err
	}
	tx.Commit()
	return nil
}

func getPendingMessages() ([]SMS, error) {
	log.Println("--- getPendingMessages ")
	query := fmt.Sprintf("SELECT uuid, message, mobile, status FROM messages WHERE status=%v", SMSPending)
	log.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var messages []SMS

	for rows.Next() {
		sms := SMS{}
		rows.Scan(&sms.UUID, &sms.Body, &sms.Mobile, &sms.Status)
		messages = append(messages, sms)
	}
	rows.Close()
	return messages, nil
}

func GetMessages(filter string) ([]SMS, error) {
	/*
	   expecting filter as empty string or WHERE clauses,
	   simply append it to the query to get desired set out of database
	*/
	log.Println("--- getPendingMessages ")
	query := fmt.Sprintf("SELECT uuid, message, mobile, status FROM messages %v", filter)
	log.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var messages []SMS

	for rows.Next() {
		sms := SMS{}
		rows.Scan(&sms.UUID, &sms.Body, &sms.Mobile, &sms.Status)
		messages = append(messages, sms)
	}
	rows.Close()
	return messages, nil
}
