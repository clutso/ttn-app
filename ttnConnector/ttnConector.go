package ttnConnector

import (
	"encoding/hex"
	"encoding/binary"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"os"
	notificator "github.com/clutso/ttn-app/notificator"
  httpServer "github.com/clutso/ttn-app/httpServer"
	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/apex"
	//"github.com/TheThingsNetwork/go-utils/random"
	"github.com/TheThingsNetwork/ttn/core/types"

)

const (
	sdkClientName = "ttn-requester"
)

func SendNotification(){
	//we still need to define when to send notifications.....
	notificator.SendMail("extra Details")
}

func DecodemoreComplexPayload(payload []byte )  map[string]float64  {
	m := make(map[string]float64)
	dataIndex:= [] int{}
	dataEntries :=int(len(payload)/3)
	//nextEntry=((x+1)*3)-1
	for x:=0; x< dataEntries;x++{
		dataIndex=append(dataIndex, 0)
	}

	fmt.Printf("Looking for %d entries\n",dataEntries)
	for x:=0; x< len(payload);x++{
			value:=int(payload[x])
			if value<= dataEntries {
				fmt.Printf("found %H in pos %d.\n",payload[x], x)
				dataIndex[value-1]=x
			}
		}
	fmt.Println(dataIndex)
	//for x:=0; x< dataEntries;x++{
	for x:=0; x< len(payload);x++{
		switch payload[x] {
			case 0x67:
				m["Temperature"]=0.
			case 0x68:
				m["Humidity"]=0.
			default:
				m["Unknown device"]=0.
		}
			fmt.Println(m)
	}
return m
}


//this particular decoder receives an array of 11 bins of uint8 numbers
//decoder SHOULD BE improved, this is only for testing porpouses
func decodePayload(payload []byte ) map[string]float64 {
	data := make(map[string]float64)

	//not so fancy decoder
	hSlice := []byte{0,payload[6]}
	humedad:=float64(binary.BigEndian.Uint16(hSlice))/2

	tSlice := []byte{payload[2],payload[3] }
	temperatura:=float64(binary.BigEndian.Uint16(tSlice))/10

	mySlice := []byte{payload[9],payload[10] }
	otro:=float64(binary.BigEndian.Uint16(mySlice))/100

	data["Temperature"]=temperatura
	data["Humidity"]=humedad
	data["Unknown device"]=otro

	return data
}

//fucntion to show output in console
	//¿deprecated?
	func PrintInConsole(payload []byte, data map[string]float64) {
		log := apex.Stdout() // We use a cli logger at Stdout
		ttnlog.Set(log)      // Set the logger as default for TTN
				hexPayload := hex.EncodeToString(payload)
				fmt.Println( "Received uplink with values:")
				log.WithField("data", hexPayload).Infof("%s: received uplink", sdkClientName)
				for entry, value := range data {
						fmt.Printf( " %v : %v \n", entry, value )
        }
    }

func GetLatLon(md types.Metadata)(float64, float64){
	rssi:=0.0
	snr:=0.0
	jsonMD, _:=json.Marshal(md)
	var dat map[string]interface{}
	if err := json.Unmarshal(jsonMD, &dat); err != nil {
				 fmt.Println(err)
		 }
	//	 fmt.Println(dat)
	strs:=dat["gateways"].([]interface{})
	for x:= range strs{
		gws:=strs[x].(map[string]interface{})
		rssi=gws["rssi"].(float64)
		snr=gws["snr"].(float64)

		//gwid:=gws["gtw_id"].(string)
/*
		locMD:=gws["LocationMetadata"]
		if locMD!= nil{
				fmt.Println("There you go!")
				fmt.Println(locMD)
				//need to substract and assing values

			}	else{
			fmt.Println("Uplink without useful info received")
			fmt.Println(dat)
		}
		//fmt.Println(gwid,rssi,snr)
*/
	}
	//	var dat1 map[string]interface{}
	//	if err := json.Unmarshal(gws, &dat1); err != nil {
	//				 panic(err)
	//		 }

	//.([]interface{})
	//for entry:= range gws{
	//	fmt.Println(string(entry))
	//}
return snr,rssi
	//return lat, lon
}



func StartConnector (pd *httpServer.PageData){
	log := apex.Stdout() // We use a cli logger at Stdout
	ttnlog.Set(log)      // Set the logger as default for TTN
	appID := os.Getenv("TTN_APP_ID")
	appAccessKey := os.Getenv("TTN_APP_ACCESS_KEY")
	//devToWatch:="pedro"
	devToWatch:="fire_monitoring"
	config := ttnsdk.NewCommunityConfig(sdkClientName)
	config.ClientVersion = "2.0.5" // The version of the application
	client := config.NewClient(appID, appAccessKey)
	defer client.Close()

	devices, err := client.ManageDevices()
	dev := new(ttnsdk.Device)
	if err != nil {
  	log.WithError(err).Fatalf("%s: could not read CA certificate file", sdkClientName)
	}

	dev, err = devices.Get(devToWatch)
	if err != nil {
		fmt.Println("Error getting Device")
  	log.WithError(err).Fatalf("%s: could not get device", sdkClientName)
	}
	fmt.Printf("got Device:%s from Application: %s and tne following description %s.\nAppEUI: %s DevEUI: %s\n",dev.DevID, dev.AppID,dev.Description,dev.AppEUI,dev.DevEUI)

	pubsub, err := client.PubSub()
	if err != nil {
	  log.WithError(err).Fatalf("%s: could not get application pub/sub", sdkClientName)
	}

	myNewDevicePubSub := pubsub.Device(devToWatch)

	uplink, err := myNewDevicePubSub.SubscribeUplink()
	if err != nil {
  	log.WithError(err).Fatalf("%s: could not subscribe to uplink messages", sdkClientName)
	}
	pd.Lat=0.0
	pd.Lon=0.0

	for message := range uplink {
		pd.Data= decodePayload (message.PayloadRaw)
		pd.Lat,pd.Lon  = GetLatLon(message.Metadata)
		//uncomment the following line to print in console
		//PrintInConsole(message.PayloadRaw, pd.Data)
		//¿call for index to refresh?
		//httpServer.Index()

  }

}
