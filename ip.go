package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"net/http"
	"time"
	"github.com/rdegges/go-ipify"
)

type secrets struct {
	Token string
	Host string
	Id string
}

type ips struct {
	Ip string
}

type record struct {
	Domain_Record struct {
		Name string
		Id int
		Data string
	}
}

func main() {
	var secret secrets
	var dns record
	logger, err := syslog.New(syslog.LOG_ERR, "dyndns")
	if err != nil {
		panic( err)
	}
	//get current IP
	ip, err := ipify.GetIp()
	if err != nil {
		logger.Err("unable to get current IP: ")
		panic (err)
	}


	//read secrets
	file, err := ioutil.ReadFile("/home/mkasun/go/src/dyndns/secrets")
	if err != nil {
		logger.Err("unable to read secrets")
		panic(err)
	}
	json.Unmarshal(file, &secret)

	//get current dns record
	url := "https://api.digitalocean.com/v2/domains/nusak.ca/records/"+secret.Id
	client := &http.Client{
		Timeout: time.Second *10,
	}
	req,_ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + secret.Token)
	response, err := client.Do(req)
	if err != nil {
		logger.Err("unable to retrieve dns record")
		panic (err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	json.Unmarshal(body, &dns)
	response.Body.Close()
	//json.Unmarshal(body,&dns)
	if dns.Domain_Record.Name != secret.Host {
		fmt.Println(dns)
		message := "wrong host name: DNS record=" + dns.Domain_Record.Name + " Host is " + secret.Host
		logger.Err(message)
		return
	}
	if dns.Domain_Record.Data == ip {
		logger.Info("IP address still the same, nothing to do")
		return
	} 

	//update dns record
	dns.Domain_Record.Data = ip
	b, _ := json.Marshal(dns)
	req,_ = http.NewRequest("PUT", url, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + secret.Token)


	response, err = client.Do(req)
	if err != nil {
		logger.Err("unable to update dns record")
		panic (err)
	}	
	message := "updated ip address for " + secret.Host + " to " + ip
	logger.Info(message)



}
