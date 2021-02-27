package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	if err := chromedp.Run(ctx, tasks()); err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Go's time.After example")
}

func tasks() chromedp.Tasks {
	return chromedp.Tasks{
		loadCookies(),
		chromedp.Navigate(`https://www.google.com/`),
		saveCookies(),
	}
}

func loadCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		if _, _err := os.Stat("cookies.json"); os.IsNotExist(_err) {
			return
		}

		cookiesData, err := ioutil.ReadFile("cookies.json")
		if err != nil {
			return
		}
		cookiesParams := network.SetCookiesParams{}
		if err = cookiesParams.UnmarshalJSON(cookiesData); err != nil {
			return
		}
		return network.SetCookies(cookiesParams.Cookies).Do(ctx)
	}
}

func saveCookies() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return
		}

		cookiesData, err := network.GetAllCookiesReturns{Cookies: cookies}.MarshalJSON()
		if err != nil {
			return
		}

		if err = ioutil.WriteFile("cookies.json", cookiesData, 0755); err != nil {
			return
		}
		return
	}
}
