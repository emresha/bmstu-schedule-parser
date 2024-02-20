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
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func parse_schedule(doc string) {
	tkn := html.NewTokenizer(strings.NewReader(doc))
	isAFlag := false
	fmt.Println(isAFlag)
	for {

		tt := tkn.Next()

		if tt == html.ErrorToken {
			return
		}
		if tt == html.StartTagToken {
			t := tkn.Token()
			if t.Data == "a" {
				isAFlag = true
				txtRaw := string(tkn.Raw())
				if strings.Contains(txtRaw, "href=\"/schedule/") {
					fmt.Println(txtRaw)
				}
			}
		}

	}
}

func main() {
	h, err := http.Get("https://lks.bmstu.ru/schedule/list") // here we connect to the website
	if err != nil || h.StatusCode == 404 {                   // if any error encountered
		fmt.Println("Couldn't connect to bmstu's website. Did you write the correct URL? Or is it possibly down?")
		return
	}

	fmt.Printf("Successfully connected to lks.bmstu.ru on %s\n", h.Header.Get("Date"))
	body, err := io.ReadAll(h.Body)
	if err != nil {
		fmt.Println("Unexpected error encountered.")
		return
	}

	parse_schedule(string(body))
}
