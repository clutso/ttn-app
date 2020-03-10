package notificator

import (
  "net/smtp"
  "fmt"
  "os"
)

func SendMail(details string){

	// Sender data.
    from :=os.Getenv("NOTIFICATION_ADDRESS")
    password :=os.Getenv("NOTIFICATION_PASS")
    to := []string{
        "clutso@gmail.com",
        //"jsgutierrez@up.edu.mx",
    }    // smtp server configuration.
    smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
		// Message.
    //message:= []byte("To: jsgutierrez@up.edu.mx\r\n"+
		message:= []byte("To: clutso@gmail.com\r\n"+
			"Subject: Smells like fire\r\n\r\n"+
			"This is a notification sent from your IoT Fire monitor app.\r\n\r\n"+
			"You better take a look at the mote located in (lat, lon) "+ details)
		// Authentication.
    auth := smtp.PlainAuth("", from, password, smtpServer.host)
		// Sending email.
    err := smtp.SendMail(smtpServer.Address(), auth, from, to, message)
		    if err != nil {
        fmt.Println(err)
        return
    }
//		fmt.Println("Email Sent!")
}

func (s *smtpServer) Address() string {
 return s.host + ":" + s.port
}

type smtpServer struct {
 host string
 port string
}
