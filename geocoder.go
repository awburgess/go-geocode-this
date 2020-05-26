package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

// Attribute struct represents the most nested unit for IN geocoder
type Attribute struct {
	ObjectID int    `json:"OBJECTID"`
	Street   string `json:"Street"`
	City     string `json:"City"`
	Zip      string `json:"ZIP"`
	State    string `json:"State"`
}

// Attributes struct is a key holder for Attribute
type Attributes struct {
	Attributes Attribute `json:"attributes"`
}

// Records struct is an array of Attributes
type Records struct {
	Records []Attributes `json:"records"`
}

// Configuration Config object for geocoder
var Configuration Config

func main() {

	Configuration = GenerateConfig()

	var wg sync.WaitGroup

	var allAddressSets Records

	var slicedAddresses []Records

	allAddressSets = ParseCSV()

	slicedAddresses = SliceRecords(allAddressSets)

	channel := make(chan Records)

	for i := 0; i < Configuration.ConcurrentRoutines; i++ {
		wg.Add(1)
		go geocodeWorker(channel, &wg)
	}

	for _, address := range slicedAddresses {
		channel <- address
	}
	close(channel)

	wg.Wait()
}

func geocodeWorker(ch <-chan Records, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range ch {
		geocode(data)
	}

}

func geocode(postBody Records) {

	requestBody, err := json.Marshal(postBody)
	if err != nil {
		panic(err)
	}

	response, err := http.PostForm(Configuration.GeocodeURL, url.Values{"addresses": {string(requestBody)}, "f": {"pjson"}, "outSR": {"4326"}})

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	uuid := uuid.New()
	name := uuid.String() + ".json"
	outFile := filepath.Join(Configuration.GeocodePath, name)

	out, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(out, response.Body)
}
