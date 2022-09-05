package scrape

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly"
)

var rdcResults []Listings

func ScrapeRDC(url string) {
	for i := 1; i <= rdcPagesFlag; i++ {
		url := fmt.Sprintf(url, cityStateFlag, priceIsUnderFlag, i)

		listings := ScrapeRDCHelper(url)
		rdcResults = append(rdcResults, listings)
	}

	for _, listings := range rdcResults {
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

		listing.Price, _ = strconv.Atoi(UnformatPrice(e.ChildText("span.rui__x3geed-0.kitA-dS")))
		listing.Status = e.ChildText("span.jsx-3853574337.statusText")
		listing.Address = e.ChildText("div.jsx-11645185.address.ellipsis.srp-page-address.srp-address-redesign") + " " + e.ChildText("div.jsx-11645185.address-second.ellipsis")

		bbsl := e.ChildText("ul.jsx-946479843.property-meta.list-unstyled.property-meta-srpPage")
		unorderedList := RDCUnmarshallPropertyMeta(bbsl)

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

// This function is broken atm; should return an arry of ints with the Bed, Bath, Sqft & LotSize but does not do so
func RDCUnmarshallPropertyMeta(in string) []int {
	var unmarshalledItems []int

	for i, x := range in {
		if x == 'b' && in[i+1] == 'e' {
			beds, _ := strconv.Atoi(string(in[i-1]))
			unmarshalledItems = append(unmarshalledItems, beds)
		} else if x == 'b' && in[i+1] == 'a' {
			baths, _ := strconv.Atoi(string(in[i-1]))
			unmarshalledItems = append(unmarshalledItems, baths)
			sqft, _ := strconv.Atoi(string(in[i+4]))
			unmarshalledItems = append(unmarshalledItems, sqft)
		} else if x == 'l' {
			lotSize, _ := strconv.Atoi(string(in[i-6]))
			unmarshalledItems = append(unmarshalledItems, lotSize)
		}
	}

	return unmarshalledItems
}
