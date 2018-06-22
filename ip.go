package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"log"
	//"log/syslog"
	"net/http"
	"time"
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
	var ip ips
	var secret secrets
	var dns record
	//logger, err := syslog.New(syslog.LOG_ERR, "ip lookup")
	//if err != nil {
	//	log.Fatal(err)
	//}
	res, _ := http.Get("https://api.ipify.org?format=json")
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &ip) 
	//ips := fmt.Sprintf("%s", ip)
	//log.Println(ips)
	//logger.Info(ips)
	fmt.Println(body,ip, ip.Ip)

	//read secrets
	file, err := ioutil.ReadFile("./secrets")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(file, &secret)
	fmt.Println(secret, secret.Token, secret.Host)

	//get current dns record
	client := &http.Client{
		Timeout: time.Second *10,
	}
	req,_ := http.NewRequest("GET","https://api.digitalocean.com/v2/domains/nusak.ca/records/"+secret.Id, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + secret.Token)
	fmt.Println(req)
	response, err := client.Do(req)
	if err != nil {
		panic (err)
	}
	body, _ = ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	json.Unmarshal(body,&dns)
	fmt.Println(dns, dns.Domain_Record.Name, dns.Domain_Record.Id, dns.Domain_Record.Data)
	if dns.Domain_Record.Data != ip.Ip {
		fmt.Println("IP address has changed, need to update it")
	} else {
		fmt.Println("IP address still the same, nothing to do")
	}

}
