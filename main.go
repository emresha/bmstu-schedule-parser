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

func parse_schedule(doc, group string) {
	tkn := html.NewTokenizer(strings.NewReader(doc))
	isAFlag := false
	txtRaw := ""
	fmt.Println(isAFlag)
	prev := []string{}
	for {

		tt := tkn.Next()
		txtRaw = string(tkn.Raw())
		if tt == html.ErrorToken {
			return
		}
		if tt == html.StartTagToken {
			t := tkn.Token()
			if t.Data == "a" {
				isAFlag = true

				if strings.Contains(txtRaw, "href=\"/schedule/") {
					fmt.Println(txtRaw)
					prev = append(prev, txtRaw)
				}
			}

		}

		if isAFlag && tt == html.TextToken {
			if group == strings.TrimSpace(txtRaw) {
				fmt.Println(txtRaw) // if name is equal to group print it
				return
			}

			isAFlag = false
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

	parse_schedule(string(body), group)
}
