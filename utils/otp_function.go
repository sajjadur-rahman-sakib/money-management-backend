package utils

import (
	"fmt"
	"math/rand"
	"money/config"
	"net/smtp"
	"strconv"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(otp)
}

func SendOTP(email, otp string) error {
	configuration := config.GetConfig()

	smtpHost := configuration.SmtpHost
	smtpPort := configuration.SmtpPort
	smtpUsername := configuration.SmtpUsername
	smtpPassword := configuration.SmtpPassword
	smtpEmail := configuration.SmtpEmail

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	to := []string{email}
	from := fmt.Sprintf("From: CashFlow <%s>\r\n", smtpEmail)
	toHeader := fmt.Sprintf("To: %s\r\n", email)
	subject := "Subject: OTP Verification\r\n"
	mime := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<title>OTP Verification</title>
</head>
<body style="margin:0; padding:0; background:#FFFFFF; font-family:Arial, Helvetica, sans-serif;">
	<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="background:#FFFFFF; padding:24px 12px;">
		<tr>
			<td align="center">
				<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="max-width:560px; background:#070B14; border:1px solid #0F1F45; border-radius:10px;">
					<tr>
						<td style="padding:48px 28px 40px; text-align:center;">
							<div style="font-size:64px; line-height:1; font-weight:700; letter-spacing:2px; color:#2B6CB8; margin-bottom:42px;">CashFlow</div>
							<div style="font-size:44px; line-height:1.2; font-weight:700; color:#FFFFFF; margin-bottom:14px;">OTP Verification</div>
							<div style="font-size:22px; line-height:1.45; color:#E6ECF7; max-width:430px; margin:0 auto 28px;">
								Use this OTP for further process.
							</div>
							<div style="background:#1A1E27; border-radius:10px; padding:22px 12px; margin:0 auto; max-width:430px;">
								<span style="display:inline-block; font-size:62px; line-height:1; letter-spacing:8px; font-weight:700; color:#FFFFFF;">%s</span>
							</div>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>
`, otp)
	msg := []byte(from + toHeader + subject + mime + "\r\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpEmail, to, msg)
	if err != nil {
		return err
	}

	return nil
}
