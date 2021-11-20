package handlers

import (
	"log"
	"net/http"
)

func badRequest(w http.ResponseWriter, err error)  {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func InternalServerError(w http.ResponseWriter, err error)  {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

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
