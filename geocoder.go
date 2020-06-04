package main

import (
	"encoding/csv"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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

// Response is JSON from geocoder
type Response struct {
	SpatialReference struct {
		Wkid       int `json:"wkid"`
		LatestWkid int `json:"latestWkid"`
	} `json:"spatialReference"`
	Locations []Location `json:"locations"`
}

// Location object in JSON response
type Location struct {
	Address  string `json:"address"`
	Location struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"location"`
	Score      int `json:"score"`
	Attributes struct {
		ResultID   int     `json:"ResultID"`
		LocName    string  `json:"Loc_name"`
		Status     string  `json:"Status"`
		Score      int     `json:"Score"`
		MatchAddr  string  `json:"Match_addr"`
		LongLabel  string  `json:"LongLabel"`
		ShortLabel string  `json:"ShortLabel"`
		AddrType   string  `json:"Addr_type"`
		Type       string  `json:"Type"`
		PlaceName  string  `json:"PlaceName"`
		PlaceAddr  string  `json:"Place_addr"`
		Phone      string  `json:"Phone"`
		URL        string  `json:"URL"`
		Rank       int     `json:"Rank"`
		AddBldg    string  `json:"AddBldg"`
		AddNum     string  `json:"AddNum"`
		AddNumFrom string  `json:"AddNumFrom"`
		AddNumTo   string  `json:"AddNumTo"`
		AddRange   string  `json:"AddRange"`
		Side       string  `json:"Side"`
		StPreDir   string  `json:"StPreDir"`
		StPreType  string  `json:"StPreType"`
		StName     string  `json:"StName"`
		StType     string  `json:"StType"`
		StDir      string  `json:"StDir"`
		BldgType   string  `json:"BldgType"`
		BldgName   string  `json:"BldgName"`
		LevelType  string  `json:"LevelType"`
		LevelName  string  `json:"LevelName"`
		UnitType   string  `json:"UnitType"`
		UnitName   string  `json:"UnitName"`
		SubAddr    string  `json:"SubAddr"`
		StAddr     string  `json:"StAddr"`
		Block      string  `json:"Block"`
		Sector     string  `json:"Sector"`
		Nbrhd      string  `json:"Nbrhd"`
		District   string  `json:"District"`
		City       string  `json:"City"`
		MetroArea  string  `json:"MetroArea"`
		Subregion  string  `json:"Subregion"`
		Region     string  `json:"Region"`
		RegionAbbr string  `json:"RegionAbbr"`
		Territory  string  `json:"Territory"`
		Zone       string  `json:"Zone"`
		Postal     string  `json:"Postal"`
		PostalExt  string  `json:"PostalExt"`
		Country    string  `json:"Country"`
		LangCode   string  `json:"LangCode"`
		Distance   int     `json:"Distance"`
		X          float64 `json:"X"`
		Y          float64 `json:"Y"`
		DisplayX   float64 `json:"DisplayX"`
		DisplayY   float64 `json:"DisplayY"`
		Xmin       float64 `json:"Xmin"`
		Xmax       float64 `json:"Xmax"`
		Ymin       float64 `json:"Ymin"`
		Ymax       float64 `json:"Ymax"`
		ExInfo     string  `json:"ExInfo"`
	} `json:"attributes"`
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

	concatJSON()
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

func concatJSON() {
	var resp Response

	files, err := ioutil.ReadDir(Configuration.GeocodePath)
	if err != nil {
		panic(err)
	}

	file := filepath.Join(Configuration.GeocodePath, "geocode_result.csv")

	_, err = os.Stat(file)
	if err == nil {
		os.Remove(file)
	}

	outfile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	resp.CSVHeader(outfile)

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			fullFile := filepath.Join(Configuration.GeocodePath, f.Name())
			byteValue, _ := ioutil.ReadFile(fullFile)
			json.Unmarshal(byteValue, &resp)

			resp.CSVRow(outfile)
		}
	}
}

// CSVHeader writes a header row
func (*Response) CSVHeader(w io.Writer) {
	cw := csv.NewWriter(w)
	cw.Write([]string{"resultid", "address", "cartesian_x", "cartesian_y", "locname", "status", "score", "matchaddr", "longlabel", "shortlabel", "addrtype", "type", "placename", "placeaddr", "phone", "url", "rank", "addbldg", "addnum", "addnumfrom", "addnumto", "addrange", "side", "stpredir", "stpretype", "stname", "sttype", "stdir", "bldgtype", "bldgname", "leveltype", "levelname", "unittype", "unitname", "subaddr", "staddr", "block", "sector", "nbrhd", "district", "city", "metroarea", "subregion", "region", "regionabbr", "territory", "zone", "postal", "postalext", "country", "langcode", "distance", "x", "y", "displayx", "displayy", "xmin", "xmax", "ymin", "ymax", "exinfo"})
	cw.Flush()
}

// CSVRow writes a row to csv
func (resp *Response) CSVRow(w io.Writer) {
	cw := csv.NewWriter(w)
	lc := resp.Locations
	for _, element := range lc {
		cw.Write([]string{strconv.Itoa(element.Attributes.ResultID), element.Address, strconv.FormatFloat(element.Location.X, 'f', -1, 64), strconv.FormatFloat(element.Location.Y, 'f', -1, 64), element.Attributes.LocName, element.Attributes.Status, strconv.Itoa(element.Attributes.Score), element.Attributes.MatchAddr, element.Attributes.LongLabel, element.Attributes.ShortLabel, element.Attributes.AddrType, element.Attributes.Type, element.Attributes.PlaceName, element.Attributes.PlaceAddr, element.Attributes.Phone, element.Attributes.URL, strconv.Itoa(element.Attributes.Rank), element.Attributes.AddBldg, element.Attributes.AddNum, element.Attributes.AddNumFrom, element.Attributes.AddNumTo, element.Attributes.AddRange, element.Attributes.Side, element.Attributes.StPreDir, element.Attributes.StPreType, element.Attributes.StName, element.Attributes.StType, element.Attributes.StDir, element.Attributes.BldgType, element.Attributes.BldgName, element.Attributes.LevelType, element.Attributes.LevelName, element.Attributes.UnitType, element.Attributes.UnitName, element.Attributes.SubAddr, element.Attributes.StAddr, element.Attributes.Block, element.Attributes.Sector, element.Attributes.Nbrhd, element.Attributes.District, element.Attributes.City, element.Attributes.MetroArea, element.Attributes.Subregion, element.Attributes.Region, element.Attributes.RegionAbbr, element.Attributes.Territory, element.Attributes.Zone, element.Attributes.Postal, element.Attributes.PostalExt, element.Attributes.Country, element.Attributes.LangCode, strconv.Itoa(element.Attributes.Distance), strconv.FormatFloat(element.Attributes.X, 'f', -1, 64), strconv.FormatFloat(element.Attributes.Y, 'f', -1, 64), strconv.FormatFloat(element.Attributes.DisplayX, 'f', -1, 64), strconv.FormatFloat(element.Attributes.DisplayY, 'f', -1, 64), strconv.FormatFloat(element.Attributes.Xmin, 'f', -1, 64), strconv.FormatFloat(element.Attributes.Xmax, 'f', -1, 64), strconv.FormatFloat(element.Attributes.Ymin, 'f', -1, 64), strconv.FormatFloat(element.Attributes.Ymax, 'f', -1, 64), element.Attributes.ExInfo})
	}
	cw.Flush()
}
