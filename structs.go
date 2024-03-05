package main

type ArticleLink struct {
	Title    string
	Subtitle string
	Date     string
	Url      string
}

type Article struct {
	Title       string
	Subtitle    string
	Author      string
	AuthorUrl   string
	Date        string
	Url         string
	ContentText string
	HasComments bool
	CommentsUrl string
}
