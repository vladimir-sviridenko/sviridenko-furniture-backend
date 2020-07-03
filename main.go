package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"

	"github.com/labstack/gommon/log"
)

func getConfig() (*Config, error) {
	file, ok := os.LookupEnv("CONFIG_FILE")
	if !ok {
		file = "config.json"
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type EmailForm struct {
	Subject     string `json:"subject"`
	SendTo      string `json:"send_to"`
	HtmlMessage string `json:"html_message"`
}

type View struct {
	config *Config
}

func (v *View) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/send-email", v.SendEmail)
}

func (v *View) SendEmail(resp http.ResponseWriter, req *http.Request) {
	var form EmailForm
	if err := json.NewDecoder(req.Body).Decode(&form); err != nil {
		log.Warn(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	auth := smtp.PlainAuth(
		"",
		v.config.Email.Username,
		v.config.Email.Password,
		v.config.Email.Host,
	)
	var message bytes.Buffer
	message.WriteString(fmt.Sprintf("From: %q\r\n", v.config.Email.Username))
	message.WriteString(fmt.Sprintf("To: %q\r\n", form.SendTo))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", form.Subject))
	message.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n")
	message.WriteString(form.HtmlMessage)
	if err := smtp.SendMail(
		fmt.Sprintf("%s:%d", v.config.Email.Host, v.config.Email.Port),
		auth,
		v.config.Email.Username,
		[]string{form.SendTo},
		message.Bytes(),
	); err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	mux := http.ServeMux{}
	srv := http.Server{
		Addr:    cfg.Server.Addr,
		Handler: &mux,
	}
	view := View{config: cfg}
	view.Register(&mux)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
