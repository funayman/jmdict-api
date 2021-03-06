package install

import (
	"io/ioutil"
	"os"
	"strings"

	"app/shared/database"
	"app/shared/logger"
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

	//insert into database
	insertWordsIntoDatabase(words)

	return nil
}

//KanjiDic2 reads in the KanjiDic2 file and inserts the data into the database
func KanjiDic2(config Config) error {
	//get the file
	data, err := os.Open(config.KanjiDic2File)
	if err != nil {
		return err
	}
	defer data.Close()

	//parse data
	kanji, err := LoadKanjiDic2(data)
	if err != nil {
		return err
	}

	insertKanjiIntoDatabase(kanji)

	return nil
}

func insertKanjiIntoDatabase(kanji []*Kanji) {
}

func insertWordsIntoDatabase(words []*Entry) {
	//open sql file
	sqlFile, err := os.Open("./sql/sqlite3_install.sql")
	if err != nil {
		logger.Fatal(err)
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
		logger.Fatal(err)
	}

	//Execute queries
	for _, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			tx.Rollback()
			logger.Fatal("shit done broke w/ the query: "+query, err)
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
			tx.Rollback()
			logger.Fatalf("Error inserting into ENTY table: %+v\n%s\n", word, err)
		}
		entyID, _ := rslt.LastInsertId()

		/*******************************************
		 * JMDict:    <k_ele>
		 * Database:  kanj
		 ******************************************/
		for _, k := range word.KEle {
			krslt, err := tx.Exec("INSERT INTO kanj (eid, kval) VALUES (?, ?)", entyID, k.Keb)
			if err != nil {
				logger.Fatalf("Error inserting into KANJ table: %+v\n%s\n", word, err)
				tx.Rollback()
			}

			kid, err := krslt.LastInsertId()
			if err != nil {
				tx.Rollback()
				logger.Fatalf("Error getting last ID from KANJ table: %+v\n%s\n", k, err)
			}
			kanjiReferenceMap[k.Keb] = kid

			/*******************************************
			 * JMDict:    <ke_inf>
			 * Database:  kinf
			 ******************************************/
			for _, ki := range k.KeInf {
				_, err := tx.Exec("INSERT INTO kinf (kid, kw) VALUES (?, ?)", kid, ki)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into KINF table: %+v\n%s\n", word, err)
				}
			}

			/*******************************************
			 * JMDict:    <ke_pri>
			 * Database:  kpri
			 ******************************************/
			for _, kp := range k.KePri {
				_, err := tx.Exec("INSERT INTO kpri (kid, kw) VALUES (?, ?)", kid, kp)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into KPRI table: %+v\n%s\n", word, err)
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
				tx.Rollback()
				logger.Fatalf("Error inserting into RDNG table: %+v\n%s\n", word, err)
			}

			rid, err := rrslt.LastInsertId()
			if err != nil {
				tx.Rollback()
				logger.Fatalf("Error getting last ID from RDNG table: %+v\n%s\n", r, err)
			}

			/*******************************************
			 * JMDict:    <re_restr>
			 * Database:  rstr
			 ******************************************/
			for _, restriction := range r.ReRestr {
				kid, ok := kanjiReferenceMap[restriction]
				if !ok {
					tx.Rollback()
					logger.Fatalf("Error looking up kanji reference: %+v\n%s\n", r, err)
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
					tx.Rollback()
					logger.Fatalf("Error inserting into RINF table: %+v\n%s\n", word, err)
				}
			}

			/*******************************************
			 * JMDict:    <re_pri>
			 * Database:  rpri
			 ******************************************/
			for _, rp := range r.RePri {
				_, err := tx.Exec("INSERT INTO rpri (rid, kw) VALUES (?, ?)", rid, rp)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into RPRI table: %+v\n%s\n", word, err)
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
				tx.Rollback()
				logger.Fatalf("Error inserting into SENS table: %+v\n%s\n", word, err)
			}

			sid, err := srslt.LastInsertId()
			if err != nil {
				tx.Rollback()
				logger.Fatalf("Error getting last ID from SENS table: %+v\n%s\n", s, err)
			}

			//stagk
			/*******************************************
			 * JMDict:    <stagk>
			 * Database:  stagk
			 ******************************************/
			for _, stagk := range s.Stagk {
				_, err := tx.Exec("INSERT INTO stagk (sid, kval) VALUES (?, ?)", sid, stagk)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into STAGK table: %+v\n%s\n", word, err)
				}
			}

			/*******************************************
			 * JMDict:    <stagr>
			 * Database:  stagr
			 ******************************************/
			for _, stagr := range s.Stagr {
				_, err := tx.Exec("INSERT INTO stagr (sid, rdng) VALUES (?, ?)", sid, stagr)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into STAGR table: %+v\n%s\n", word, err)
				}
			}

			/*******************************************
			 * JMDict:    <pos>
			 * Database:  pos
			 ******************************************/
			for _, pos := range s.Pos {
				_, err := tx.Exec("INSERT INTO pos (sid, kw) VALUES (?, ?)", sid, pos)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into POS table: %+v\n%s\n", word, err)
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
					tx.Rollback()
					logger.Fatalf("Error inserting into XREF table: %+v\n%s\n", word, err)
				}
			}

			/*******************************************
			 * JMDict:    <ant>
			 * Database:  ant
			 ******************************************/
			for _, ant := range s.Ant {
				_, err := tx.Exec("INSERT INTO ant (sid, rdng) VALUES (?, ?)", sid, ant)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into ANT table: %+v\n%s\n", word, err)
				}
			}

			/*******************************************
			 * JMDict:    <field>
			 * Database:  field
			 ******************************************/
			for _, field := range s.Field {
				_, err := tx.Exec("INSERT INTO field (sid, ctg) VALUES (?, ?)", sid, field)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into FIELD table: %+v\n%s\n", word, err)
				}
			}

			/*******************************************
			 * JMDict:    <misc>
			 * Database:  misc
			 ******************************************/
			for _, misc := range s.Misc {
				_, err := tx.Exec("INSERT INTO misc (sid, text) VALUES (?, ?)", sid, misc)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into MISC table: %+v\n%s\n", word, err)
				}
			}

			/*******************************************
			 * JMDict:    <s_inf>
			 * Database:  sinf
			 ******************************************/
			for _, sinf := range s.SInf {
				_, err := tx.Exec("INSERT INTO sinf (sid, text) VALUES (?, ?)", sid, sinf)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into SINF table: %+v\n%s\n", word, err)
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
					tx.Rollback()
					logger.Fatalf("Error inserting into LSOURCE table: %+v\n%s\n", word, err)
				}
			}

			/*******************************************
			 * JMDict:    <dial>
			 * Database:  dial
			 ******************************************/
			for _, dial := range s.Dial {
				_, err := tx.Exec("INSERT INTO dial (sid, ben) VALUES (?, ?)", sid, dial)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into SINF table: %+v\n%s\n", word, err)
				}
			}

			/*******************************************
			 * JMDict:    <gloss>
			 * Database:  gloss
			 ******************************************/
			for _, gloss := range s.Gloss {
				_, err := tx.Exec("INSERT INTO gloss (sid, text, lang, gender) VALUES (?, ?, ?, ?)", sid, gloss.Value, gloss.Lang, gloss.Gender)
				if err != nil {
					tx.Rollback()
					logger.Fatalf("Error inserting into LSOURCE table: %+v\n%s\n", word, err)
				}
			}

		}
	}
	tx.Commit()
}
