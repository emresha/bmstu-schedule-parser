package main

/*

this project is made entirely
for the purpose of practice
in parsing HTML and working with
net/http library in golang.
it is not intended to be
professionaly done nor it is
going to be updated in the future

*/

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/net/html"
)

func findSchedule(doc, group string) (string, error) {
	tkn := html.NewTokenizer(strings.NewReader(doc))
	for {

		tt := tkn.Next()
		txtRaw := string(tkn.Raw())
		if tt == html.ErrorToken { // if end of html file encountered then link wasn't found
			return "", errors.New("Group was not found. Perhaps you made a mistake?")
		}
		if tt == html.StartTagToken {
			t := tkn.Token()
			if t.Data == "a" {

				if strings.Contains(txtRaw, "href=\"/schedule/") {
					// the line below gets the href link from the raw txt.
					href := "https://lks.bmstu.ru" + strings.Trim(strings.Split(txtRaw, " ")[1][6:], "\"")
					tkn.Next()
					hrefGroup := strings.TrimSpace(string(tkn.Raw()))            // this fetches the group name.
					if strings.Contains(hrefGroup, group) || group == "--TEST" { // <- finally, if group name matches the query, we print it.
						fmt.Println(href + " " + hrefGroup)
						if group != "--TEST" {
							return href, nil
						}

					}
				}
			}

		}

	}

}

func openBrowser(url string) { // this is temporarily
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

// not done yet
// func readConfig(path string) bool {
// 	f, err := os.Open(path)
// 	if err != nil { // if any error than set everything to default
// 		return false
// 	}

// 	return false
// }

func main() {
	browser := false // if true opens browser when finds schedule, else does nothing.

	h, err := http.Get("https://lks.bmstu.ru/schedule/list") // here we connect to the website
	if err != nil || h.StatusCode == 404 {                   // if any error encountered
		fmt.Println("Couldn't connect to bmstu's website. Did you write the correct URL? Or is it possibly down?")
		return
	}

	fmt.Printf("Successfully connected to lks.bmstu.ru on %s\n", h.Header.Get("Date"))
	body, err := io.ReadAll(h.Body)

	defer h.Body.Close()

	if err != nil {
		fmt.Println("Unexpected error encountered.")
		return
	}

	var group string

	fmt.Println("Please enter your group in Russian, e.g \"ИУ7-21Б\": ")
	fmt.Scanln(&group)

	link, err := findSchedule(string(body), strings.ToUpper(group))

	if err != nil {
		fmt.Println(err)
	}

	if group != "--test" && browser {
		openBrowser(link)
	}

	fmt.Println("Press CTRL + C or ENTER to exit")
	fmt.Scanf("h")
}
