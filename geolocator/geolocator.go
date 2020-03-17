package geolocator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"net/http"
	"bytes"
	"github.com/TheThingsNetwork/ttn/core/types"
)




type LocationEst struct  {
	Latitude float32 `json:"latitude",omitempty`
	Longitude float32 `json:"longitude",omitempty`
	ToleranceHoriz int32 `json:"toleranceHoriz",omitempty`
}

type GeoLocRespose struct{
	Result GeoLocResposeData `json:"result",omitempty`
}
type GeoLocResposeData struct {
  NumUsedGateways uint32 `json:"numUsedGateways",omitempty`
  HDOP float32 `json:"HDOP":,omitempty`
  algorithmType string `json:"algorithmType",omitempty`
	LocEst LocationEst `json:"locationEst",omitempty`

}

type InternalGateway struct {
	GatewayId string `json:"gatewayId"`
	Latitude float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Altitude int32 `json:"altitude",omitempty`
}

type InternalFrame struct {
	GatewayId string
	Antenna uint8
	TDOA uint32 //should be set to null... I Guess ...
	Rssi float32
	Snr float32
}

type InternalData struct {
  Gateways []InternalGateway
  Frames []InternalFrame
}
type Request struct {
    Gateways []InternalGateway `json:"gateways"`
    Frames [][]interface{}		`json:"frame"`
}




func GetFrameData (structMD types.Metadata )([]InternalGateway,[]InternalFrame){
	//drop the metadata on an appropiate strcuture
		var structGW []types.GatewayMetadata
		structGW= structMD.Gateways

		var gws []InternalGateway
		var tempGW InternalGateway
		var frames []InternalFrame
		var tempFrame InternalFrame

		for x:=range structGW{
			tempFrame.GatewayId=structGW[x].GtwID
			tempFrame.Antenna=structGW[x].Antenna
			tempFrame.TDOA=structGW[x].Timestamp
			tempFrame.Rssi=structGW[x].RSSI
			tempFrame.Snr=structGW[x].SNR
			frames=append(frames, tempFrame)

			tempGW.GatewayId=structGW[x].GtwID
			tempGW.Latitude=structGW[x].LocationMetadata.Latitude
			tempGW.Longitude=structGW[x].LocationMetadata.Longitude
			tempGW.Altitude=structGW[x].LocationMetadata.Altitude
			gws=append(gws, tempGW)


		}
return gws, frames
}

func SimLatLon(gw *InternalGateway){
switch gw.GatewayId{
	case "eui-7276ff00440101a7":
		gw.Latitude=21.822710
		gw.Longitude=-102.283705
		gw.Altitude=1
		//fmt.Println("Success :D")
	case "eui-7276ff0044010031":
		gw.Latitude=21.824142
		gw.Longitude=-102.283651
		gw.Altitude=1
		//fmt.Println("Success :D")
	case "eui-7276ff0044010051":
		gw.Latitude=21.824144
		gw.Longitude=-102.285204
		gw.Altitude=1
		//fmt.Println("Success :D")
	case "eui-7276ff00440101db":
		gw.Latitude=21.823264
		gw.Longitude=-102.284576
		gw.Altitude=1
		//fmt.Println("Success :D")
	default :
		fmt.Println("no GW declared.... unable to calc position")
	}
}

func RequestGeoloc(myGateWays []InternalGateway, intFrames []InternalFrame)(float32, float32){
	timeout:= time.Duration(5*time.Second)
	var myFrames[][]interface{}
	for x:= range intFrames{
		var myFrame[]interface{}
		myFrame=append(myFrame, intFrames[x].GatewayId)
		myFrame=append(myFrame, intFrames[x].Antenna)
		//myValue=nil //if having issues try this one
		myFrame=append(myFrame, intFrames[x].TDOA)
		myFrame=append(myFrame, intFrames[x].Rssi)
		myFrame=append(myFrame, intFrames[x].Snr)
		myFrames=append(myFrames, myFrame)
		}
	myBody:= Request{myGateWays,myFrames}
	requestBody, err := json.Marshal(myBody)
	//fmt.Println(string(requestBody))
	client:= http.Client{
		Timeout: timeout,
		}
	request, err:= http.NewRequest("POST", "https://gls.loracloud.com/api/v3/solve/singleframe", bytes.NewBuffer(requestBody))
	geoLocKey:=os.Getenv("GEOLOCATION_KEY")
	request.Header.Set("Ocp-Apim-Subscription-Key",geoLocKey)
	request.Header.Set("Accept", "application/json")
	resp, err:= client.Do(request)
	defer resp.Body.Close()
	body, err:= ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println(err)
		}

	fmt.Println(string (body))
	var response GeoLocRespose
	err =json.Unmarshal(body, &response)
	if err!=nil{
		fmt.Println (err)
		}

	return response.Result.LocEst.Latitude, response.Result.LocEst.Longitude

}
