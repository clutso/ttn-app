

function sayHI(){
  alert("HI")
}


function validateForm() {
if (document.forms[0]["Name"].value == "" || document.forms[0]["Email"].value == "" || document.forms[0]["Subject"].value == "" ||document.forms[0]["Message"].value == "") {
    alert("All fields must be filled out");
    return false;
  }
document.forms[0].submit()
alert("Email Sent");
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
