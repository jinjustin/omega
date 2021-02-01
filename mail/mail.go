package mail

import (
	"gopkg.in/gomail.v2"
)

//Send is a function that use to send email parameter contain reciever email and message
func Send(reciever string,subject string, message string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "omega.online.test.platform@gmail.com")
	m.SetAddressHeader("Reply-To", "noreply@gmail.com", "omega.online.test.platform@gmail.com")
	m.SetHeader("To", reciever)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)
	
	d := gomail.NewDialer("smtp.gmail.com", 587, "omega.online.test.platform@gmail.com", "gr]dfvdg9viN")
	
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}