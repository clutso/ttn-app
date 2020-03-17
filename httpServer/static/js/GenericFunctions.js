

function sayHI(){
  alert("HI")
}

setInterval(function(){
  updtDash();
}, 10000);

function getLocation(){
  $.get("/updateLocation", function(data) {
    $('#currentLocation').html(data);
  });

}

function getMap(){
  $.get("/showmap", function(data) {
    $('#googleMap').html(data);
  });

}

function updtDash(){

  $.get('/updateDashboard', function(data) {
    $('#dashContent').html(data);
  });

}
