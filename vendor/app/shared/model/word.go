package model

import (
	"app/shared/database"
	"encoding/xml"
	"errors"
	"log"
)

/********************
 * Case Studies
 * - 頼む（たのむ）
 * - 明かん（あかん）
 * - 正座（せいざ）
 ********************/

type Word struct {
	XMLName    xml.Name `json:"-" xml:"word"`
	ID         int      `json:"id" xml:"id"`
	Kanji      string   `json:"kanji,omitempty" xml:"kanji,omitempty"`
	Reading    string   `json:"reading" xml:"reading"`
	OtherForms []string `json:"otherForms,omitempty" xml:"otherForms>reading,omitempty"`
}

func (w *Word) BuildSelf() error {
	//make sure ID has been set
	if w.ID == 0 {
		return errors.New("ID cannot be 0")
	}

	var isFirst bool = true

	//see if there is a kanji for this word
	var count int
	err := database.SQL.QueryRow("SELECT COUNT(kanj.kval) FROM kanj WHERE kanj.eid = ?", w.ID).Scan(&count)
	if err != nil {
		log.Print(err)
	}

	if count > 0 {
		//we got a kanji!
		rows, err := database.SQL.Query("SELECT kanj.kval FROM kanj WHERE kanj.eid = ?", w.ID)
		if err != nil {
			log.Print(err)
		}

		//if theres more than one, store the first one as the main reading and others as OtherForms
		isFirst = true
		for rows.Next() {
			if isFirst {
				rows.Scan(&w.Kanji)
				isFirst = false
			} else {
				var k string
				rows.Scan(&k)
				w.OtherForms = append(w.OtherForms, k)
			}
		}

	}

	//there will always be a reading
	isFirst = true
	rows, err := database.SQL.Query("SELECT rdng.rval FROM rdng WHERE rdng.eid = ? ORDER BY rdng.id;", w.ID)
	for rows.Next() {
		if isFirst {
			rows.Scan(&w.Reading)
			isFirst = false
		} else {
			var r string
			rows.Scan(&r)
			w.OtherForms = append(w.OtherForms, r)
		}
	}

	return nil
}
