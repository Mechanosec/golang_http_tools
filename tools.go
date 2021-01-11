package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func JsonDecode(data io.Reader) (map[string]string, error) {
	var responseData map[string]string

	body, _ := ioutil.ReadAll(data)
	err := json.Unmarshal(body, &responseData)

	return responseData, err
}

func JsonDecodeArray(data io.Reader) ([]map[string]interface{}, error) {
	var responseData []map[string]interface{}

	body, _ := ioutil.ReadAll(data)
	err := json.Unmarshal(body, &responseData)

	return responseData, err
}

func CheckRequiredParams(data map[string]string, filter []string) error {
	var missingParams []string
	for _, filterKey := range filter {
		if strings.Contains(filterKey, "|") {
			dependParams := strings.Split(filterKey, "|")
			found := false
			for _, dependParamsKey := range dependParams {
				value, ok := data[dependParamsKey]
				if ok || value != "" {
					found = true
					break
				}
			}
			if !found {
				missingParams = append(missingParams, filterKey)
			}
		} else {
			value, ok := data[filterKey]
			if !ok || value == "" {
				missingParams = append(missingParams, filterKey)
			}
		}
	}

	if len(missingParams) > 0 {
		return errors.New(fmt.Sprintf("Missing %s param", missingParams))
	}

	return nil
}

func ConvertMap(arr map[string][]string) map[string]string {
	data := make(map[string]string)
	if arr != nil {
		for key, val := range arr {
			data[key] = val[0]
		}
	}
	return data
}

func ShowError(err error, message string, w http.ResponseWriter) {
	log.Printf("%s: %s\n", message, err)
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "application/json")

	msg, err := json.Marshal(map[string]string{
		"error": fmt.Sprintf("%s: %s", message, err.Error()),
	})
	if err != nil {
		return
	}
	w.Write(msg)
}

func GetHttpError(data io.Reader) error {
	body, _ := ioutil.ReadAll(data)
	return errors.New(string(body))
}
