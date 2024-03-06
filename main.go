package main

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var baseUrl = "https://www.theregister.com"

var dataDir string
var databaseDir string
var db *badger.DB

func main() {
	appDataDir, err := os.UserConfigDir()

	dataDir = filepath.Join(appDataDir, "elreg-cli")

	err = os.Mkdir(dataDir, 0755)

	if err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	// make a database folder
	databaseDir = filepath.Join(dataDir, "v1_database")

	err = os.Mkdir(databaseDir, 0755)

	if err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	opts := badger.DefaultOptions(databaseDir)

	opts.Logger = nil

	go func() {
		// run the rubbish collection every 5 minutes
		for {
			time.Sleep(5 * time.Minute)
			_ = db.RunValueLogGC(0.5)
		}
	}()

	// initialize badger db in this folder
	db, err = badger.Open(opts)

	//var currentArticle *Article = nil
	allLinks, err := FetchTheRegisterHomepageArticleLinks()

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	// save the viewed articles to a file

	if err != nil {
		panic(err)
	}

	// print first 5 links then wait for user input (one line per link)
	var linksShown int = 0

	var looping = true
	var menuNew = true
	var articleNew = false
	var currentArticle *Article = nil
	var articleLines []string = nil
	var currentArticleLine int = 0
	var input string

	for looping {
		input = ""
		// show a prompt for user input
		if !articleNew && !menuNew {
			_, _ = fmt.Scanln(&input)
		}

		input = strings.TrimSpace(strings.ToLower(input))

		if input == "?" || input == "help" || input == "h" {
			fmt.Println("Help for El Reg CLI")
			fmt.Println("Main Menu Controls:")
			fmt.Println("- `<number>`: Type the number corresponding to the headline you want to read.")
			fmt.Println("- `a`: Show all headlines at once.")
			fmt.Println("- `s`: Show the next 5 headlines.")
			fmt.Println("- `q`: Quit the application.")
			fmt.Println()
			fmt.Println("Article Controls:")
			fmt.Println("- `Enter`: Advance to the next paragraph/line of the article.")
			fmt.Println("- `a`: Display the entire article at once.")
			fmt.Println("- `q`: Go back to the main menu.")
			fmt.Println("- `r`: Mark the article as read and go back to the main menu.")
			fmt.Println("- `q`: Go back to the main menu without marking the article as read.")
			fmt.Println("- `u`: Mark the article as unread and go back to the main menu.")
			fmt.Println()
			fmt.Println("© github.com/CADawg 2024 - Licenced under the GPL v3.0")
			continue
		}

		if currentArticle != nil {
			if articleNew {
				// limit line width to 80 characters for readability
				boldText := color.New(color.Bold)
				_, _ = boldText.Printf(currentArticle.Title + "\n")
				fmt.Println(currentArticle.Subtitle)
				italicText := color.New(color.Italic)
				_, _ = italicText.Printf("By %s\n", currentArticle.Author)
				fmt.Println(currentArticle.Date)

				articleLines = strings.Split(WrapText(currentArticle.ContentText), "\n")

				for _, line := range articleLines[:5] {
					fmt.Println(line)
					currentArticleLine++
				}

				articleNew = false
				continue
			}

			if input == "" {
				if currentArticleLine >= len(articleLines) {
					fmt.Println("End of article")
					continue
				}
				for i, line := range articleLines[currentArticleLine : currentArticleLine+1] {
					if i == 0 {
						fmt.Printf("\u001B[F%s\n", line)
					} else {
						fmt.Println(line)
					}
					currentArticleLine++
				}
			} else if input == "a" {
				if currentArticleLine >= len(articleLines) {
					fmt.Println("End of article")
					continue
				}
				// show all lines of the article
				for i, line := range articleLines[currentArticleLine:] {
					if i == 0 {
						fmt.Printf("\u001B[F%s\n", line)
					} else {
						fmt.Println(line)
					}
					currentArticleLine++
				}
			} else if input == "s" {
				if currentArticleLine >= len(articleLines) {
					fmt.Println("End of article")
					continue
				}
				// show next 5 lines of the article
				for i, line := range articleLines[currentArticleLine : currentArticleLine+5] {
					if i == 0 {
						fmt.Printf("\u001B[F%s\n", line)
					} else {
						fmt.Println(line)
					}
					currentArticleLine++
				}
			} else if input == "q" {
				// go back to the homepage
				currentArticle = nil
				articleNew = false
				currentArticleLine = 0
				articleLines = nil
				menuNew = true
				linksShown = 0
				for i := 0; i < 5; i++ {
					fmt.Println()
				}
			} else if input == "r" {
				// mark the article as read
				SetReadStatusBadger(currentArticle.Url, true)
				// go back to the homepage
				currentArticle = nil
				articleNew = false
				currentArticleLine = 0
				articleLines = nil
				menuNew = true
				linksShown = 0
				for i := 0; i < 5; i++ {
					fmt.Println()
				}
			} else if input == "q" {
				// go back to the homepage
				currentArticle = nil
				articleNew = false
				currentArticleLine = 0
				articleLines = nil
				menuNew = true
				linksShown = 0
				for i := 0; i < 5; i++ {
					fmt.Println()
				}
			} else if input == "u" {
				SetReadStatusBadger(currentArticle.Url, false)
				// go back to the homepage
				currentArticle = nil
				articleNew = false
				currentArticleLine = 0
				articleLines = nil
				menuNew = true
				linksShown = 0
				for i := 0; i < 5; i++ {
					fmt.Println()
				}
			} else {
				fmt.Println("Unknown command. Type `?` for help.")
			}
		} else {
			if menuNew {
				links := allLinks[:5]
				for i, link := range links {
					PrintArticleLink(link, linksShown, i == 0)
					linksShown++
				}
				menuNew = false
				continue
			}

			// standard homepage loop
			if input == "" {
				PrintArticleLink(allLinks[linksShown], linksShown, true)
				linksShown++
			} else if input == "a" {
				for i, link := range allLinks[linksShown:] {
					PrintArticleLink(link, linksShown, i == 0)
					linksShown++
				}
			} else if input == "s" {
				// show next 5 links
				for i, link := range allLinks[linksShown : linksShown+5] {
					PrintArticleLink(link, linksShown, i == 0)
					linksShown++
				}
			} else if input == "q" {
				looping = false
			} else {
				parsedInt, err := strconv.Atoi(input)

				if err != nil {
					fmt.Println("Invalid input. Type `?` for help.")
					continue
				}

				if parsedInt < 0 || parsedInt >= len(allLinks) {
					fmt.Println("Unknown article. Type `?` for help.")
					continue
				}

				// get the article
				currentArticle, err = ParseArticle(allLinks[parsedInt].Url)

				if err != nil {
					fmt.Println("Error fetching article: ", err)
				} else {
					// print the article
					articleNew = true
				}
			}
		}
	}

	// close the db
	err = db.Close()

	if err != nil {
		panic(err)
	}
}

type LocalPostInfo struct {
	Url    string
	IsRead bool
}

func PrintArticleLink(article ArticleLink, index int, isFirstLine ...bool) {
	if len(isFirstLine) > 0 && isFirstLine[0] {
		fmt.Printf("\u001B[F%d) %s - %s (%s) %s\n", index, article.Title, article.Subtitle, article.Date, GetReadText(GetReadStatusBadger(article.Url)))
	} else {
		fmt.Printf("%d) %s - %s (%s) %s\n", index, article.Title, article.Subtitle, article.Date, GetReadText(GetReadStatusBadger(article.Url)))
	}
}

func WrapText(text string) string {
	const maxLineLength = 80
	var result strings.Builder
	words := strings.Fields(text)
	currentLineLength := 0

	for _, word := range words {
		wordLength := len(word)
		if currentLineLength+wordLength+1 > maxLineLength && wordLength < maxLineLength {
			result.WriteString("\n")
			currentLineLength = 0
		}
		result.WriteString(word + " ")
		currentLineLength += wordLength + 1
	}

	return strings.TrimSpace(result.String())
}

func GetReadText(read bool) string {
	if read {
		return "[Read]"
	} else {
		return ""
	}
}
