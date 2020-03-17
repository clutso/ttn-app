package main
import (
	httpServer "github.com/clutso/ttn-app/httpServer"
	ttnConnector "github.com/clutso/ttn-app/ttnConnector"
	geolocator "github.com/clutso/ttn-app/geolocator"

)

func main (){

	var pageData httpServer.PageData
	var pd *httpServer.PageData
	pd= &pageData

	var geoRequest geolocator.InternalData
	var greq *geolocator.InternalData

	greq = &geoRequest

	go httpServer.StartServer(pd, greq)
	go ttnConnector.StartConnector(pd, greq)

for {}



}
