package scrape

import (
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
	c                = colly.NewCollector()

	urlList = []string{
		"https://www.realtor.com/realestateandhomes-search/%s/type-single-family-home/price-na-%d/pg-%d",
		"https://www.zillow.com/homes/Laramie,WY/",
	}
)

type Listing struct {
	Status  string
	Address string
	Bed     int
	Bath    int
	Sqft    int
	LotSize int
	Price   int
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

func Scrape() {
	ScrapeRDC(urlList[0])
}
