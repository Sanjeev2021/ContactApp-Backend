package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"contactapp/errors"
)

// RespondJSON Make response with json formate.
func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(response))
}

// RespondError check error type and Write to ResponseWriter
func RespondError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case *errors.ValidationError:
		RespondJSON(w, http.StatusBadRequest, err)
	case *errors.HTTPError:
		httpError := err.(*errors.HTTPError)
		RespondJSON(w, httpError.HTTPStatus, httpError.ErrorKey)
	default:
		RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
	}
}

// RespondErrorMessage Make error response with Payload.
func RespondErrorMessage(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}

// RespondJSONWithXTotalCount Make response with json formate.
func RespondJSONWithXTotalCount(w http.ResponseWriter, code int, totalCount int, payload interface{}) {
	reponse, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	SetNewHeader(w, "X-Total-Count", strconv.Itoa(totalCount))
	w.WriteHeader(code)
	w.Write([]byte(reponse))
}

// SetNewHeader set new header to response
func SetNewHeader(w http.ResponseWriter, headerName, value string) {
	w.Header().Add("Access-control-expose-headers", headerName)
	w.Header().Set(headerName, value)
}
