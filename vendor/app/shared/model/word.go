package model

import (
	"app/shared/database"
	"database/sql"
	"encoding/xml"
	"errors"
	"strings"
)

var (
	emptyString      = ""
	emptyStringArray = []string{}
)

/********************
 * Case Studies
 * - 頼む（たのむ）
 * - 明かん（あかん）
 * - 正座（せいざ）
 ********************/

type Word struct {
	XMLName      xml.Name   `json:"-" xml:"word"`
	ID           int        `json:"id" xml:"id"`
	JmdictEntSeq int        `json:"jmdictEntSeq" xml:"jmdictEntSeq"`
	Kanji        string     `json:"kanji,omitempty" xml:"kanji,omitempty"`
	Reading      string     `json:"reading" xml:"reading"`
	Meanings     []*Meaning `json:"meaning" xml:"meaning"`
	OtherForms   []string   `json:"otherForms,omitempty" xml:"otherForms>reading,omitempty"`
}

type Meaning struct {
	Definition   string   `json:"def" xml:"def"`
	PartOfSpeech []string `json:"pos,omitempty" xml:"pos,omitempty"`
	Field        []string `json:"ctg,omitempty" xml:"ctg,omitempty"`
}

func (w *Word) BuildSelf() error {
	//make sure ID has been set
	if w.ID == 0 {
		return errors.New("ID cannot be 0")
	}

	//get entrysequence along with kanji and reading elements
	var kanjis, readings sql.NullString
	err := database.SQL.QueryRow(database.QueryKanjiAndReading, w.ID, w.ID).Scan(&kanjis, &readings)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	//if we have more than one kanji or reading, split and assign properly
	kanjiSlice := splitIntoArray(kanjis.String)
	readingSlice := splitIntoArray(readings.String)

	if len(kanjiSlice) > 0 {
		w.Kanji = kanjiSlice[0]
		w.OtherForms = append(w.OtherForms, kanjiSlice[1:]...)
	}
	w.Reading = readingSlice[0]
	w.OtherForms = append(w.OtherForms, readingSlice[1:]...)

	//query for meanings
	rows, err := database.SQL.Query(database.QueryGlossPosField, w.ID, w.ID, w.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	defer rows.Close()

	//put each meaning into the struct
	for rows.Next() {
		var def string
		var pos, ctg sql.NullString
		err = rows.Scan(&def, &pos, &ctg)
		if err != nil {
			return err
		}

		m := Meaning{
			Definition:   def,
			PartOfSpeech: splitIntoArray(pos.String),
			Field:        splitIntoArray(ctg.String),
		}

		w.Meanings = append(w.Meanings, &m)
	}
	return nil
}

func splitIntoArray(s string) []string {
	if s == emptyString {
		return emptyStringArray
	}

	return strings.Split(s, database.ResultDelimeter)
}
