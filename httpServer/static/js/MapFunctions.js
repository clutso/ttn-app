var myLatlng;
var marker;
var map;
var mapOptions;
var lat;
var lng;

function myMap() {
  lat = parseFloat(document.getElementById('lat').textContent);
  lng = parseFloat(document.getElementById('lng').textContent);
  myLatlng=new google.maps.LatLng(lat,lng);


  mapOptions = {
      zoom: 17,
      center: myLatlng
    }

  map = new google.maps.Map(document.getElementById("googleMap"), mapOptions);
  marker = new google.maps.Marker({position: myLatlng,title:"Monitor here"});
  marker.setMap(map);
  }
