package controller

import (
	"net/http"

	"app/model"
	"app/shared/database"
	"app/shared/logger"
	"app/shared/router"
)

var (
	qFormat = "format"
)

func init() {
	router.Route("/word/{word}", GetWordsByChar)
}

func GetWordsByChar(w http.ResponseWriter, r *http.Request) {
	vars := router.GetParams(r)
	q := vars["word"]
	format := r.URL.Query().Get(qFormat)
	words := []*model.Word{}

	//grab the base ID for the word(s)
	rows, err := database.SQL.Query(database.QuerySearchForID, q, q, q)
	if err != nil {
		logger.Error(err)
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
			logger.Error(err)
		}
	}

	writeToWriter(w, words, format)
}
