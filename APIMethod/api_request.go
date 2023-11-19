package APIMethod

import (
	tool "FlonBotApi/Helper"
	komica "FlonBotApi/KomicaReply"
	"encoding/json"
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
		result.ErrorMessage = "Need Param word example:http://127.0.0.1/reply?word=hello%20Word"
		returnData, _ := json.Marshal(result)
		w.Write(returnData)
		return
	}
	splitWord := tool.SpliteWord(word)
	replyOnKomica := komica.GetReplyOnKomica(splitWord)
	result.Reply = replyOnKomica
	returnData, _ := json.Marshal(result)
	w.Write(returnData)
	return
}
