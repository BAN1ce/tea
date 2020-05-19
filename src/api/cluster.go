package api

import (
	"encoding/json"
	"net/http"
	"tea/src/distributed"
)

func GetHandle(w http.ResponseWriter, r *http.Request) {

}

func GetBroadcastTotalCount(w http.ResponseWriter, r *http.Request) {

	response := make(map[string]interface{})

	response["total_count"] = distributed.BroadcastTotalCount
	response["total_handle_count"] = distributed.BroadcastHandleTotalCount
	response["total_broadcasted"] = distributed.BroadcastedCount

	b, _ := json.Marshal(response)

	w.Write(b)

}
