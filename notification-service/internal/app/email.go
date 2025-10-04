package app

import (
    "fmt"
    "log"
    "net/smtp"
)

func SendEmail(to, msg string) {
    from := "socialmediashitebs@gmail.com"
    password := "ekvn aykp oihc uoui"

    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    body := "Subject: Notification\n\n" + msg
    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(body))
    if err != nil {
        log.Printf("Failed to send email: %v", err)
        return
    }
    fmt.Printf("✅ Email sent to %s\n", to)
}
