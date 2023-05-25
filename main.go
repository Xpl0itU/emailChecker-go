package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"

	"github.com/joho/godotenv"
)

type Filter struct {
	Mail    string
	Subject string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	filtersContent, err := ioutil.ReadFile("filters.json")
	if err != nil {
		log.Fatal("Error loading filters.json file")
	}
	var filters []Filter
	_ = json.Unmarshal([]byte(filtersContent), &filters)

	server := os.Getenv("SERVER")
	username := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	c, err := client.DialTLS(fmt.Sprintf("%s:%d", server, 993), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	if err := c.Login(username, password); err != nil {
		log.Fatal(err)
	}

	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	ret := 0

	for i := 0; i < len(filters); i++ {
		criteria := imap.NewSearchCriteria()
		fromAddress := filters[i].Mail
		subject := filters[i].Subject
		today, err := time.Parse("02-Jan-2006", time.Now().Format("02-Jan-2006"))
		if err != nil {
			log.Fatal(err)
		}
		criteria.Header.Set("From", fromAddress)
		criteria.Header.Set("Subject", subject)
		criteria.SentSince = today

		seqNums, err := c.Search(criteria)
		if err != nil {
			log.Fatal(err)
		}

		if len(seqNums) == 0 {
			fmt.Printf("No matching emails found for sender %s and subject %s\n", fromAddress, subject)
			ret = 1
			continue
		}
	}
	os.Exit(ret)
}
