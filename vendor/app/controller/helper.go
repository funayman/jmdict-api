package controller

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"strings"

	"app/shared/logger"
)

func writeToWriter(w io.Writer, data interface{}, format string) {
	var err error

	switch strings.ToLower(format) {
	case "xml":
		err = xml.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error(err)
		}
	default:
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error(err)
		}
	}
}
