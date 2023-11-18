package APIMethod

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestData struct {
	ErrorMessage string
	Reply        string
}

func RequestReply(w http.ResponseWriter, r *http.Request) {
	result := RequestData{}
	if r.Method != "GET" {
		result.ErrorMessage = "Worng Method"
		returnData, _ := json.Marshal(result)
		w.Write(returnData)
		return
	}
	// if only one expected
	word := r.URL.Query().Get("word")
	if word == "" {
		result.ErrorMessage = "Need Param word example:http://127.0.0.1/reply?word=hello Word"
		returnData, _ := json.Marshal(result)
		w.Write(returnData)
		return
	}
	fmt.Printf("input:%s", word)
	result.Reply = word
	returnData, _ := json.Marshal(result)
	w.Write(returnData)
	return
}
