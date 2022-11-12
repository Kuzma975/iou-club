package database

import (
	"database/sql"
	"fmt"
	"kuzma975/iou-club/database/logging"
	"os"

	_ "github.com/mattn/go-sqlite3"
	// "encoding/json"
	// "errors"
	// "net/http"
	// "fmt"
	// "io/ioutil"
	// "log"
	// "time"
)

const (
	DatabaseFile string = "./test.db"
	LogFile      string = "./testSqlite.log"
)

func InitializeDatabase() *sql.DB {
	if _, err := os.Stat(DatabaseFile); err != nil {
		if os.IsNotExist(err) {
			logging.Warning.Printf("File %s not exist\n", DatabaseFile)
			logging.Info.Printf("We create %s file for You", DatabaseFile)
			file, err := os.Create(DatabaseFile)
			logging.CheckErr(err)
			file.Close()
		} else {
			logging.CheckErr(err)
		}
	}

	logging.Info.Printf("Database file (%s) exists and open\n", DatabaseFile)
	db, err := sql.Open("sqlite3", DatabaseFile)
	logging.CheckErr(err)
	InitializeTable(db, "user_info")
	return db
}

func InitializeTable(db *sql.DB, tableName string) {
	var create string
	switch tableName {
	case "user_info":
		logging.Debug.Printf("Initializing %s table name", tableName)
		create = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS user_info (
											id			INTEGER PRIMARY KEY AUTOINCREMENT,
											user_id		INTEGER UNIQUE		NULL,
											user_name	TEXT				NULL,
											first_name	TEXT				NULL,
											last_name	TEXT				NULL,
											joined		DATE  				NULL
											)`)

	}
	stmt, err := db.Prepare(create)
	logging.CheckErr(err)

	_, err = stmt.Exec()
	logging.CheckErr(err)
	err = stmt.Close()
	logging.CheckErr(err)
	logging.Info.Printf("Table '%s' is created\n", tableName)
}

func CreateUser(db *sql.DB, userId int64, userName, firstName, lastName string, date int) int64 {
	InitializeTable(db, "user_info")
	insert := fmt.Sprintf("INSERT OR REPLACE INTO user_info(user_id, user_name, first_name, last_name, joined) VALUES (?, ?, ?, ?, ?)")
	stmt, err := db.Prepare(insert)
	logging.CheckErr(err)
	res, err := stmt.Exec(userId, userName, firstName, lastName, date)
	logging.CheckErr(err)
	rows, err := res.RowsAffected()
	err = stmt.Close()
	logging.CheckErr(err)
	return rows
}

// func GetUpdates() (telegram.Result, error) {
// 	client := http.Client{}
// 	responce, err := client.Get("https://api.telegram.org/token/getUpdates")
// 	logging.CheckErr(err)
// 	defer func() {
// 		err = responce.Body.Close()
// 		logging.CheckErr(err)
// 	}()
// 	bodyBytes, err := ioutil.ReadAll(responce.Body)
// 	logging.CheckErr(err)
// 	var resp telegram.Result
// 	err = json.Unmarshal(bodyBytes, &resp)
// 	logging.CheckErr(err)
// 	if resp.Ok {
// 		return resp, nil
// 	} else {
// 		return telegram.Result{}, errors.New("something went wrong with getting update from telegram")
// 	}
// }

func TestDatabase() {
	logFile := logging.InitializeLogging(LogFile)
	defer func() {
		err := logFile.Close()
		logging.CheckErr(err)
		logging.Info.Printf("Logfile %s is closed", LogFile)
	}()

	db := InitializeDatabase()
	defer func() {
		err := db.Close()
		logging.CheckErr(err)
		logging.Info.Printf("Database %s is closed\n", DatabaseFile)
	}()
}

// func Execute() int {
// 	logFile := logging.InitializeLogging()
// 	defer func() {
// 		err := logFile.Close()
// 		logging.CheckErr(err)
// 		logging.Info.Printf("Logfile %s is closed", LogFile)
// 	}()

// 	db := InitializeDatabase()
// 	defer func() {
// 		err := db.Close()
// 		logging.CheckErr(err)
// 		logging.Info.Printf("Database %s is closed\n", DatabaseFile)
// 	}()
// 	response, err := GetUpdates()
// 	logging.CheckErr(err)
// 	for _, update := range response.Updates {
// 		if update.Message.MessageId != 0 {
// 			if update.Message.Text != "" {
// 				WriteToDatabase(db, &update.Message, update.UpdateId)
// 			}
// 		} else {
// 			logging.Debug.Printf("This is not regular message: %+v\n", update)
// 		}
// 	}
// 	result, err := db.Query("PRAGMA foreign_keys;")
// 	logging.CheckErr(err)
// 	resp := 5
// 	for result.Next() {
// 		err = result.Scan(&resp)
// 		logging.CheckErr(err)
// 		if resp != 5 {
// 			logging.Info.Println(resp)
// 		} else {
// 			logging.Warning.Println("foreign key not supported")
// 		}
// 	}

// 	return 0
// }
