package project

import (
	"fmt"

	"github.com/mvrahden/go-enumer/pkg/utils"
)

func ExampleCountryCode() {
	cc := utils.Must(CountryCodeFromString("USA"))
	fmt.Println("Country Name:", cc.GetCountryName())
	fmt.Println("Country Code:", cc.GetCountryCode())
	fmt.Println("Area (in km^2):", cc.GetAreaInSquareKilometer())
	fmt.Println("GDP in USD:", cc.GetGdpInBillion()*1e12)
	fmt.Println("Population:", cc.GetPopulation())
	// output:
	// Country Name: United States
	// Country Code: 1
	// Area (in km^2): 9629091
	// GDP in USD: 1.672e+15
	// Population: 310232863
}
