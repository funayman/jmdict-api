package database

import (
	"database/sql"

	"app/shared/logger"
)

const (
	QueryKanjiAndReading = `SELECT (SELECT vk.kanjis FROM vkanjicc vk WHERE vk.entyid=?) AS "kanjis", (SELECT vr.readings FROM vreadingcc vr WHERE vr.entyid=?) AS "readings"`
	QueryGlossPosField   = `SELECT g.gloss, p.pos, f.ctg FROM
		(SELECT s.eid, s.id, group_concat(g.text, "; ") AS "gloss" FROM gloss g
		 INNER JOIN sens s ON s.id = g.sid AND s.eid = ?
		 GROUP BY s.id) AS g
		LEFT JOIN
			(SELECT s.id, group_concat(p.kw, "; ") AS "pos" FROM pos p
			 INNER JOIN sens s ON s.id = p.sid AND s.eid = ?
			 GROUP BY p.sid) AS p ON p.id = g.id
		LEFT JOIN
			(SELECT s.eid, s.id, group_concat(f.ctg, "; ") AS "ctg" FROM field f
			 INNER JOIN sens s ON s.id = f.sid AND s.eid = ?
			 GROUP BY s.id) AS f ON f.id = g.id`
	QuerySearchForID = `SELECT DISTINCT t.eid FROM (SELECT r.eid AS "eid" FROM rdng r WHERE r.rval=? UNION SELECT k.eid AS "eid" FROM kanj k WHERE k.kval=?) AS t`
	ResultDelimeter  = "; "
)

var (
	//SQL is a wrapper for database/sql
	SQL *sql.DB

	//Driver is the database type
	Driver = "sqlite3"

	//Connection to the database
	Connection = "./test.db"
)

//Setup the SQL connection
type Setup struct {
	Driver     string `json:"driver"`
	Connection string `json:"connect"`
}

//Connect to database of choice
func Connect(info Setup) {
	var err error
	SQL, err = sql.Open(info.Driver, info.Connection)
	if err != nil {
		logger.Fatal("SQL Open error: ", err)
	}

	//we good?
	if err = SQL.Ping(); err != nil {
		logger.Fatal("Database connection error: ", err)
	}
}

func Close() {
	logger.Info("closing database connection...")
	if err := SQL.Close(); err != nil {
		logger.Fatal(err)
	}
}
