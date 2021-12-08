package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
)

func badRequest(w http.ResponseWriter, err error) {
	log.Printf("%+v\n", err)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func InternalServerError(w http.ResponseWriter, err error) {
	log.Printf("%+v\n", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func GetIdFromContext(ctx context.Context) (id int64, err error) {
	id, ok := ctx.Value(AuthenticateContextKey).(int64)
	if !ok {
		return 0, errors.New("parsing error")
	}
	return id, nil
}

//func FormatAndSending()  {
//	data, err := json.Marshal(item)
//	if err != nil {
//		InternalServerError(w, err)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	_, err = w.Write(data)
//	if err != nil {
//		InternalServerError(w, err)
//		return
//	}
//}
