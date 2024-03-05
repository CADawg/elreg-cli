package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var baseUrl = "https://www.theregister.com"

var dataDir string
var databaseDir string
var db *badger.DB

func main() {
	//var currentArticle *Article = nil
	allLinks, err := FetchTheRegisterHomepageArticleLinks()

	if err != nil {
		panic(err)
	}

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

	// hide logging output as this is a cli app so that'll clog up the terminal
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

	if err != nil {
		panic(err)
	}

	// save the viewed articles to a file

	if err != nil {
		panic(err)
	}

	// print first 5 links then wait for user input (one line per link)
	links := allLinks[:5]
	var linksShown int = 0

	for _, link := range links {
		fmt.Printf("%d) %s - %s (%s)\n", linksShown, link.Title, link.Subtitle, link.Date)
		linksShown++
	}

	var looping bool = true
	var currentArticle *Article = nil

	for looping {
		// show a prompt for user input
		var input string
		_, _ = fmt.Scanln(&input)

		if currentArticle != nil {

		} else {
			// standard homepage loop
			if input == "" {
				fmt.Printf("\033[F%d) %s %s - %s\n", linksShown, allLinks[linksShown].Date, allLinks[linksShown].Title, allLinks[linksShown].Subtitle)
				linksShown++
			} else if input == "a" {
				for i, link := range allLinks[linksShown:] {
					if i == 0 {
						fmt.Printf("\u001B[F%d) %s - %s (%s)\n", linksShown, link.Title, link.Subtitle, link.Date)
					} else {
						fmt.Printf("%d) %s - %s (%s)\n", linksShown, link.Title, link.Subtitle, link.Date)
					}
					linksShown++
				}
			} else if input == "s" {
				// show next 5 links
				for i, link := range allLinks[linksShown : linksShown+5] {
					if i == 0 {
						fmt.Printf("\u001B[F%d) %s - %s (%s)\n", linksShown, link.Title, link.Subtitle, link.Date)
					} else {
						fmt.Printf("%d) %s - %s (%s)\n", linksShown, link.Title, link.Subtitle, link.Date)
					}
					linksShown++
				}
			} else if input == "q" {
				looping = false
			} else {
				parsedInt, err := strconv.Atoi(input)

				if err != nil {
					fmt.Println("Invalid input")
				}

				if parsedInt < 0 || parsedInt >= len(allLinks) {
					fmt.Println("Unknown article")
				}
			}
		}
	}

	// close the db
	err = db.Close()

	if err != nil {
		panic(err)
	}

	/*article, err := ParseArticle("/2024/03/01/in_the_vanguard_of_21st/")

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
	fmt.Println("=====================================")*/
}

type LocalPostInfo struct {
	Url    string
	IsRead bool
}

func GetReadStatusBadger(homepageUrl string) bool {
	var localPostInfo LocalPostInfo

	_ = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("v1_" + homepageUrl))

		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &localPostInfo)
		})
	})

	return localPostInfo.IsRead
}

func SetReadStatusBadger(homepageUrl string, readStatus bool) {
	localPostInfo := LocalPostInfo{Url: homepageUrl, IsRead: readStatus}

	_ = db.Update(func(txn *badger.Txn) error {
		encoded, err := json.Marshal(localPostInfo)

		if err != nil {
			return err
		}

		return txn.Set([]byte("v1_"+homepageUrl), encoded)
	})
}
