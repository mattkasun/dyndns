package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
)

func main() {
	logger, err := syslog.New(syslog.LOG_ERR, "ip lookup")
	if err != nil {
		log.Fatal(err)
	}
	res, _ := http.Get("https://api.ipify.org")
	ip, _ := ioutil.ReadAll(res.Body)
	ips := fmt.Sprintf("%s", ip)
	log.Println(ips)
	logger.Info(ips)
}
