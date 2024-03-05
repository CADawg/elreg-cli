package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

var baseUrl = "https://www.theregister.com"

func main() {
	//var currentArticle *Article = nil

	links, err := FetchTheRegisterHomepageArticleLinks()

	if err != nil {
		fmt.Println("Error fetching article links: ", err)
	}

	for _, link := range links {
		fmt.Println("Title: ", link.Title)
		fmt.Println("Subtitle: ", link.Subtitle)
		fmt.Println("Date: ", link.Date)
		fmt.Println("Url: ", link.Url)
		fmt.Println("=====================================")
	}

	article, err := ParseArticle("/2024/03/01/in_the_vanguard_of_21st/")

	if err != nil {
		fmt.Println("Error fetching article: ", err)
	}

	fmt.Println("Title: ", article.Title)
	fmt.Println("Subtitle: ", article.Subtitle)
	fmt.Println("Date: ", article.Date)
	fmt.Println("Url: ", article.Url)
	fmt.Println("Content: ", article.ContentText)
	fmt.Println("Author: ", article.Author)
	fmt.Println("Author URL: ", article.AuthorUrl)
	fmt.Println("=====================================")
}

func FetchTheRegisterHomepageArticleLinks() ([]ArticleLink, error) {
	// fetch the register homepage
	resp, err := GetWithoutBotDetection(baseUrl, "")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

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
		subtitle := s.Find("div.standfirst").Text()
		urlLink, _ := s.Attr("href")

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
