package processXml

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Urlset struct {
	XMLUrlSet xml.Name `xml:"urlset"`
	Urls      []Url    `xml:"url"`
}

type Url struct {
	Url xml.Name `xml:"url"`
	Loc string   `xml:"loc"`
}

func ReadSiteMap(sitekey int) (urlSet Urlset) {
	sitemap := "processxml/sitemap/sitemap-" + strconv.Itoa(sitekey) + ".xml"
	xmlFile, err := os.Open(sitemap)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.xml")
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(byteValue, &urlSet)

	return
}
