package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	const Url = "https://www.skyscanner.net/transport/flights-from/mdea/?adults=1&children=0&adultsv2=1&childrenv2=&infants=0&cabinclass=economy&rtn=1&preferdirects=false&outboundaltsenabled=false&inboundaltsenabled=false&oym=1808&iym=1808&ref=home"

	fmt.Println("Collect start at: " + time.Now().Format(time.Stamp))
	cLinks := colly.NewCollector(
		colly.AllowedDomains("www.skyscanner.net"),
	)
	cDetails := colly.NewCollector(
		colly.AllowedDomains("www.skyscanner.net"),
		colly.Async(true),
	)
	cLinks.Limit(&colly.LimitRule{
		DomainGlob:  "*www.skyscanner.net*",
		Delay:       4 * time.Second,
		RandomDelay: 2 * time.Second,
	})
	cDetails.Limit(&colly.LimitRule{
		DomainGlob:  "*www.skyscanner.net*",
		Parallelism: 5,
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	})
	cDetails.OnHTML("li", func(e *colly.HTMLElement) {
		TextTitle := strings.TrimSpace(e.ChildText("b.col_50"))

		TextTitle = strings.TrimSpace(e.ChildText("a"))
		fmt.Println(TextTitle)
	})

	cDetails.OnScraped(func(r *colly.Response) {
		if strings.HasPrefix(r.Request.URL.String(), "http://www.arrendamientossantafe.com/webs/santafe/inmueble") {
			//fmt.Println("Finished ", r.Request.URL)
			fmt.Println("a")
		}
	})

	// On every a element which has href attribute call callback
	cLinks.OnHTML("a[href]", func(e *colly.HTMLElement) {
		/*link := e.Attr("href")
		if strings.HasPrefix(link, "/webs/santafe/inmueble") {
			cDetails.Visit(e.Request.AbsoluteURL(link))
		} else if strings.HasPrefix(link, "/webs/santafe/pages/basico") {
			cLinks.Visit(e.Request.AbsoluteURL(link))
		} else {
			return
		}*/
	})
	cLinks.OnResponse(func(r *colly.Response) {
		resp := string(r.Body)
		target := "United States"
		if strings.Contains(resp, target) {
			fmt.Println("Yes found " + target)
		}
		//fmt.Println(string(r.Body))
	})
	cLinks.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	cLinks.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		fmt.Println("RETRYING")
		cLinks.Visit(r.Request.URL.String())
	})
	cLinks.Visit(Url)

	fmt.Println("Collect end at: " + time.Now().Format(time.Stamp))
}
