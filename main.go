package main

import (
	"encoding/hex"
	"encoding/binary"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"os"

	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/apex"
	//"github.com/TheThingsNetwork/go-utils/random"
	//"github.com/TheThingsNetwork/ttn/core/types"

)

const (
	sdkClientName = "my-amazing-app"
)

//this particular application receives an array of 11 bins of uint8 numbers
func decodePayload(payload []byte ) (float64,float64,float64) {

	//mySlice := []byte{,byte(payload[3]) }
	hSlice := []byte{0,payload[6] }
	humedad:=float64(binary.BigEndian.Uint16(hSlice))/2
	//binary.BigEndian.Uint32(mySlice)
	//humedad:=0.0
	tSlice := []byte{payload[2],payload[3] }
	temperatura:=float64(binary.BigEndian.Uint16(tSlice))/10
	mySlice := []byte{payload[9],payload[10] }
	otro:=float64(binary.BigEndian.Uint16(mySlice))/100
	return temperatura,humedad,  otro
}
func main (){

	log := apex.Stdout() // We use a cli logger at Stdout
	ttnlog.Set(log)      // Set the logger as default for TTN
	appID := os.Getenv("TTN_APP_ID")
	appAccessKey := os.Getenv("TTN_APP_ACCESS_KEY")
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
	fmt.Printf("got Device:%s from Application: %s and tne following description %s.\nAppEUI: %s DevEUI: %s\n",dev.DevID,dev.AppID,dev.Description, dev.AppEUI,dev.DevEUI)

	pubsub, err := client.PubSub()
	if err != nil {
	  log.WithError(err).Fatalf("%s: could not get application pub/sub", sdkClientName)
	}

	myNewDevicePubSub := pubsub.Device(devToWatch)

	uplink, err := myNewDevicePubSub.SubscribeUplink()
	if err != nil {
  	log.WithError(err).Fatalf("%s: could not subscribe to uplink messages", sdkClientName)
	}
	go func() {
  	for message := range uplink {
    	hexPayload := hex.EncodeToString(message.PayloadRaw)
			fmt.Println( "Received uplink with values:")
			temperatura, humedad,otro:= decodePayload (message.PayloadRaw)
			fmt.Printf("Temperatura: %2f, Humedad: %2f , no me acuerdo como se llama: %2f\n", temperatura, humedad,otro)
			log.WithField("data", hexPayload).Infof("%s: received uplink", sdkClientName)

  	}
	}()
for{

}
/*
	dev.AppID = appID
	dev.DevID = "fire_monitoring"
	dev.Description = "A new device in my amazing app"
	dev.AppEUI = types.AppEUI{0x70, 0xB3, 0xD5, 0x7E, 0xF0, 0x00, 0x00, 0x24} // Use the real AppEUI here
	dev.DevEUI = types.DevEUI{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08} // Use the real DevEUI here
	random.FillBytes(dev.AppKey[:])
*/

}
