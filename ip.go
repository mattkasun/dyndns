package main

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"github.com/digitalocean/godo"
	"github.com/rdegges/go-ipify"
	"io/ioutil"
	"log"
	"log/syslog"
	"os"

)

type secrets struct {
	Token string
	Host string
	Id int
}

type ips struct {
	Ip string
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

func checkError (logger *syslog.Writer, err error, message string) {
	if err != nil {
		message = message + err.Error()
		logger.Err(message)
		log.Fatal(err)
	}
}

func main() {
	var secret secrets
	logger, err := syslog.New(syslog.LOG_ERR, "dyndns")
	if err != nil {
		log.Fatal(err)
	}
	
	//get current IP
	ip, err := ipify.GetIp()
	checkError(logger, err, "unable to get current IP: ")

	//read secrets
	file, err := ioutil.ReadFile(os.Getenv("HOME")+"/.config/dyndns")
	checkError(logger, err, "unable to read config: ")
	json.Unmarshal(file, &secret)

	//setup
	tokenSource := &TokenSource{
		AccessToken: secret.Token,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	ctx := context.TODO()

	//get current dns record
	record, _, err := client.Domains.Record(ctx, "nusak.ca", secret.Id)
	checkError(logger, err, "unable to retrieve dns record: ")

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
	_, _, err = client.Domains.EditRecord(ctx, "nusak.ca", secret.Id, editRequest)

	checkError(logger, err, "unable to update dns record: ")
	message := "updated ip address for " + secret.Host + " to " + ip
	logger.Info(message)
}
