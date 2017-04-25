package types

// 初始化登陆的公众号消息体
type MPSubscribeMsg struct {
	UserName       string
	MPArticleCount int
	MPArticleList  []MPArticle
	Time           int
	NickName       string
}

type MPArticle struct {
	Title  string
	Digest string
	Cover  string
	Url    string
}
