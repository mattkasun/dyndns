package main

import (
	//"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"github.com/digitalocean/godo"
	"io/ioutil"
	"log/syslog"
	//"net/http"
	"github.com/rdegges/go-ipify"

)

type secrets struct {
	Token string
	Host string
	Id int
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

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token {
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func main() {
	var secret secrets
	//var dns record
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
	fmt.Println(ip)


	//read secrets
	file, err := ioutil.ReadFile("/home/mkasun/go/src/dyndns/secrets")
	if err != nil {
		logger.Err("unable to read secrets")
		panic(err)
	}
	json.Unmarshal(file, &secret)

	//setup
	tokenSource := &TokenSource{
		AccessToken: secret.Token,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	ctx := context.TODO()

	//get current dns record
	record, response, err := client.Domains.Record(ctx, "nusak.ca", secret.Id)
	fmt.Println(record, response)

	if err != nil {
		logger.Err("unable to retrieve dns record")
		panic (err)
	}
	if record.Name != secret.Host {
		message := "wrong host name: DNS record=" + record.Name + " Host is " + secret.Host
		logger.Err(message)
		return
	}
	if record.Data == ip {
		logger.Info("IP address still the same, nothing to do")	
		return
	} 

	//update dns record
	editRequest := &godo.DomainRecordEditRequest{
		Type: "A",
		Data: ip, 
	}
	updatedRecord, response, err := client.Domains.EditRecord(ctx, "nusak.ca", secret.Id, editRequest)

	if err != nil {
		logger.Err("unable to update dns record")
		panic (err)
	}	
	fmt.Println(updatedRecord, response)
	message := "updated ip address for " + secret.Host + " to " + ip
	logger.Info(message)
}
