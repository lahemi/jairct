package main

import (
	"database/sql"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SourceTable struct {
	ID                    int64
	Network, Target, Nick string
}

type MsgTable struct {
	ID        int64
	Timestamp int64
	Content   string
	Source    int64 // source table id
}

var DB *sql.DB

func setupDB() {
	var err error // This is so that no local DB variable is created.
	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		die("Failed to open DB connection pool")
	}

	DB.SetMaxIdleConns(100)

	err = DB.Ping()
	if err != nil {
		die("Failed to open DB connection: " + err.Error())
	}

	if _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER NOT NULL PRIMARY KEY,
            timestamp INTEGER NOT NULL, -- UNIX timestamp
            content TEXT,
            source INTEGER,
            FOREIGN KEY(source) REFERENCES sources(source)
        );
        CREATE TABLE IF NOT EXISTS sources (
            id INTEGER NOT NULL PRIMARY KEY,
            network TEXT,
            target TEXT,
            nick TEXT
        );
    `); err != nil {
		die("Failed to execute SQLite3.")
	}
}

func saveSource(fields SourceTable) (int64, error) {
	err := DB.Ping()
	if err != nil {
		die("Failed to open DB connection: " + err.Error())
	}
	stmt, err := DB.Prepare(`
        INSERT INTO sources(network, target, nick)
        VALUES(?, ?, ?)
    `)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(fields.Network, fields.Target, fields.Nick)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func saveMsg(fields MsgTable) error {
	err := DB.Ping()
	if err != nil {
		die("Failed to open DB connection: " + err.Error())
	}
	msgStmt, err := DB.Prepare(`
        INSERT INTO messages(timestamp, content, source)
        VALUES(?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer msgStmt.Close()

	_, err = msgStmt.Exec(fields.Timestamp, fields.Content, fields.Source)
	if err != nil {
		return err
	}
	return nil
}

// MsgLine : nick, cmd, target, msg
func saveDB(msg MsgLine) {
	srcFields := SourceTable{}
	msgFields := MsgTable{}

	srcFields.Network = initNetwork
	srcFields.Target = msg.Target
	srcFields.Nick = msg.Nick

	srcId, err := saveSource(srcFields)
	if err != nil {
		stderr(err)
		return
	}

	msgFields.Timestamp = int64(time.Now().Unix())
	msgFields.Content = msg.Msg
	msgFields.Source = srcId

	err = saveMsg(msgFields)
	if err != nil {
		stderr(err)
		return
	}
}

func getMsgs() []string {
	err := DB.Ping()
	if err != nil {
		die("Failed to open DB connection: " + err.Error())
	}
	rows, err := DB.Query(`
        SELECT messages.timestamp, sources.target, sources.nick, messages.content
        FROM messages INNER JOIN sources ON messages.source = sources.id;
    `)
	if err != nil {
		stderr(err)
		return []string{}
	}
	defer rows.Close()

	var rets []string
	sep := " | "

	for rows.Next() {
		var (
			timestamp             int64
			target, nick, content string
		)
		if err := rows.Scan(&timestamp, &target, &nick, &content); err != nil {
			stderr(err)
			return []string{}
		}
		rets = append(rets, strconv.FormatInt(timestamp, 10)+sep+target+sep+nick+content)
	}
	if err := rows.Err(); err != nil {
		stderr(err)
		return []string{}
	}

	return rets
}
