package scrape

import (
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

var (
	//rdc flags
	cityStateFlag    string
	priceIsUnderFlag int
	rdcPagesFlag     int

	//trulia flags
	cityFlag        string
	stateFlag       string
	minPriceFlag    int
	maxPriceFlag    int
	truliaPagesFlag int

	c = colly.NewCollector()

	urlList = []string{
		"https://www.realtor.com/realestateandhomes-search/%s/type-single-family-home/price-na-%d/pg-%d", //CityState, Price, Page #
		"https://www.zillow.com/homes/%s/",              //CityState,
		"https://www.trulia.com/%s/%s/%d-%d_price/%d_p", //State, City, Min Price, Max Price, Page #
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

	cityStateFlag = os.Getenv("rdc-citystate")
	priceIsUnderFlag, err = strconv.Atoi(os.Getenv("rdc-price"))

	if err != nil {
		log.Fatal(err)
	}

	rdcPagesFlag, err = strconv.Atoi(os.Getenv("rdc-pages"))

	if err != nil {
		log.Fatal(err)
	}

	cityFlag = os.Getenv("trulia-city")
	stateFlag = os.Getenv("trulia-state")
	minPriceFlag, err = strconv.Atoi(os.Getenv("trulia-min"))

	if err != nil {
		log.Fatal(err)
	}

	maxPriceFlag, err = strconv.Atoi(os.Getenv("trulia-max"))

	if err != nil {
		log.Fatal(err)
	}

	truliaPagesFlag, err = strconv.Atoi(os.Getenv("trulia-pages"))

	if err != nil {
		log.Fatal(err)
	}

	c.OnError(func(_ *colly.Response, err error) {
		log.Fatal(err)
	})
}

func Scrape() {
	//ScrapeRDC(urlList[0])
	ScrapeTrulia(urlList[2])
}

func UnformatPrice(price string) string {
	var newInt []rune

	for i, c := range price {
		if c == '$' {
			if i > 0 {
				break
			}
		}

		if c == '$' {
			continue
		} else if c == ',' {
			continue
		} else if c == '+' {
			continue
		} else {
			newInt = append(newInt, c)
		}
	}

	return string(newInt)
}
