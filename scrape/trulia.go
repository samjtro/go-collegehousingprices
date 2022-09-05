package scrape

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly"
)

var truliaResults []Listings

func ScrapeTrulia(url string) {
	for i := 1; i <= truliaPagesFlag; i++ {
		url := fmt.Sprintf(url, stateFlag, cityFlag, minPriceFlag, maxPriceFlag, i)

		listings := ScrapeTruliaHelper(url)
		truliaResults = append(truliaResults, listings)
	}

	for _, listings := range truliaResults {
		for _, listing := range listings {
			fmt.Println(listing.Status, listing.Address, listing.Price, listing.Bed, listing.Bath, listing.Sqft, listing.LotSize)
		}
	}
}

func ScrapeTruliaHelper(url string) Listings {
	var listings Listings

	c.OnHTML("div.PropertyCard__InteractivePropertyCardContainer-m1ur0x-2.CEYnN", func(e *colly.HTMLElement) {
		var listing Listing

		bbsq := TruliaUnmarshallBedBathSqft(e.ChildText("div.FlexContainers__Columns-zvngfq-2.eCjeDf"))

		if len(bbsq) >= 1 {
			listing.Bed = bbsq[0]
		} else if len(bbsq) >= 2 {
			listing.Bath = bbsq[1]
		} else if len(bbsq) >= 3 {
			listing.Sqft = bbsq[2]
		}

		listing.Price, _ = strconv.Atoi(UnformatPrice(e.ChildText("div.Text__TextBase-sc-1cait9d-0-div.Text__TextContainerBase-sc-1cait9d-1.keMYfJ")))

	})

	c.Visit(url)

	return listings
}

func TruliaUnmarshallBedBathSqft(in string) []int {
	var unmarshalledItems []int

	return unmarshalledItems
}
