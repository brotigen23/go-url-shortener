// Code from Yandex.Practicum

package main

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"

	"math/rand"

	"github.com/brotigen23/go-url-shortener/internal/dto"
	"github.com/go-resty/resty/v2"
)

const target = "http://localhost:8080/"

var wg sync.WaitGroup

func newRandomString(size int) string {

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}

const urlCounts = 10000 / 2

func main() {

	//--------------------------------------------------
	// Register user
	//--------------------------------------------------
	client := resty.New()
	resp, err := client.R().
		Get(target)
	if err != nil {
		log.Println(err)
		return
	}
	client.SetCookie(resp.Cookies()[0])

	//--------------------------------------------------
	// CREATE URLS
	//--------------------------------------------------
	urls := []string{}
	aliases := []string{}

	log.Println("Create urls")
	for i := 0; i < urlCounts; i++ {
		urls = append(urls, newRandomString(8))
	}

	//--------------------------------------------------
	// POST URLs
	//--------------------------------------------------
	wg.Add(2)

	// /
	go func() {
		log.Println("Post urls to /")
		for i := 0; i < urlCounts/2; i++ {
			resp, err = client.R().
				SetBody(urls[i]).
				Post(target)
			aliases = append(aliases, string(resp.Body()))
		}
		wg.Done()
	}()

	// /api/shorten/batch
	go func() {
		log.Println("Post urls to /api/shorten/batch")
		for i := 0; i < urlCounts/2-1; i += 2 {
			body, err := json.Marshal([]dto.BatchRequest{
				{
					ID:  strconv.Itoa(i),
					URL: urls[i],
				},
				{
					ID:  strconv.Itoa(i + 1),
					URL: urls[i+1]},
			})
			if err != nil {
				log.Println(err)
				return
			}
			resp, err = client.R().
				SetBody(body).
				Post(target + "api/shorten/batch")
			if err != nil {
				log.Println(err)
				return
			}
		}
		wg.Done()
	}()
	wg.Wait()

	//--------------------------------------------------
	// GET URLs
	//--------------------------------------------------
	wg.Add(2)

	// /api/user/urls
	go func() {
		log.Println("Get urls from /api/user/urls")
		for i := 0; i < urlCounts; i++ {
			resp, err = client.R().
				Get(target + "api/user/urls")
			if err != nil {
				log.Println(err)
				return
			}
		}
		wg.Done()
	}()

	// /{id}
	go func() {
		log.Println("Get urls from /{id}")
		for i := 0; i < urlCounts/2; i++ {
			if err != nil {
				log.Println(err)
				return
			}
			resp, err = client.R().
				Get(aliases[i])
			if err != nil {
				log.Println(err)
				return
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
