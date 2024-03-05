package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

func GetWithoutBotDetection(bUrl string, urlPath string) (*http.Response, error) {
	fullUrl, err := url.JoinPath(bUrl, urlPath)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", fullUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Referer", "https://kagi.com/")

	res, err := http.DefaultClient.Do(req)

	return res, err
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

func ParseArticle(articleUrl string) (*Article, error) {
	resp, err := GetWithoutBotDetection(baseUrl, articleUrl)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	article := &Article{}

	article.Title = doc.Find("h1").Text()
	article.Subtitle = doc.Find(".header_right h2").Text()
	article.Author = strings.TrimSpace(doc.Find(".byline").Text())
	article.AuthorUrl, _ = doc.Find("a.byline").Attr("href")
	article.Date = TrimSpaceExtra(doc.Find(".dateline").Text())
	article.Url = articleUrl
	doc.Find("#body p").Each(func(i int, s *goquery.Selection) {
		var paragraphText strings.Builder
		s.Contents().Each(func(j int, node *goquery.Selection) {
			if node.Is("a") {
				href, _ := node.Attr("href")
				text := node.Text()
				paragraphText.WriteString(fmt.Sprintf(" [%s](%s) ", text, href))
			} else {
				if node.HasClass("label") {
					paragraphText.WriteString(strings.TrimSpace(node.Text()) + " ")
					return
				}

				paragraphText.WriteString(strings.TrimSpace(node.Text()))
			}
		})
		article.ContentText += strings.TrimSpace(paragraphText.String()) + "\n\n"
	})
	article.HasComments = doc.Find(".comment_count").Length() > 0
	article.CommentsUrl, _ = doc.Find(".comment_count").Attr("href")

	return article, nil
}

func TrimSpaceExtra(untidyString string) string {
	// split on //, then trim space then put them back together with " "
	parts := strings.Split(untidyString, "//")

	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	return strings.Join(parts, " ")
}
