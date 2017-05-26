package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/funayman/jmdict-api/parser"
	_ "github.com/mattn/go-sqlite3"
)

var (
	installFlag = flag.Bool("install", false, "install to a local db")
)

func main() {
	flag.Parse()
	if *installFlag {
		installStuff()
		os.Exit(0)
	}

	//kanjidic2
	kanjiFile, err := os.Open("./data/kanjidic2.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer kanjiFile.Close()

	kanji, err := parser.LoadKanjiDic2(kanjiFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Len of Kanji: %d\n", len(kanji))
	fmt.Printf("%+v\n", kanji[100])
}

func installStuff() {
	//install JMDict to DB

	//get the file
	data, err := os.Open("./data/JMdict_e")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	//parse data
	words, err := parser.LoadJMDict(data)
	if err != nil {
		log.Fatal(err)
	}

	//open database
	//Open DB
	dbPath := "./test.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	//open sql file
	sqlFile, err := os.Open("sql/install.sql")
	if err != nil {
		log.Fatal(err)
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

	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}

	//insert into db
	wtx, _ := db.Begin()
	for _, word := range words {
		kanjiReferenceMap := make(map[string]int64) //used for re_restr

		//enty
		rslt, err := wtx.Exec("INSERT INTO enty (entseq) VALUES (?)", word.EntSeq)
		if err != nil {
			log.Printf("Error inserting into ENTY table: %+v\n", word)
			wtx.Rollback()
		}
		entyID, _ := rslt.LastInsertId()

		//kanj
		for _, k := range word.KEle {
			krslt, err := wtx.Exec("INSERT INTO kanj (eid, kval) VALUES (?, ?)", entyID, k.Keb)
			if err != nil {
				log.Printf("Error inserting into KANJ table: %+v\n", word)
				wtx.Rollback()
				log.Fatal(err)
			}

			kid, err := krslt.LastInsertId()
			if err != nil {
				log.Printf("Error getting last ID from KANJ table: %+v\n", k)
				wtx.Rollback()
				log.Fatal(err)
			}
			kanjiReferenceMap[k.Keb] = kid

			//kinf
			for _, ki := range k.KeInf {
				_, err := wtx.Exec("INSERT INTO kinf (kid, kw) VALUES (?, ?)", kid, ki)
				if err != nil {
					log.Printf("Error inserting into KINF table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//kpri
			for _, kp := range k.KePri {
				_, err := wtx.Exec("INSERT INTO kpri (kid, kw) VALUES (?, ?)", kid, kp)
				if err != nil {
					log.Printf("Error inserting into KPRE table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}
		}

		//rdng
		for _, r := range word.Rele {
			rrslt, err := wtx.Exec("INSERT INTO rdng (rval, eid, nokj) VALUES (?, ?, ?)", r.Reb, entyID, r.ReNokanji)
			if err != nil {
				log.Printf("Error inserting into RDNG table: %+v\n", word)
				wtx.Rollback()
				log.Fatal(err)
			}

			rid, err := rrslt.LastInsertId()
			if err != nil {
				log.Printf("Error getting last ID from RDNG table: %+v\n", r)
				wtx.Rollback()
				log.Fatal(err)
			}

			//rstr
			for _, restriction := range r.ReRestr {
				kid, ok := kanjiReferenceMap[restriction]
				if !ok {
					log.Printf("Error looking up kanji reference: %+v\n", r)
					wtx.Rollback()
					log.Fatal(err)
				}
				wtx.Exec("INSERT INTO rstr (kid, rid) VALUES (?, ?)", kid, rid)
			}

			//rinf
			for _, ri := range r.ReInf {
				_, err := wtx.Exec("INSERT INTO rinf (rid, kw) VALUES (?, ?)", rid, ri)
				if err != nil {
					log.Printf("Error inserting into RINF table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//rpri
			for _, rp := range r.RePri {
				_, err := wtx.Exec("INSERT INTO rpri (rid, kw) VALUES (?, ?)", rid, rp)
				if err != nil {
					log.Printf("Error inserting into RPRE table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}
		}

		//sens
		for _, s := range word.Sense {
			srslt, err := wtx.Exec("INSERT INTO sens (eid) VALUES (?)", word.EntSeq)
			if err != nil {
				log.Printf("Error inserting into SENS table: %+v\n", word)
				wtx.Rollback()
				log.Fatal(err)
			}

			sid, err := srslt.LastInsertId()
			if err != nil {
				log.Printf("Error getting last ID from SENS table: %+v\n", s)
				wtx.Rollback()
				log.Fatal(err)
			}

			//stagk
			for _, stagk := range s.Stagk {
				_, err := wtx.Exec("INSERT INTO stagk (sid, kval) VALUES (?, ?)", sid, stagk)
				if err != nil {
					log.Printf("Error inserting into STAGK table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//stagr
			for _, stagr := range s.Stagr {
				_, err := wtx.Exec("INSERT INTO stagr (sid, rdng) VALUES (?, ?)", sid, stagr)
				if err != nil {
					log.Printf("Error inserting into STAGR table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//pos
			for _, pos := range s.Pos {
				_, err := wtx.Exec("INSERT INTO pos (sid, kw) VALUES (?, ?)", sid, pos)
				if err != nil {
					log.Printf("Error inserting into POS table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//xref
			for _, xref := range s.Xref {
				// TODO: Fix this issue later
				//split on the "center-dot" char and store only the keb/reb
				ref := strings.Split(xref, "・")

				_, err := wtx.Exec("INSERT INTO xref (sid, rdng) VALUES (?, ?)", sid, ref[0])
				if err != nil {
					log.Printf("Error inserting into XREF table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//ant
			for _, ant := range s.Ant {
				_, err := wtx.Exec("INSERT INTO ant (sid, rdng) VALUES (?, ?)", sid, ant)
				if err != nil {
					log.Printf("Error inserting into ANT table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//field
			for _, field := range s.Field {
				_, err := wtx.Exec("INSERT INTO field (sid, ctg) VALUES (?, ?)", sid, field)
				if err != nil {
					log.Printf("Error inserting into FIELD table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//misc
			for _, misc := range s.Misc {
				_, err := wtx.Exec("INSERT INTO misc (sid, text) VALUES (?, ?)", sid, misc)
				if err != nil {
					log.Printf("Error inserting into MISC table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//sinf
			for _, sinf := range s.SInf {
				_, err := wtx.Exec("INSERT INTO sinf (sid, text) VALUES (?, ?)", sid, sinf)
				if err != nil {
					log.Printf("Error inserting into SINF table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//lsource
			for _, lsource := range s.Lsource {
				var wasei bool
				if "y" == lsource.Wasei {
					wasei = true
				}

				_, err := wtx.Exec("INSERT INTO lsource (sid, text, lang, type, wasei) VALUES (?, ?, ?, ?, ?)", sid, lsource.Value, lsource.Lang, lsource.Type, wasei)
				if err != nil {
					log.Printf("Error inserting into LSOURCE table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//dial
			for _, dial := range s.Dial {
				_, err := wtx.Exec("INSERT INTO dial (sid, ben) VALUES (?, ?)", sid, dial)
				if err != nil {
					log.Printf("Error inserting into SINF table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

			//gloss
			for _, gloss := range s.Gloss {
				_, err := wtx.Exec("INSERT INTO gloss (sid, text, lang, gender) VALUES (?, ?, ?, ?)", sid, gloss.Value, gloss.Lang, gloss.Gender)
				if err != nil {
					log.Printf("Error inserting into LSOURCE table: %+v\n", word)
					wtx.Rollback()
					log.Fatal(err)
				}
			}

		}
	}
	wtx.Commit()
}
