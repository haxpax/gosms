package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB(driver, dbname string) error {
	var err error
	db, err = sql.Open(driver, dbname)
	return err
}

func insertMessage(sms *SMS) error {

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO messages(uuid, message, mobile) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(sms.uuid, sms.body, sms.mobile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	tx.Commit()
	return nil
}

func updateMessageStatus(sms SMS) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}
	stmt, err := tx.Prepare("UPDATE messages SET status=?")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(sms.status)
	if err != nil {
		fmt.Println(err)
		return err
	}
	tx.Commit()
	return nil

}

func getPendingMessages() ([]SMS, error) {
	query := fmt.Sprintf("SELECT uuid, message, mobile, status FROM messages WHERE status=%v", SMSPending)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var messages []SMS

	for rows.Next() {
		sms := SMS{}
		rows.Scan(&sms.uuid, &sms.body, &sms.mobile, &sms.status)
		messages = append(messages, sms)
	}
	rows.Close()
	return messages, nil
}
