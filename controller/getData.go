package controller

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Scrape get telegra page data
func Scrape(cors bool, url string) ([]string, string, bool) {
	// Request the HTML page.
	var rootURL string = "https://telegra.ph"
	// var findthis []string
	if cors {
		url = strings.Join([]string{"https://cors.izumana.ml", url}, "/?url=")
		rootURL = strings.Join([]string{"https://cors.izumana.ml", rootURL}, "/?url=")
	}
	res, err := http.Get(url)
	if err != nil {
		return []string{"bad"}, "bad", false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return []string{res.Status}, res.Status, false
	}
	var list []string
	// Load the HTML document
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	title := doc.Find("article#_tl_editor  h1").Text()
	if strings.Contains(title, "视频") {
		n := strings.LastIndex(title, "视频")
		title = title[0:n]
	}
	title = strings.Replace(title, " ", "_", -1)
	if strings.LastIndex(title, "_") == len(title)-1 {
		title = strings.TrimRight(title, "_")
	}
	doc.Find("article#_tl_editor img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		src = strings.Join([]string{rootURL, src}, "")
		list = append(list, src)
	})
	// fmt.Println(findthis)
	return list, title, true
}
