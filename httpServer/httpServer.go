package httpServer

import (
  "os"
  "fmt"
  "net/http"
  "net/smtp"
  "html/template"
)
var internalPD *PageData

type PageData struct {
Author string
PageDescription string
Title string
Uri string
Data map[string]float64
Lat float64
Lon float64
Snr float64
Rssi float64
}


//please remove this function or look for a better way to send emails...
func SendEmail(w http.ResponseWriter, r *http.Request){

  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  r.ParseForm()
  from :=os.Getenv("NOTIFICATION_ADDRESS")
  password :=os.Getenv("NOTIFICATION_PASS")
  host:="smtp.gmail.com"
  recipient:="clutso@gmail.com"
  subject:=r.Form["Subject"][0]
  name:= r.Form["Name"][0]
  email:=r.Form["Email"][0]
  message:=r.Form["Message"][0]

  to := []string{recipient}
  msg := []byte("To: clutso@hotmail.com\r\n" +
  		"Subject:"+ subject+"\r\n\r\n"+
  		"Hi pedro, "+name+" has sent you the following message:\r\n"+message+"\r\n\r\nContact him back to the following addres: "+email)
    auth := smtp.PlainAuth("",from, password, host)
    err := smtp.SendMail("smtp.gmail.com:587", auth, recipient, to, msg)
  	if err != nil {
  		fmt.Println(err)
  	}
    t, _ :=template.ParseFiles("./httpServer/static/index.html")
    internalPD.PageDescription="Peter's protofolio developed in go"
    internalPD.Title="Peter's way to go"
    t.Execute(w, internalPD)
  }

func Index (w http.ResponseWriter, r *http.Request){
  //t, _ :=template.ParseFiles("/home/clutso/go/src/github.com/clutso/ttn-app/httpServer/static/index.html")
  t, _ :=template.ParseFiles("./httpServer/static/index.html")
  internalPD.PageDescription="Peter's protofolio developed in go"
  internalPD.Title="Peter's way to go"
  t.Execute(w, internalPD)
}

func FireMonitor (w http.ResponseWriter, r *http.Request){
  //t, _ :=template.ParseFiles("/home/clutso/go/src/github.com/clutso/ttn-app/httpServer/static/fireMonitor.html")
  t, _ :=template.ParseFiles("./httpServer/static/fireMonitor.html")
  internalPD.PageDescription="LoRa Application to monitor the probability of fire in some location"
  internalPD.Title="Fire Monitoring"
//  internalPD.Uri="https://maps.googleapis.com/maps/api/js?key="+process.env.GOOGLE_API_KEY+"&callback=myMap"
  internalPD.Uri="https://maps.googleapis.com/maps/api/js?key="+os.Getenv("GOOGLE_API_KEY")+"&callback=myMap"
  t.Execute(w, internalPD)
}

func UpdateDash(w http.ResponseWriter, r *http.Request){
t, _ :=template.ParseFiles("./httpServer/static/dashboard.html")
t.Execute(w, internalPD)

}

func UpdateLocation (w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "text/html; charset=utf-8")

  fmt.Fprintf(w, "<span> Latitude: <span id=\"lat\"> %f </span> Longitude:  <span id=\"lng\"> %f </span>  </span>", internalPD.Lat, internalPD.Lon)
}

func StartServer(pd *PageData){
internalPD=pd
internalPD.Author = "Pedro Luna"
fmt.Println("Starting Server")
http.HandleFunc("/", Index)
http.HandleFunc("/firemonitor", FireMonitor)
http.HandleFunc("/SendEmail", SendEmail)
http.HandleFunc("/updateDashboard", UpdateDash)
http.HandleFunc("/updateLocation", UpdateLocation)

fs := http.FileServer(http.Dir("./httpServer/static/assets/"))
http.Handle("/assets/", http.StripPrefix("/assets/", fs))
fs1 := http.FileServer(http.Dir("./httpServer/static/documents/"))
http.Handle("/documents/", http.StripPrefix("/documents/", fs1))
http.ListenAndServe(":80", nil)
}
