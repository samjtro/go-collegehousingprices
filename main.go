package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

var (
	cityStateFlag    string
	priceIsUnderFlag int
	pagesFlag        int
	results          []Listings
	urlList          = []string{}
	c                = colly.NewCollector()
)

type Listing struct {
	Status  string
	Address string
	Bed     string
	Bath    string
	Sqft    int
	LotSize int
	Price   float64
}

type Listings []Listing

func init() {
	err := godotenv.Load("config.env")

	if err != nil {
		log.Fatal(err)
	}

	cityStateFlag = os.Getenv("citystate")
	priceIsUnderFlag, err = strconv.Atoi(os.Getenv("priceisunder"))

	if err != nil {
		log.Fatal(err)
	}

	pagesFlag, err = strconv.Atoi(os.Getenv("pages"))

	if err != nil {
		log.Fatal(err)
	}

	c.OnError(func(_ *colly.Response, err error) {
		log.Fatal(err)
	})
}

func main() {
	urlList = append(urlList, fmt.Sprintf("https://www.realtor.com/realestateandhomes-search/%s/type-single-family-home/price-na-%d", cityStateFlag, priceIsUnderFlag))

	for _, url := range urlList {
		for i := 1; i <= pagesFlag; i++ {
			url += "/pg-%d"
			listings := ScrapeRDC(fmt.Sprintf(url, i))
			results = append(results, listings)
		}
	}

	for _, listings := range results {
		for _, listing := range listings {
			fmt.Println(listing.Status, listing.Address, listing.Bath, listing.Bed, listing.Price, listing.Sqft)
		}
	}
}

/*
for states:     'Colorado' or 'California'
for cities:     'city1-city2_state'   i.e.    'Los-Angeles_CA' or 'Denver_CO'
*/
func ScrapeRDC(url string) Listings {
	var listings Listings

	c.OnHTML("li.jsx-1881802087.component_property-card", func(e *colly.HTMLElement) {
		var listing Listing
		bedBathSqftLotSize := 0

		c.OnHTML("a.jsx-1534613990.card-anchor", func(e *colly.HTMLElement) {
			/*link := e.Attr("href")
			c.Visit(e.Request.AbsoluteURL(link))*/
		})

		c.OnHTML("span.jsx-3853574337.statusText", func(e *colly.HTMLElement) {
			listing.Status = e.Text
		})

		c.OnHTML("div.jsx-11645185.address.ellipsis.srp-page-address.srp-address-redesign", func(e *colly.HTMLElement) {
			listing.Address = e.Text
		})

		c.OnHTML("ul.jsx-946479843.property-meta.list-unstyled.property-meta-srpPage", func(e *colly.HTMLElement) {
			c.OnHTML("span.jsx-946479843.meta-value", func(e *colly.HTMLElement) {
				if bedBathSqftLotSize == 0 {
					listing.Bed = e.Text
				} else if bedBathSqftLotSize == 1 {
					listing.Bath = e.Text
				} else if bedBathSqftLotSize == 2 {
					sqft, err := strconv.Atoi(e.Text)

					if err != nil {
						log.Fatal(err)
					}

					listing.Sqft = sqft
				} else if bedBathSqftLotSize == 3 {
					lotSize, err := strconv.Atoi(e.Text)

					if err != nil {
						log.Fatal(err)
					}

					listing.LotSize = lotSize
				}

				bedBathSqftLotSize++
			})
		})

		c.OnHTML("span.rui_x3geed-0.kitA-dS", func(e *colly.HTMLElement) {
			price, err := strconv.ParseFloat(e.Text, 64)

			if err != nil {
				log.Fatal(err)
			}

			listing.Price = price
		})

		listings = append(listings, listing)
	})

	c.Visit(url)

	return listings
}
