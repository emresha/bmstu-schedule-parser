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
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func parse_schedule(doc, group string) error {
	tkn := html.NewTokenizer(strings.NewReader(doc))
	for {

		tt := tkn.Next()
		txtRaw := string(tkn.Raw())
		if tt == html.ErrorToken {
			return errors.New("Group was not found. Maybe you made a mistake?")
		}
		if tt == html.StartTagToken {
			t := tkn.Token()
			if t.Data == "a" {

				if strings.Contains(txtRaw, "href=\"/schedule/") {
					href := "\"https://lks.bmstu.ru" + strings.Trim(strings.Split(txtRaw, " ")[1][6:], "\"") // this line gets the href link from the raw txt
					tkn.Next()
					hrefGroup := strings.TrimSpace(string(tkn.Raw())) // this fetches the group name
					if hrefGroup == group {
						fmt.Println(href + " " + hrefGroup)
						return nil
					}
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

	var group string

	fmt.Println("Please enter your group in Russian, e.g \"ИУ7-21Б\": ")
	fmt.Scanln(&group)

	err = parse_schedule(string(body), strings.ToUpper(group))

	if err != nil {
		fmt.Println(err)
	}
}
