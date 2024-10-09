package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Book represents a book result from LibGen with all required fields
type Book struct {
	Title     string
	Author    string
	Publisher string
	Year      string
	Language  string
	Pages     string
	Filesize  string
	Extension string
	BookURL   string
	CoverURL  string
	ID        string
}

// gsSearch searches LibGen and returns the book details as a slice of Book structs
func gsSearch(query string) ([]*Book, error) {
	// Replace spaces with "+" for URL encoding
	query = strings.ReplaceAll(query, " ", "+")

	// Construct the LibGen search URL
	url := fmt.Sprintf("https://libgen.gs/index.php?req=%s&columns%%5B%%5D=t&columns%%5B%%5D=a&columns%%5B%%5D=s&columns%%5B%%5D=y&columns%%5B%%5D=p&columns%%5B%%5D=i&objects%%5B%%5D=f&objects%%5B%%5D=e&objects%%5B%%5D=s&objects%%5B%%5D=a&objects%%5B%%5D=p&objects%%5B%%5D=w&topics%%5B%%5D=l&topics%%5B%%5D=c&topics%%5B%%5D=f&topics%%5B%%5D=a&topics%%5B%%5D=m&topics%%5B%%5D=r&topics%%5B%%5D=s&res=100&covers=on&filesuns=all", query)

	// Make an HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	// Parse the response HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Slice to hold the books with detailed information
	var books []*Book

	// Debug: check if the table with ID "tablelibgen" exists
	if doc.Find("#tablelibgen").Length() == 0 {
		log.Println("Table with ID #tablelibgen not found!")
		return nil, nil
	}

	// Debug: Print number of rows found in the table
	rows := doc.Find("#tablelibgen tbody tr")
	log.Printf("Found %d rows in the table", rows.Length())

	// Find all book information in the table and extract the details
	rows.Each(func(i int, s *goquery.Selection) {
		// Extract the cover image URL from the <img> tag
		coverTag := s.Find("td").Eq(0).Find("img").First()
		coverURL, exists := coverTag.Attr("src")
		if !exists {
			coverURL = "" // Handle cases where there's no image
		}

		// Extract title, author, year, publisher, pages, language, filesize, and extension
		titleTag := s.Find("td").Eq(1).Find("a").First()
		title := strings.TrimSpace(titleTag.Text())

		author := strings.TrimSpace(s.Find("td").Eq(2).Text())    // Author is in the third <td>
		publisher := strings.TrimSpace(s.Find("td").Eq(3).Text()) // Publisher is in the fourth <td>
		year := strings.TrimSpace(s.Find("td").Eq(4).Text())      // Year is in the fifth <td>
		language := strings.TrimSpace(s.Find("td").Eq(5).Text())  // Language is in the sixth <td>
		pages := strings.TrimSpace(s.Find("td").Eq(6).Text())     // Pages is in the seventh <td>
		filesize := strings.TrimSpace(s.Find("td").Eq(7).Text())  // Filesize is in the eighth <td>
		extension := strings.TrimSpace(s.Find("td").Eq(8).Text()) // Extension is in the ninth <td>

		// Variable to store the final URL
		var bookURL string
		baseURL := "https://libgen.gs/"

		// Check for both edition.php and series.php links
		var editionURL, seriesURL string

		s.Find("a").Each(func(index int, link *goquery.Selection) {
			href, hrefExists := link.Attr("href")
			if hrefExists {
				// Check if the link contains "edition.php"
				if strings.Contains(href, "edition.php") {
					editionURL = baseURL + href
				}
				// Check if the link contains "series.php"
				if strings.Contains(href, "series.php") {
					seriesURL = baseURL + href
				}
			}
		})

		// Prioritize edition.php if found; otherwise, use series.php
		if editionURL != "" {
			bookURL = editionURL
		} else if seriesURL != "" {
			bookURL = seriesURL
		}

		// Extract the ID from the edition link or fallback to series link
		id := ""
		if editionURL != "" {
			id = strings.Split(editionURL, "=")[1]
		} else if seriesURL != "" {
			id = strings.Split(seriesURL, "=")[1]
		}

		// Add the book to the list if it has a title and ID
		if title != "" && id != "" {
			books = append(books, &Book{
				Title:     title,
				Author:    author,
				Publisher: publisher,
				Year:      year,
				Language:  language,
				Pages:     pages,
				Filesize:  filesize,
				Extension: extension,
				CoverURL:  fmt.Sprintf("https://libgen.gs%s", coverURL), // Construct full cover URL
				ID:        id,
				BookURL:   bookURL, // Set the correct URL
			})
		}
	})

	return books, nil
}
