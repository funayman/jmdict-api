package install

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"app/shared/database"
)

type Config struct {
	JMDictFile    string `json:"jmdict"`
	KanjiDic2File string `json:"kanjidic2"`
}

//JMDict reads in the JMdict file and inserts the data into the database
func JMDict(config Config) error {
	//get the file
	data, err := os.Open(config.JMDictFile)
	if err != nil {
		return err
	}
	defer data.Close()

	//parse data
	words, err := LoadJMDict(data)
	if err != nil {
		return err
	}

	insertIntoDatabase(words)

	return nil
}

func insertIntoDatabase(words []*Entry) {
	//open sql file
	sqlFile, err := os.Open("./sql/sqlite3_install.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlFile.Close()

	/*******************************************
	 * CREATE TABLES
	 ******************************************/
	//Break up individual queries, Sqlite3 cannot do multi-statements
	bigAssQuery, _ := ioutil.ReadAll(sqlFile)
	queries := strings.Split(string(bigAssQuery), ";\n")
	tx, err := database.SQL.Begin()
	if err != nil {
		log.Fatal(err)
	}

	//Execute queries
	for _, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			log.Println("shit done broke w/ the query: " + query)
			tx.Rollback()
			log.Fatal(err)
		}
	}

	/*******************************************
	 * INSERT JMDict to database
	 ******************************************/
	for _, word := range words {
		kanjiReferenceMap := make(map[string]int64) //used for <re_restr> (rstr)

		/*******************************************
		 * JMDict:    <entry>
		 * Database:  enty
		 ******************************************/
		rslt, err := tx.Exec("INSERT INTO enty (entseq) VALUES (?)", word.EntSeq)
		if err != nil {
			log.Printf("Error inserting into ENTY table: %+v\n", word)
			tx.Rollback()
		}
		entyID, _ := rslt.LastInsertId()

		/*******************************************
		 * JMDict:    <k_ele>
		 * Database:  kanj
		 ******************************************/
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

			/*******************************************
			 * JMDict:    <ke_inf>
			 * Database:  kinf
			 ******************************************/
			for _, ki := range k.KeInf {
				_, err := tx.Exec("INSERT INTO kinf (kid, kw) VALUES (?, ?)", kid, ki)
				if err != nil {
					log.Printf("Error inserting into KINF table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			/*******************************************
			 * JMDict:    <ke_pri>
			 * Database:  kpri
			 ******************************************/
			for _, kp := range k.KePri {
				_, err := tx.Exec("INSERT INTO kpri (kid, kw) VALUES (?, ?)", kid, kp)
				if err != nil {
					log.Printf("Error inserting into KPRI table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}
		}

		/*******************************************
		 * JMDict:    <r_ele>
		 * Database:  rdng
		 ******************************************/
		for _, r := range word.Rele {
			rrslt, err := tx.Exec("INSERT INTO rdng (eid, rval, nokj) VALUES (?, ?, ?)", entyID, r.Reb, r.ReNokanji)
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

			/*******************************************
			 * JMDict:    <re_restr>
			 * Database:  rstr
			 ******************************************/
			for _, restriction := range r.ReRestr {
				kid, ok := kanjiReferenceMap[restriction]
				if !ok {
					log.Printf("Error looking up kanji reference: %+v\n", r)
					tx.Rollback()
					log.Fatal(err)
				}
				tx.Exec("INSERT INTO rstr (kid, rid) VALUES (?, ?)", kid, rid)
			}

			/*******************************************
			 * JMDict:    <re_inf>
			 * Database:  rinf
			 ******************************************/
			for _, ri := range r.ReInf {
				_, err := tx.Exec("INSERT INTO rinf (rid, kw) VALUES (?, ?)", rid, ri)
				if err != nil {
					log.Printf("Error inserting into RINF table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			/*******************************************
			 * JMDict:    <re_pri>
			 * Database:  rpri
			 ******************************************/
			for _, rp := range r.RePri {
				_, err := tx.Exec("INSERT INTO rpri (rid, kw) VALUES (?, ?)", rid, rp)
				if err != nil {
					log.Printf("Error inserting into RPRI table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}
		}

		/*******************************************
		 * JMDict:    <sense>
		 * Database:  sens
		 ******************************************/
		for _, s := range word.Sense {
			srslt, err := tx.Exec("INSERT INTO sens (eid) VALUES (?)", entyID)
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
			/*******************************************
			 * JMDict:    <stagk>
			 * Database:  stagk
			 ******************************************/
			for _, stagk := range s.Stagk {
				_, err := tx.Exec("INSERT INTO stagk (sid, kval) VALUES (?, ?)", sid, stagk)
				if err != nil {
					log.Printf("Error inserting into STAGK table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			/*******************************************
			 * JMDict:    <stagr>
			 * Database:  stagr
			 ******************************************/
			for _, stagr := range s.Stagr {
				_, err := tx.Exec("INSERT INTO stagr (sid, rdng) VALUES (?, ?)", sid, stagr)
				if err != nil {
					log.Printf("Error inserting into STAGR table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			/*******************************************
			 * JMDict:    <pos>
			 * Database:  pos
			 ******************************************/
			for _, pos := range s.Pos {
				_, err := tx.Exec("INSERT INTO pos (sid, kw) VALUES (?, ?)", sid, pos)
				if err != nil {
					log.Printf("Error inserting into POS table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			/*******************************************
			 * JMDict:    <xref>
			 * Database:  xref
			 ******************************************/
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

			/*******************************************
			 * JMDict:    <ant>
			 * Database:  ant
			 ******************************************/
			for _, ant := range s.Ant {
				_, err := tx.Exec("INSERT INTO ant (sid, rdng) VALUES (?, ?)", sid, ant)
				if err != nil {
					log.Printf("Error inserting into ANT table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			/*******************************************
			 * JMDict:    <field>
			 * Database:  field
			 ******************************************/
			for _, field := range s.Field {
				_, err := tx.Exec("INSERT INTO field (sid, ctg) VALUES (?, ?)", sid, field)
				if err != nil {
					log.Printf("Error inserting into FIELD table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			/*******************************************
			 * JMDict:    <misc>
			 * Database:  misc
			 ******************************************/
			for _, misc := range s.Misc {
				_, err := tx.Exec("INSERT INTO misc (sid, text) VALUES (?, ?)", sid, misc)
				if err != nil {
					log.Printf("Error inserting into MISC table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			/*******************************************
			 * JMDict:    <s_inf>
			 * Database:  sinf
			 ******************************************/
			for _, sinf := range s.SInf {
				_, err := tx.Exec("INSERT INTO sinf (sid, text) VALUES (?, ?)", sid, sinf)
				if err != nil {
					log.Printf("Error inserting into SINF table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			/*******************************************
			 * JMDict:    <lsource>
			 * Database:  lsource
			 ******************************************/
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

			/*******************************************
			 * JMDict:    <dial>
			 * Database:  dial
			 ******************************************/
			for _, dial := range s.Dial {
				_, err := tx.Exec("INSERT INTO dial (sid, ben) VALUES (?, ?)", sid, dial)
				if err != nil {
					log.Printf("Error inserting into SINF table: %+v\n", word)
					tx.Rollback()
					log.Fatal(err)
				}
			}

			/*******************************************
			 * JMDict:    <gloss>
			 * Database:  gloss
			 ******************************************/
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
}
