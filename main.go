package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	const Url = "https://www.skyscanner.net/transport/flights/mdea/syda/180801?adults=1&children=0&adultsv2=1&childrenv2=&infants=0&cabinclass=economy&rtn=0&preferdirects=false&outboundaltsenabled=false&inboundaltsenabled=false#results"

	fmt.Println("Collect start at: " + time.Now().Format(time.Stamp))
	cLinks := colly.NewCollector(
		colly.AllowedDomains("www.skyscanner.net"),
	)
	cLinks.Limit(&colly.LimitRule{
		DomainGlob:  "*www.skyscanner.net*",
		Delay:       4 * time.Second,
		RandomDelay: 2 * time.Second,
	})

	cLinks.OnHTML("a.CTASection__price-2bc7h.price", func(e *colly.HTMLElement) {
		fmt.Println("Price found: " + e.Text)
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
		fmt.Println(string(r.Body))
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

	cDetails := colly.NewCollector(
		colly.AllowedDomains("www.skyscanner.net"),
		colly.Async(true),
	)
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
	fmt.Println("Collect end at: " + time.Now().Format(time.Stamp))
}
