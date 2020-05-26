package api

import (
	"encoding/json"
	"net/http"
	"tea/src/distributed"
	"tea/src/status"
)

func GetHandle(w http.ResponseWriter, r *http.Request) {

}

func GetBroadcastTotalCount(w http.ResponseWriter, r *http.Request) {

	response := make(map[string]interface{})

	response["total_count"] = distributed.BroadcastTotalCount
	response["total_handle_count"] = distributed.BroadcastHandleTotalCount
	response["total_broadcasted"] = distributed.BroadcastedCount
	response["total_hb_count"] = distributed.HBTotalCount
	response["client_count"] = status.GetClientCount()
	b, _ := json.Marshal(response)

	w.Write(b)

}
