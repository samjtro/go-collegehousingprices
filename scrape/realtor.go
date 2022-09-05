package scrape

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
	c                = colly.NewCollector()
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

func ScrapeRDC() {
	for i := 1; i <= pagesFlag; i++ {
		url := fmt.Sprintf("https://www.realtor.com/realestateandhomes-search/%s/type-single-family-home/price-na-%d/pg-%d", cityStateFlag, priceIsUnderFlag, i)

		listings := ScrapeRDCHelper(url)
		results = append(results, listings)
	}

	for _, listings := range results {
		for _, listing := range listings {
			fmt.Println(listing.Status, listing.Address, listing.Price, listing.Bed, listing.Bath, listing.Sqft, listing.LotSize)
		}
	}
}

/*
for states:     'Colorado' or 'California'
for cities:     'city1-city2_state'   i.e.    'Los-Angeles_CA' or 'Denver_CO'
*/
func ScrapeRDCHelper(url string) Listings {
	var listings Listings

	c.OnHTML("li.jsx-1881802087.component_property-card", func(e *colly.HTMLElement) {
		var listing Listing

		listing.Price, _ = strconv.Atoi(RDCPriceStringToDigits(e.ChildText("span.rui__x3geed-0.kitA-dS")))
		listing.Status = e.ChildText("span.jsx-3853574337.statusText")
		listing.Address = e.ChildText("div.jsx-11645185.address.ellipsis.srp-page-address.srp-address-redesign") + " " + e.ChildText("div.jsx-11645185.address-second.ellipsis")

		bbsl := e.ChildText("ul.jsx-946479843.property-meta.list-unstyled.property-meta-srpPage")
		unorderedList := RDCBedBathSqftLotSize(bbsl)

		if len(unorderedList) >= 1 {
			listing.Bed = unorderedList[0]
		}

		if len(unorderedList) >= 2 {
			listing.Bath = unorderedList[1]
		}

		if len(unorderedList) >= 3 {
			listing.Sqft = unorderedList[2]
		}

		if len(unorderedList) >= 4 {
			listing.LotSize = unorderedList[3]
		}

		/*sqft, err := strconv.Atoi((e.ChildText("span.jsx-946479843.meta-value")))

		if err != nil {
			log.Fatal(err)
		}*/

		listings = append(listings, listing)
	})

	c.Visit(url)

	return listings
}

func RDCPriceStringToDigits(price string) string {
	var newInt []rune
	fmt.Println(price)

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
		} else {
			newInt = append(newInt, c)
		}
	}

	fmt.Println(string(newInt))

	return string(newInt)
}

// This function is broken atm; should return an arry of ints with the Bed, Bath, Sqft & LotSize but does not do so
func RDCBedBathSqftLotSize(in string) []int {
	var BedBathSqftLotSize []int

	for i, x := range in {
		if x == 'b' && in[i+1] == 'e' {
			beds, _ := strconv.Atoi(string(in[i-1]))
			BedBathSqftLotSize = append(BedBathSqftLotSize, beds)
		} else if x == 'b' && in[i+1] == 'a' {
			baths, _ := strconv.Atoi(string(in[i-1]))
			BedBathSqftLotSize = append(BedBathSqftLotSize, baths)
			sqft, _ := strconv.Atoi(string(in[i+4]))
			BedBathSqftLotSize = append(BedBathSqftLotSize, sqft)
		} else if x == 'l' {
			lotSize, _ := strconv.Atoi(string(in[i-6]))
			BedBathSqftLotSize = append(BedBathSqftLotSize, lotSize)
		}
	}

	fmt.Println(BedBathSqftLotSize)

	return BedBathSqftLotSize
}
