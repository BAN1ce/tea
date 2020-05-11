package api

import (
	"encoding/json"
	"net/http"
	"tea/src/distributed"
)

func GetHandle(w http.ResponseWriter, r *http.Request) {

	response := make(map[string]string)

	for k, v := range distributed.Cluster.M {

		response[k] = v.Name
	}
	b, _ := json.Marshal(response)
	w.Write(b)
}
