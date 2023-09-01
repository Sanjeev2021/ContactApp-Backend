package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"contactapp/errors"
)

// UnmarshalJSON unmarshals the json data into the given interface
//unmarshalJSON basically helps in decoding the json data into the given interface

func UnmarshalJSON(request *http.Request, out interface{}) error {
	if request.Body == nil {
		fmt.Println("===========err request.Body == nil=============")
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}
	//The ioutil package contains functions to read from and write to files, as well as functions to work with byte slices (like reading/writing to and from byte slices).
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println("===========err ioutil.ReadALL(request.Body)=============")
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	if len(body) == 0 {
		fmt.Println("===========err len(body) == 0=============")
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		fmt.Println("===========err json.Unmarshal(body ,out)=============")
		fmt.Println(body)
		fmt.Println(out)
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	return nil
}
