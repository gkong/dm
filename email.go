// Copyright 2017 George S. Kong. All rights reserved.
// Use of this source code is governed by a license that can be found in the LICENSE.txt file.

package main

import (
	"net/smtp"
	"strings"
)

// NOTE: can NOT use raw string literals (back quotes) - drops carriage return

var htmlStart = "<!doctype html><html lang='en'><head><style type='text/css'>\r\n" +
	"body { font-family: verdana,tahoma,sans-serif,arial; }\r\n" +
	"a { background-color:#348eda; border-radius:5px; box-sizing:border-box;\r\n" +
	"color:#fff; cursor:pointer; text-decoration:none; font-size:14px;\r\n" +
	"font-weight:bold; margin:0; padding:20px;} </style></head><body><br/>\r\n"

var htmlEnd = "</body></html>\r\n"

func sendTextEmail(recipient string, subject string, textbody string) error {
	msg := []byte("From: " + Config.EmailSender + "\r\n" +
		"To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=utf-8; format=flowed\r\n" +
		"Content-Transfer-Encoding: 7bit\r\n" +
		"\r\n" +
		textbody)
	return sendEmail(recipient, msg)
}

func sendTextAndHTMLEmail(recipient string, subject string, textbody string, htmlbody string) error {
	msg := []byte("From: " + Config.EmailSender + "\r\n" +
		"To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Mime-Version: 1.0\r\n" +
		"Content-Type: multipart/alternative; boundary=\"cb75e46f89d377ffaf699aff\"\r\n" +
		"\r\n" +
		"--cb75e46f89d377ffaf699aff\r\n" +
		"Content-Type: text/plain; charset=\"ascii\"\r\n" +
		"Mime-Version: 1.0\r\n" +
		"Content-Transfer-Encoding: 7bit\r\n" +
		"\r\n" +
		textbody +
		"\r\n" +
		"--cb75e46f89d377ffaf699aff\r\n" +
		"Content-Type: text/html; charset=\"ascii\"\r\n" +
		"Mime-Version: 1.0\r\n" +
		"Content-Transfer-Encoding: 7bit\r\n" +
		"\r\n" +
		htmlbody +
		"\r\n" +
		"--cb75e46f89d377ffaf699aff--\r\n")
	return sendEmail(recipient, msg)
}

func sendEmail(recipient string, msg []byte) error {
	to := []string{recipient}
	auth := smtp.PlainAuth("", Config.EmailSender, Config.EmailPassword, Config.EmailServer)
	sp := Config.EmailServer + ":" + Config.EmailPort
	return smtp.SendMail(sp, auth, Config.EmailSender, to, msg)
}

func sendContactEmail(from, name, msg string) error {
	body := "  from: " + from + "\r\n\r\n  name: " + name + "\r\n\r\n" +
		strings.Replace(msg, "\n", "\r\n", -1) + "\r\n"
	return sendTextEmail(Config.AdminEmailAddr, "Delight/Meditate contact form submission", body)
}

func sendActivationEmail(recipient string, url string) error {
	subject := "Activate Your Delight/Meditate Account"

	tbody := "Please visit the following web page, " +
		"to activate your Delight/Meditate account:\r\n\r\n" +
		url + "\r\n\r\nThanks!\r\n\r\n"

	hbody := htmlStart + "<p>Click to activate your " +
		"Delight/Meditate account:</p><br/>\r\n" +
		"<a href='" + url + "'>Activate My Account</a>" + htmlEnd

	return sendTextAndHTMLEmail(recipient, subject, tbody, hbody)
}

func sendEChangeEmail(recipient string, url string) error {
	subject := "Confirm Email Address Change"

	tbody := "Please visit the following web page, to confirm " +
		"this email address for your Delight/Meditate account:\r\n\r\n" +
		url + "\r\n\r\nThanks!\r\n\r\n"

	hbody := htmlStart + "<p>Click to confirm this email address " +
		"for your Delight/Meditate account:</p><br/>\r\n" +
		"<a href='" + url + "'>Confirm</a>" + htmlEnd

	return sendTextAndHTMLEmail(recipient, subject, tbody, hbody)
}

func sendRecoveryEmail(recipient string, url string) error {
	subject := "Delight/Meditate Password Recovery"

	tbody := "Please visit the following web page, " +
		"to reset the password for your Delight/Meditate account:\r\n\r\n" +
		url + "\r\n\r\nThanks!\r\n\r\n"

	hbody := htmlStart + "<p>Click to reset the password for " +
		"your Delight/Meditate account:</p><br/>\r\n" +
		"<a href='" + url + "'>Reset Password</a></body></html>"

	return sendTextAndHTMLEmail(recipient, subject, tbody, hbody)
}
