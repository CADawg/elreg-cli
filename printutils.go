package main

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

func PrintNextXArticleLines(articleLines []string, currentIndex int, x int, replacePrevious ...bool) int {
	if currentIndex >= len(articleLines) {
		fmt.Println("End of article")
		return currentIndex
	}

	if len(replacePrevious) > 0 && replacePrevious[0] {
		fmt.Print("\u001B[F")
	}

	// loop through the next x lines or until the end of the list if x is -1
	for i := currentIndex; i < currentIndex+x || x == -1; i++ {
		if i >= len(articleLines) {
			break
		}

		fmt.Println(articleLines[i])
	}

	return currentIndex + x
}

func PrintArticleLink(article ArticleLink, index int) {
	fmt.Printf("%d) %s - %s (%s) %s\n", index, article.Title, article.Subtitle, article.Date, GetReadText(GetReadStatusBadger(article.Url)))
}

func PrintNextXLinks(links []ArticleLink, currentIndex int, x int, replacePrevious ...bool) int {
	if len(replacePrevious) > 0 && replacePrevious[0] {
		fmt.Print("\u001B[F")
	}

	// loop through the next x links or until the end of the list if x is -1
	for i := currentIndex; i < currentIndex+x || x == -1; i++ {
		if i >= len(links) {
			break
		}

		PrintArticleLink(links[i], i)
	}

	return currentIndex + x
}

func PrintBlankLines(x int) {
	for i := 0; i < x; i++ {
		fmt.Println()
	}
}

func PrintHelp() {
	fmt.Println("El Reg CLI")
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
}

func PrintArticleHeading() {
	// limit line width to 80 characters for readability
	boldText := color.New(color.Bold)
	_, _ = boldText.Printf(currentArticle.Title + "\n")
	fmt.Println(currentArticle.Subtitle)
	italicText := color.New(color.Italic)
	_, _ = italicText.Printf("By %s\n", currentArticle.Author)
	fmt.Println(currentArticle.Date)
	if currentArticle.HasComments {
		_, _ = boldText.Printf("Comments: %s\n", currentArticle.CommentsUrl)
	}

	fmt.Println()

	articleLines = strings.Split(WrapText(currentArticle.ContentText), "\n")

	for _, line := range articleLines[:5] {
		fmt.Println(line)
		currentArticleLine++
	}
}
