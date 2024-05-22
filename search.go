package main

import (
	"fmt"
	"github.com/ciehanski/libgen-cli/libgen"
	"net/url"
)

// SearchLibgen searches libgen for a provided title
func SearchLibgen(title string) []*libgen.Book {
	//fmt.Println("Searching libgen for: ", title)
	// SearchMirror: libgen.GetWorkingMirror(libgen.SearchMirrors)
	results, err := libgen.Search(&libgen.SearchOptions{Query: title, SearchMirror: url.URL{Host: "libgen.rs", Scheme: "https"}, Results: 10}) //, Extension: []string{"pdf"}})
	if err != nil {
		fmt.Println("Error searching libgen: ", err)
		//return results
	}
	/*
		fmt.Println(fmt.Sprintf("Search Results for %s", title))

		for _, result := range results {
			fmt.Println("\n\n--------------------------------------------------")
			fmt.Println(fmt.Sprintf("Title:%s", result.Title))
			fmt.Println(fmt.Sprintf("Extension:%s", result.Extension))
			fmt.Println(fmt.Sprintf("MD5:%s", result.Md5))
			fmt.Println(fmt.Sprintf("ID:%s", result.ID))

			filesize, err := strconv.Atoi(result.Filesize)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(fmt.Sprintf("Filesize:%s", prettyByteSize(filesize)))
			fmt.Println(fmt.Sprintf("Year:%s", result.Year))
			fmt.Println(fmt.Sprintf("Author:%s", result.Author))
			fmt.Println(fmt.Sprintf("Publisher:%s", result.Publisher))
			fmt.Println(fmt.Sprintf("Edition:%s", result.Edition))
			fmt.Println(fmt.Sprintf("Language:%s", result.Language))
			fmt.Println(fmt.Sprintf("Pages:%s", result.Pages))
			fmt.Println(fmt.Sprintf("Cover URL: https://library.lol/covers/%s", result.CoverURL))
			fmt.Println(fmt.Sprintf("Download: https://library.lol/main/%s", result.Md5))
			fmt.Println("--------------------------------------------------")
		}

	*/
	return results
}
