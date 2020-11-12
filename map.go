package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const length = 10

func generateURL(in string) string {
	in = strings.TrimSpace(in)
	in = strings.ToLower(in)
	in = strings.Replace(in, " ", "+", -1)
	url := fmt.Sprintf("https://www.music-map.com/%s", in)

	return url
}

func menu(options *goquery.Selection, in *bufio.Reader) string {
	//add all related names to list
	names := make([]string, 0)
	options.Each(func(i int, s *goquery.Selection) {
		names = append(names, s.Text())
	})

	//present options
	if len(names) == 0 {
		return "bad"
	}
	fmt.Println("\nRelated artists:")
	start, end := 1, 1+length
	for i := start; i < end; i++ {
		fmt.Printf("%d: %s\n", i, names[i])
	}

	//interactive prompting
	var result string
	loop := true
	for loop {
		fmt.Printf("\n> ")

		cmd, _ := in.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if cmd == "more" {
			start += length
			end = start + length
			if end > len(names) {
				end = len(names)
				start = end - 10
			}
			for i := start; i < end; i++ {
				fmt.Printf("%d: %s\n", i+1, names[i])
			}

			if end == len(names) {
				fmt.Println("END OF RESULTS")
			}
		}

		if cmd == "go" {
			fmt.Printf("Enter artist name: ")
			arg, _ := in.ReadString('\n')
			result = generateURL(arg)
			loop = false
		}

		if cmd == "exit" {
			return "none"
		}

		if num, err := strconv.Atoi(cmd); err == nil {
			return generateURL(names[num])
		}
	}

	return result
}

func main() {
	//setup reader
	in := bufio.NewReader(os.Stdin)

	//make first url
	fmt.Printf("Enter an artist to begin: ")
	arg, _ := in.ReadString('\n')
	url := generateURL(arg)

	loop := true

	for loop {
		//get response
		response, err := http.Get(url)
		if err != nil {
			log.Fatal("couldnt find page: ", err)
		}
		defer response.Body.Close()

		//gather body
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatal("couldn't parse page: ", err)
		}

		//parse and prompt
		options := doc.Find("div#gnodMap a")
		next := menu(options, in)

		//loop logic
		if next == "none" {
			loop = false
		} else if next == "bad" {
			fmt.Println("No results found")
			fmt.Printf("Try a different artist: ")
			arg, _ := in.ReadString('\n')
			url = generateURL(arg)
		} else {
			url = next
		}
	}
}
