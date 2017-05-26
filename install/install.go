package install

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	JMDictFile = "./data/JMdict_e"
	DBDriver   = "sqlite3"
	DBLoc      = "./test.db"
)

func JMDict() error {
	//install JMDict to DB

	//get the file
	data, err := os.Open(JMDictFile)
	if err != nil {
		return err
	}
	defer data.Close()

	//parse data
	words, err := LoadJMDict(data)
	if err != nil {
		return err
	}

	//open database
	//Open DB
	db, err := sql.Open("sqlite3", DBLoc)
	if err != nil {
		return err
	}
	defer db.Close()

	//open sql file
	sqlFile, err := os.Open("./sql/" + DBDriver + "_install.sql")
	if err != nil {
		return err
	}
	defer sqlFile.Close()

	bigAssQuery, _ := ioutil.ReadAll(sqlFile)
	queries := strings.Split(string(bigAssQuery), ";\n")
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			log.Println("shit done broke w/ the query: " + query)
			tx.Rollback()
			log.Fatal(err)
		}
	}

	//insert into db
	for _, word := range words {
		kanjiReferenceMap := make(map[string]int64) //used for re_restr

		//enty
		rslt, err := tx.Exec("INSERT INTO enty (entseq) VALUES (?)", word.EntSeq)
		if err != nil {
			log.Printf("Error inserting into ENTY table: %+v\n", word)
			tx.Rollback()
		}
		entyID, _ := rslt.LastInsertId()

		//kanj
		for _, k := range word.KEle {
			krslt, err := tx.Exec("INSERT INTO kanj (eid, kval) VALUES (?, ?)", entyID, k.Keb)
			if err != nil {
				log.Printf("Error inserting into KANJ table: %+v\n", word)
				tx.Rollback()
				log.Fatal(err)
			}

			kid, err := krslt.LastInsertId()
			if err != nil {
				log.Printf("Error getting last ID from KANJ table: %+v\n", k)
				tx.Rollback()
				log.Fatal(err)
			}
			kanjiReferenceMap[k.Keb] = kid

			//kinf
			for _, ki := range k.KeInf {
				_, err := tx.Exec("INSERT INTO kinf (kid, kw) VALUES (?, ?)", kid, ki)
				if err != nil {
					log.Printf("Error inserting into KINF table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//kpri
			for _, kp := range k.KePri {
				_, err := tx.Exec("INSERT INTO kpri (kid, kw) VALUES (?, ?)", kid, kp)
				if err != nil {
					log.Printf("Error inserting into KPRE table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}
		}

		//rdng
		for _, r := range word.Rele {
			rrslt, err := tx.Exec("INSERT INTO rdng (rval, eid, nokj) VALUES (?, ?, ?)", r.Reb, entyID, r.ReNokanji)
			if err != nil {
				log.Printf("Error inserting into RDNG table: %+v\n", word)
				tx.Rollback()
				log.Fatal(err)
			}

			rid, err := rrslt.LastInsertId()
			if err != nil {
				log.Printf("Error getting last ID from RDNG table: %+v\n", r)
				tx.Rollback()
				log.Fatal(err)
			}

			//rstr
			for _, restriction := range r.ReRestr {
				kid, ok := kanjiReferenceMap[restriction]
				if !ok {
					log.Printf("Error looking up kanji reference: %+v\n", r)
					tx.Rollback()
					log.Fatal(err)
				}
				tx.Exec("INSERT INTO rstr (kid, rid) VALUES (?, ?)", kid, rid)
			}

			//rinf
			for _, ri := range r.ReInf {
				_, err := tx.Exec("INSERT INTO rinf (rid, kw) VALUES (?, ?)", rid, ri)
				if err != nil {
					log.Printf("Error inserting into RINF table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//rpri
			for _, rp := range r.RePri {
				_, err := tx.Exec("INSERT INTO rpri (rid, kw) VALUES (?, ?)", rid, rp)
				if err != nil {
					log.Printf("Error inserting into RPRE table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}
		}

		//sens
		for _, s := range word.Sense {
			srslt, err := tx.Exec("INSERT INTO sens (eid) VALUES (?)", word.EntSeq)
			if err != nil {
				log.Printf("Error inserting into SENS table: %+v\n", word)
				tx.Rollback()
				log.Fatal(err)
			}

			sid, err := srslt.LastInsertId()
			if err != nil {
				log.Printf("Error getting last ID from SENS table: %+v\n", s)
				tx.Rollback()
				log.Fatal(err)
			}

			//stagk
			for _, stagk := range s.Stagk {
				_, err := tx.Exec("INSERT INTO stagk (sid, kval) VALUES (?, ?)", sid, stagk)
				if err != nil {
					log.Printf("Error inserting into STAGK table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//stagr
			for _, stagr := range s.Stagr {
				_, err := tx.Exec("INSERT INTO stagr (sid, rdng) VALUES (?, ?)", sid, stagr)
				if err != nil {
					log.Printf("Error inserting into STAGR table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//pos
			for _, pos := range s.Pos {
				_, err := tx.Exec("INSERT INTO pos (sid, kw) VALUES (?, ?)", sid, pos)
				if err != nil {
					log.Printf("Error inserting into POS table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//xref
			for _, xref := range s.Xref {
				// TODO: Fix this issue later
				//split on the "center-dot" char and store only the keb/reb
				ref := strings.Split(xref, "・")

				_, err := tx.Exec("INSERT INTO xref (sid, rdng) VALUES (?, ?)", sid, ref[0])
				if err != nil {
					log.Printf("Error inserting into XREF table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//ant
			for _, ant := range s.Ant {
				_, err := tx.Exec("INSERT INTO ant (sid, rdng) VALUES (?, ?)", sid, ant)
				if err != nil {
					log.Printf("Error inserting into ANT table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//field
			for _, field := range s.Field {
				_, err := tx.Exec("INSERT INTO field (sid, ctg) VALUES (?, ?)", sid, field)
				if err != nil {
					log.Printf("Error inserting into FIELD table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//misc
			for _, misc := range s.Misc {
				_, err := tx.Exec("INSERT INTO misc (sid, text) VALUES (?, ?)", sid, misc)
				if err != nil {
					log.Printf("Error inserting into MISC table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//sinf
			for _, sinf := range s.SInf {
				_, err := tx.Exec("INSERT INTO sinf (sid, text) VALUES (?, ?)", sid, sinf)
				if err != nil {
					log.Printf("Error inserting into SINF table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//lsource
			for _, lsource := range s.Lsource {
				var wasei bool
				if "y" == lsource.Wasei {
					wasei = true
				}

				_, err := tx.Exec("INSERT INTO lsource (sid, text, lang, type, wasei) VALUES (?, ?, ?, ?, ?)", sid, lsource.Value, lsource.Lang, lsource.Type, wasei)
				if err != nil {
					log.Printf("Error inserting into LSOURCE table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//dial
			for _, dial := range s.Dial {
				_, err := tx.Exec("INSERT INTO dial (sid, ben) VALUES (?, ?)", sid, dial)
				if err != nil {
					log.Printf("Error inserting into SINF table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			//gloss
			for _, gloss := range s.Gloss {
				_, err := tx.Exec("INSERT INTO gloss (sid, text, lang, gender) VALUES (?, ?, ?, ?)", sid, gloss.Value, gloss.Lang, gloss.Gender)
				if err != nil {
					log.Printf("Error inserting into LSOURCE table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

		}
	}
	tx.Commit()

	return nil
}
