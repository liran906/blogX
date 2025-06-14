// Path: ./service/redis_service/redis_article/get.go

package redis_article

import (
	"blogX_server/global"
	"blogX_server/models"
	"strconv"
)

func get(t articleCacheType, articleID uint) int {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(articleID))).Int()
	return num
}

func GetArticleRead(articleID uint) int {
	return get(ArticleReadCount, articleID)
}
func GetArticleLike(articleID uint) int {
	return get(ArticleLikeCount, articleID)
}
func GetArticleCollect(articleID uint) int {
	return get(ArticleCollectCount, articleID)
}
func GetArticleComment(articleID uint) int {
	return get(ArticleCommentCount, articleID)
}

func getAllArticleCache(t articleCacheType) map[uint]int {
	res, err := global.Redis.HGetAll(string(t)).Result()
	if err != nil {
		return nil
	}
	mps := make(map[uint]int)
	for k, v := range res {
		key, err1 := strconv.Atoi(k)
		val, err2 := strconv.Atoi(v)
		if err1 != nil || err2 != nil {
			continue // skip this invalid entry
		}
		mps[uint(key)] = val
	}
	return mps
}

func GetAllReadCounts() map[uint]int {
	return getAllArticleCache(ArticleReadCount)
}
func GetAllLikeCounts() map[uint]int {
	return getAllArticleCache(ArticleLikeCount)
}
func GetAllCollectCounts() map[uint]int {
	return getAllArticleCache(ArticleCollectCount)
}
func GetAllCommentCounts() map[uint]int {
	return getAllArticleCache(ArticleCommentCount)
}

func GetAllFields(articleID uint) map[string]int {
	mps := map[string]int{
		"read":    GetArticleRead(articleID),
		"like":    GetArticleLike(articleID),
		"collect": GetArticleCollect(articleID),
		"comment": GetArticleComment(articleID),
	}
	return mps
}
func UpdateCachedFieldsForArticle(article *models.ArticleModel) (ok bool) {
	if article == nil || article.ID == 0 {
		return false
	}
	mps := GetAllFields(article.ID)
	article.ReadCount += mps["read"]
	article.LikeCount += mps["like"]
	article.CollectCount += mps["collect"]
	article.CommentCount += mps["comment"]
	return true
}
