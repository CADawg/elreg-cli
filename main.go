package main

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var baseUrl = "https://www.theregister.com"

var dataDir string
var databaseDir string
var db *badger.DB

var lengthMap = map[string]int{
	"":  1,
	"a": -1,
	"s": 5,
}

// print first 5 links then wait for user input (one line per link)
var linksShown = 0

var looping = true
var menuNew = true
var articleNew = false
var currentArticle *Article = nil
var articleLines []string = nil
var currentArticleLine = 0
var input string

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

	for looping {
		input = ""
		// show a prompt for user input
		if !articleNew && !menuNew {
			_, _ = fmt.Scanln(&input)
		}

		input = strings.TrimSpace(strings.ToLower(input))

		if input == "?" || input == "help" || input == "h" {
			PrintHelp()
			continue
		}

		if currentArticle != nil {
			if articleNew {
				PrintArticleHeading()

				articleNew = false
				continue
			}

			switch input {
			case "", "a", "s":
				nLines := GetNLines(input)
				currentArticleLine = PrintNextXArticleLines(articleLines, currentArticleLine, nLines, true)
			case "q", "r", "u":
				if input == "r" {
					SetReadStatusBadger(currentArticle.Url, true)
				}
				if input == "u" {
					SetReadStatusBadger(currentArticle.Url, false)
				}
				resetPostState()
				PrintBlankLines(5)
			default:
				fmt.Println("Unknown command. Type `?` for help.")
			}
		} else {
			if menuNew {
				linksShown = PrintNextXLinks(allLinks, 0, 5, false)
				menuNew = false
				continue
			}

			// standard homepage loop
			switch input {
			case "", "a", "s":
				nLines := GetNLines(input)
				linksShown = PrintNextXLinks(allLinks, linksShown, nLines, true)
			case "q":
				looping = false
			default:
				currentArticle, err = TryFetchArticleFromArticleLink(input, allLinks)

				if err != nil {
					fmt.Println(err)
					continue
				} else {
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
