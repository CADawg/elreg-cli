package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func GetWithoutBotDetection(urlPath string) (*http.Response, error) {
	req, err := http.NewRequest("GET", urlPath, nil)

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

func ParseArticle(articleUrl string) (*Article, error) {
	resp, err := GetWithoutBotDetection(articleUrl)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

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

func FetchTheRegisterHomepageArticleLinks() ([]ArticleLink, error) {
	// fetch the register homepage
	resp, err := GetWithoutBotDetection(baseUrl)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// parse the HTML to find the article links
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	// find the article links
	var articleLinks []ArticleLink

	doc.Find("a.story_link").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title, subtitle, timestamp and URL
		title := s.Find("h4").Text()
		subtitle := RemoveAllNewlines(s.Find("div.standfirst").Text())
		urlLink, _ := s.Attr("href")

		// remove any query string as url.join path escapes it, and it breaks the url
		urlLink = strings.Split(urlLink, "?")[0]

		urlLink, err := url.JoinPath(baseUrl, urlLink)

		if err != nil {
			// we want consistency for db saving read status so forget it.
			return
		}

		// get 2 parents above to check it's not sponsored
		parent := s.ParentsFiltered("article").ParentsFiltered("div")

		// it's sponsored, skip
		if parent.HasClass("other_stories") {
			return
		}

		// also check the .section_name's text isn't "Webinar" as that's sponsored too
		sectionName := s.ParentsFiltered("article").Find(".section_name").Text()

		// it's sponsored, skip
		if sectionName == "Webinar" {
			return
		}

		date := s.Find(".time_stamp").Text()

		if len(date) == 0 {
			// fetch from url
			date, _ = DateFromUrl(urlLink)
		}

		articleLinks = append(articleLinks, ArticleLink{
			Title:    title,
			Subtitle: subtitle,
			Date:     date,
			Url:      urlLink,
		})
	})

	return articleLinks, nil
}

// RemoveAllNewlines removes all newlines from a string
// as well as applying some other fixes to make text more uniform
// there's probably a better way to do this
func RemoveAllNewlines(s string) string {
	return FeatureExtraSpaceRemove(UpdatedExtraSpaceRemove(strings.ReplaceAll(s, "\n", "")))
}

func UpdatedExtraSpaceRemove(s string) string {
	split := strings.Split(s, "Updated")

	for i, part := range split {
		split[i] = strings.TrimSpace(part)
	}

	return strings.Join(split, "Updated ")
}

func FeatureExtraSpaceRemove(s string) string {
	split := strings.Split(s, "Feature")

	for i, part := range split {
		split[i] = strings.TrimSpace(part)
	}

	return strings.Join(split, "Feature ")
}

func TryFetchArticleFromArticleLink(input string, allLinks []ArticleLink) (*Article, error) {
	parsedInt, err := strconv.Atoi(input)

	if err != nil {
		return nil, errors.New("invalid input. Type `?` for help")
	}

	if parsedInt < 0 || parsedInt >= len(allLinks) {
		return nil, errors.New("invalid input. Type `?` for help")
	}

	// get the article
	currentArticle, err := ParseArticle(allLinks[parsedInt].Url)

	if err != nil {
		return nil, errors.New("error fetching article")
	}

	return currentArticle, err
}
