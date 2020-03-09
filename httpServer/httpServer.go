package httpServer

import (
  "fmt"
  "net/http"
  "html/template"
)

var internalPD *PageData

type PageData struct {
Data map[string]float64
Author string
PageDescription string
Title string
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
  t.Execute(w, internalPD)
}

func StartServer(pd *PageData){
internalPD=pd
internalPD.Author = "Pedro Luna"
fmt.Println("Starting Server")
http.HandleFunc("/", Index)
http.HandleFunc("/firemonitor", FireMonitor)
http.ListenAndServe(":80", nil)
}
