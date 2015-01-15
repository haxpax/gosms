package gosms

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB

func InitDB(driver, dbname string) (*sql.DB, error) {
	var err error
	db, err = sql.Open(driver, dbname)
	return db, err
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

func getMessages(filter string) ([]SMS, error) {
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
