package controller

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"app/model"
	"app/shared/database"
	"app/shared/router"
)

var (
	qFormat = "format"
)

func init() {
	router.Route("/word/{word}", GetWordsByChar)
}

func GetWordsByChar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	q := vars["word"]
	format := r.URL.Query().Get(qFormat)
	words := []*model.Word{}

	//grab the base ID for the word(s)
	rows, err := database.SQL.Query(database.QuerySearchForID, q, q, q)
	if err != nil {
		log.Print(err)
		writeToWriter(w, words, format)
	}
	defer rows.Close()

	var id int
	for rows.Next() {
		rows.Scan(&id)
		words = append(words, &model.Word{ID: id})
	}

	for _, word := range words {
		err = word.BuildSelf()
		if err != nil {
			log.Println(err)
		}
	}

	writeToWriter(w, words, format)
}

func writeToWriter(w io.Writer, data interface{}, format string) {
	var output []byte
	var err error

	switch strings.ToLower(format) {
	case "xml":
		output, err = xml.Marshal(data)
		if err != nil {
			log.Println(err)
		}
	default:
		output, err = json.Marshal(data)
		if err != nil {
			log.Println(err)
		}
	}

	fmt.Fprintf(w, "%s\n", output)
}
