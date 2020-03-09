package main
import (
	httpServer "github.com/clutso/ttn-app/httpServer"
	ttnConnector "github.com/clutso/ttn-app/ttnConnector"
)

func main (){
	var pageData httpServer.PageData
	var pd *httpServer.PageData
	pd= &pageData

	go httpServer.StartServer(pd)
	go ttnConnector.StartConnector(pd)

for {}

}
