package main

import (
	"fmt"
	"net/url"
	"strings"
)

// DateFromUrl Get the date from the URL
// if there is a domain in the url, ignore it and skip until the path
func DateFromUrl(urlAdd string) (string, error) {
	parts, err := url.Parse(urlAdd)

	if err != nil {
		return "", err
	}

	// split the path
	pathParts := strings.Split(parts.Path, "/")

	fmt.Println(pathParts)

	// convert month to name
	switch pathParts[2] {
	case "01":
		pathParts[2] = "Jan"
	case "02":
		pathParts[2] = "Feb"
	case "03":
		pathParts[2] = "Mar"
	case "04":
		pathParts[2] = "Apr"
	case "05":
		pathParts[2] = "May"
	case "06":
		pathParts[2] = "Jun"
	case "07":
		pathParts[2] = "Jul"
	case "08":
		pathParts[2] = "Aug"
	case "09":
		pathParts[2] = "Sep"
	case "10":
		pathParts[2] = "Oct"
	case "11":
		pathParts[2] = "Nov"
	case "12":
		pathParts[2] = "Dec"
	}

	// 0 is year, 1 is month, 2 is day
	return pathParts[3] + " " + pathParts[2] + " " + pathParts[1], nil
}
