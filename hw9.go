package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
)

/*
main function to which contains the goroutine for the file itself. Parsing through the url's webpage to then pass to functions to do
operations to get the url (if accessible), length of the body, and title names (if accessible)
*/
func main() {
	fileRead, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer fileRead.Close()
	// loading file URLs into a slice string
	sc := bufio.NewScanner(fileRead)
	var url_lines []string
	for sc.Scan() {
		url_lines = append(url_lines, sc.Text())
	}

	var waiting sync.WaitGroup

	// starting a goroutine function for each element in the string slice
	for _, url := range url_lines {
		waiting.Add(1)
		url := url
		// go routine
		go func() {
			defer waiting.Done()
			err = searchingThr(url)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	waiting.Wait()
}

// function to which is going through each url to find the title's for the url
func searchingThr(url_single string) error {
	url_page, error := searchingUrl(url_single)

	if error != nil {
		return error
	}
	body_length := len(url_page)
	// regular expressions to locate the title itself alongside the ending of the title
	var title = regexp.MustCompile(`(<title)([\s\S]*)(<\/title>)`)
	var lastChar = regexp.MustCompile(`<.+?>`)

	extract := title.FindString(url_page)
	extract = lastChar.ReplaceAllString(extract, "")

	if extract != "" {
		fmt.Println("The length of the url " + url_single + " body is: " + strconv.Itoa(body_length))
		fmt.Println("The title(s) found in " + url_single + " are: " + extract)
	}
	return error
}

// function to handle the error messaging with self-defined error message alongsid system error message displayed within the terminal
func searchingUrl(url string) (string, error) {
	resp, err := http.Get(url)
	var empty string
	if err != nil {
		fmt.Println("There was no title nor body length found for the url titled: " + url)
		return empty, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Println()

	if err != nil {
		return empty, err
	}
	return string(body), nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
