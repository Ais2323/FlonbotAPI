package KomicaReply

import (
	db "FlonBotApi/Database"
	tool "FlonBotApi/Helper"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// 參考API https://github.com/Nekosyndrome/yonkoma/wiki/Api
var _ALL_KOMICA_API_URL = []string{
	"https://gaia.komica1.org/00b/pixmicat.php?mode=module&load=mod_ajax&action=posts", // 綜合
	"https://grea.komica1.org/64/pixmicat.php?mode=module&load=mod_ajax&action=posts",  // 掛圖
	"https://tomo.komica1.org/12/pixmicat.php?mode=module&load=mod_ajax&action=posts",  // 歡樂惡搞
	"https://tomo.komica1.org/42a/pixmicat.php?mode=module&load=mod_ajax&action=posts", // 四格
	"https://grea.komica1.org/30/pixmicat.php?mode=module&load=mod_ajax&action=posts",  // 塗鴉王國
}

const (
	_TEMP_POST_COUNT   = 10000
	_REQUEST_NEXT_TIME = 5
	_REQUEST_ALL_TIME  = 30
)

type komicaPosts struct {
	Posts []komicaPost `json:"posts"`
}

type komicaPost struct {
	No    int    `json:"no"`
	Resto int    `json:"resto"`
	Com   string `json:"com"`
}
type komicaBoard struct {
	lastNo int
	Posts  []komicaPost `json:"posts"`
}

var komicaAllPosts = []komicaBoard{}
var isStart = false

// 資料取得
func getKomicaDatabase() {
	// 補滿
	for len(komicaAllPosts) < len(_ALL_KOMICA_API_URL) {
		komicaAllPosts = append(komicaAllPosts, komicaBoard{})
	}
	for idx, url := range _ALL_KOMICA_API_URL {
		// 延遲5秒請求下一格版面
		time.Sleep(5 * time.Second)
		borad := &komicaAllPosts[idx]
		request, err := http.Get(url)
		if err != nil {
			// 連線錯誤
			continue
		}
		defer request.Body.Close()
		if request.StatusCode != 200 {
			// StatusCode非200錯誤
			continue
		}
		body, err := io.ReadAll(request.Body)
		if err != nil {
			continue
		}
		posts := komicaPosts{}
		err = json.Unmarshal(body, &posts)
		if err != nil {
			continue
		}
		endIdx := len(posts.Posts)
		maxNo := borad.lastNo
		for idx := 0; idx < endIdx; idx++ {
			posts.Posts[idx].Com = tool.HtmlEscapeSign(posts.Posts[idx].Com, false)
			// 舊資料 和 過濾詞
			if posts.Posts[idx].No <= borad.lastNo || db.HasIgnoreWord(posts.Posts[idx].Com) {
				posts.Posts = removePost(posts.Posts, idx)
				endIdx--
				idx--
				continue
			}
			// 紀錄該次拿到的最新ID
			if posts.Posts[idx].No > maxNo {
				maxNo = posts.Posts[idx].No
			}
			fmt.Println(posts.Posts[idx])
		}
		// 更新最後拿到的No
		borad.lastNo = maxNo
		// 從前面拚
		borad.Posts = append(posts.Posts, borad.Posts...)
		// 移除舊資料
		if len(borad.Posts) > _TEMP_POST_COUNT {
			borad.Posts = borad.Posts[:_TEMP_POST_COUNT]
		}
	}
}
func removePost(s []komicaPost, i int) []komicaPost {
	// 移除 並把最後移到該位置
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func StartRequest() {
	isStart = true
	go func() {
		for isStart {
			getKomicaDatabase()
			time.Sleep(_REQUEST_ALL_TIME * time.Second)
		}
	}()
}

func CloseRequest() {
	isStart = false
}

// K島回復
func GetReplyOnKomica(words []string) string {
	findMatchComment()
	findReplyComment()
	randomReplyComment()
	// TODO: 還沒做 先把他第一個關鍵字丟回去
	return words[0]
}

// TODO: 找符合該訊息的文章 預計輸出 {符合分數(字數 數量) 文章物件}
func findMatchComment() {

}

// TODO: 找回覆該文章的 預計輸出 {符合分數(字數 數量) 回覆字串}
func findReplyComment() {

}

// TODO: 按分數比分輸出回覆 預計輸出 {回覆字串}
func randomReplyComment() {

}
