package main

import (
	"encoding/xml"
	"fmt"
	"log"
)

func main() {
	xmlStr := `<Autodiscover xmlns="https://schemas.microsoft.com/exchange/autodiscover/outlook/requestschema/2006">
	<Request>
	  <EMailAddress>user@contoso.com</EMailAddress>
	  <AcceptableResponseSchema>https://schemas.microsoft.com/exchange/autodiscover/outlook/responseschema/2006a</AcceptableResponseSchema>
	</Request>
  	</Autodiscover>`
	var xmlRec struct {
		XMLName xml.Name `xml:"Autodiscover"`
		Request struct {
			XMLName                  xml.Name `xml:"Request"`
			EMailAddress             string   `xml:"EMailAddress"`
			AcceptableResponseSchema string   `xml:"AcceptableResponseSchema"`
		}
	}
	if err := xml.Unmarshal([]byte(xmlStr), &xmlRec); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s, %s\n", xmlRec.Request.EMailAddress, xmlRec.Request.AcceptableResponseSchema)
}
