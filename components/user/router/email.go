package router

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"regexp"
	"strconv"
	"time"
)

func GenEmailCode(digit int) string { //生成几位验证码

	AuthCode := ""
	seed := time.Now().Unix()
	rand.Seed(seed)
	for i := 0; i < digit; i++ {

		AuthCode += strconv.Itoa(rand.Intn(10))
	}
	return AuthCode

}
func EmailSend(to []string, title string, context string, from string) error {
	userEmail := "1273014435@qq.com"
	mailSmtpPort := ":587"
	mailPassword := "tfmbksrpxxfvhjig"
	mailHost := "smtp.qq.com"
	auth := smtp.PlainAuth("", userEmail, mailPassword, mailHost)
	for _, v := range to {
		if v != "" {
			header := make(map[string]string)
			header["From"] = from
			header["To"] = v
			header["Subject"] = title
			header["Content-Type"] = "text/html;charset=UTF-8"
			body := context
			to := []string{v}
			messageStr := ""
			for k, v := range header {
				messageStr += fmt.Sprintf("%s: %s\r\n", k, v)
			}
			messageStr += "\r\n" + body
			msg := []byte(messageStr)
			err := smtp.SendMail(mailHost+mailSmtpPort, auth, userEmail, to, msg)
			if err != nil {
				return err
			}
		}
	}

	return nil

}
func CheckEmail(targetEmail string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(targetEmail)
}
