package KomicaReply

import (
	db "FlonBotApi/Database"
	tool "FlonBotApi/Helper"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
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
	_REPLY_REGEXP      = "(\\>\\>|\\>No\\.)"
)

type komicaPosts struct {
	Posts []KomicaPost `json:"posts"`
}

type KomicaPost struct {
	No    int    `json:"no"`
	Resto int    `json:"resto"`
	Com   string `json:"com"`
}
type komicaBoard struct {
	lastNo int
	Posts  []KomicaPost `json:"posts"`
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
func removePost(s []KomicaPost, i int) []KomicaPost {
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
	result := []ReplyData{}
	for _, board := range komicaAllPosts {
		boardResult := findMatchComment(words, &board)
		boardResult = findReplyComment(boardResult, &board)
		if len(boardResult) > 0 {
			result = append(result, boardResult...)
		}
	}
	// TODO: 先隨機丟回去
	reply := ""
	if len(result) > 0 {
		randSeed := rand.NewSource(time.Now().UnixNano())
		randData := rand.New(randSeed)
		startIdx := randData.Intn(len(result))
		for i := 0; i < len(result); i++ {
			post := result[startIdx]
			if len(post.ReplyString) > 0 {
				targetIdx := randData.Intn(len(post.ReplyString))
				reply = post.ReplyString[targetIdx]
				break
			}
			startIdx++
			if startIdx >= len(result) {
				startIdx = 0
			}
		}
	}
	return reply
}

type ReplyData struct {
	Keywords    []string     // 對應回覆關鍵字
	Post        KomicaPost   // 符合關鍵字文章物件
	ReplyPost   []KomicaPost // 對應回覆文章物件
	ReplyString []string     // 拆出回覆結果
}

// 找符合該訊息的文章
func findMatchComment(words []string, targetBoard *komicaBoard) []ReplyData {
	posts := targetBoard.Posts
	var result = []ReplyData{}
	for _, post := range posts {
		nowPost := ReplyData{}
		nowPost.Post = post
		com := post.Com // TODO: 需要排除不必要資訊
		for _, keyword := range words {
			// 檢查是否有比較多文字符合的關鍵字包含在裡面
			isFind := false
			for _, findKeyword := range nowPost.Keywords {
				isFind = strings.Index(findKeyword, keyword) >= 0
				if isFind {
					break
				}
			}
			if isFind {
				continue
			}
			// 檢查是否包含
			isSuccess := strings.Index(com, keyword) >= 0
			if isSuccess == false {
				continue
			}
			nowPost.Keywords = append(nowPost.Keywords, keyword)
		}
		// 有成功符合關鍵字取出
		if len(nowPost.Keywords) > 0 {
			result = append(result, nowPost)
		}
	}
	return result
}

// TODO: 找回覆該文章的 預計輸出 {符合分數(字數 數量) 回覆字串}
func findReplyComment(result []ReplyData, targetBoard *komicaBoard) []ReplyData {
	posts := targetBoard.Posts
	for _, post := range posts {
		com := post.Com // TODO: 需要排除不必要資訊
		for idx, replyData := range result {
			// 找>>123456
			reOnDirReply := fmt.Sprintf("%s%d", _REPLY_REGEXP, replyData.Post.No)
			isReplyOnComment, err := regexp.MatchString(reOnDirReply, com)
			if err == nil && isReplyOnComment {
				result[idx].ReplyPost = append(result[idx].ReplyPost, post)
				result[idx].ReplyString = append(result[idx].ReplyString, replyCommentSplit(com, reOnDirReply))
				continue
			}
			// 找>>關鍵字
			reReply := fmt.Sprintf("%s.*%s.*\\n", "(\\>\\>)", replyData.Keywords[0])
			if len(replyData.Keywords) > 1 {
				for _, keyword := range replyData.Keywords[1:] {
					reData := fmt.Sprintf("%s.*%s.*\\n", "(\\>\\>)", keyword)
					reReply = fmt.Sprintf("%s||%s", reReply, reData)
				}
			}
			isReplyKeywordOnComment, err := regexp.MatchString(reReply, com)
			if err == nil && isReplyKeywordOnComment {
				result[idx].ReplyPost = append(result[idx].ReplyPost, post)
				result[idx].ReplyString = append(result[idx].ReplyString, replyCommentSplit(com, reReply))
				continue
			}
			// 找文章開頭
			if replyData.Post.Resto == 0 && post.Resto == replyData.Post.No { // 是開頭文章 且回覆對象是目標No
				result[idx].ReplyPost = append(result[idx].ReplyPost, post)
				result[idx].ReplyString = append(result[idx].ReplyString, post.Com)
				continue
			}
		}
	}
	return result
}

func replyCommentSplit(com string, targetNo string) string {
	// 擷取其中文字
	re := regexp.MustCompile(fmt.Sprintf("(%s\\n)[\\s\\S]*(%s|[\\s\\S]$)", targetNo, _REPLY_REGEXP))
	result := re.FindString(com)
	// 刪去回覆的>>
	re2 := regexp.MustCompile(fmt.Sprintf("([^((%s)\\n)|%s])[\\s\\S]*", targetNo, _REPLY_REGEXP))
	result = re2.FindString(result)
	return result
}
