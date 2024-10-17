package servercommon

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

func ProcessBodyData(r *http.Request, dto any) error {
	dataType := r.Header.Get(HEADER_CONTENT_TYPE)
	if dataType == APPLICATION_JSON_TYPE {
		err := json.NewDecoder(r.Body).Decode(dto)
		return err
	}
	err := r.ParseMultipartForm(MAX_UPLOAD_FILE_SIZE)
	if err != nil {
		return err
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(dto, r.PostForm)
	return err
}

func SendJsonResponse(w http.ResponseWriter, resp any) error {
	respStr, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error. Can't generate json. Error: " + err.Error())
		return errors.New(ERROR_INTERNAL)
	}
	w.Header().Set(HEADER_CONTENT_TYPE, APPLICATION_JSON_TYPE)
	w.Write(respStr)
	return nil
}
