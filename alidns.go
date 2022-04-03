package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

type Record struct {
	Type string `json:"type"`
	RR   string `json:"RR"`
}

type Config struct {
	RegionId     string   `json:"regionId"`
	AccessKeyId  string   `json:"accessKeyId"`
	AccessSecret string   `json:"accessSecret"`
	DomainName   string   `json:"domainName"`
	Records      []Record `json:"records"`
}

type IpJson struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Query   string `json:"query"`
}

func getPubIp() (ip string) {
	res, err := http.Get("https://ifconfig.me")
	defer res.Body.Close()
	if err != nil || res.StatusCode != 200 {
		log.Fatal("http get error:", err, " http response:", res.StatusCode)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("read http body error:", err)
	}
	ip = string(data)
	fmt.Println("ip:", ip)
	return
}

func main() {

	configPath := flag.String("config", "./config.json", "path to config file")
	flag.Parse()
	fmt.Println("config file path:", *configPath)

	configJson, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	var config Config
	json.Unmarshal(configJson, &config)

	ip := getPubIp()

	client, err := alidns.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessSecret)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	describeDomainRequest := alidns.CreateDescribeDomainRecordsRequest()
	describeDomainRequest.Scheme = "https"

	describeDomainRequest.DomainName = config.DomainName

	domainRecords, err := client.DescribeDomainRecords(describeDomainRequest)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	var found bool = false

	// find subdomain
	for _, record := range config.Records {
		found = false
		for _, record_ := range domainRecords.DomainRecords.Record {

			if record.Type == record_.Type &&
				record.RR == record_.RR {
				found = true             //found subdomain record
				if ip == record_.Value { // ip not changed
					fmt.Println(record.RR+"."+config.DomainName, "unchanged:", ip)
				} else {
					fmt.Println(record.RR+"."+config.DomainName, "changed:", record_.Value, "->", ip)
					fmt.Println("updating record...")

					updateDomainRecordRequest := alidns.CreateUpdateDomainRecordRequest()
					updateDomainRecordRequest.Scheme = "https"
					updateDomainRecordRequest.RecordId = record_.RecordId
					updateDomainRecordRequest.RR = record.RR
					updateDomainRecordRequest.Type = record.Type
					updateDomainRecordRequest.Value = ip
					response, err := client.UpdateDomainRecord(updateDomainRecordRequest)
					if err != nil {
						fmt.Println(err.Error())
					}
					if response.IsSuccess() {
						fmt.Println("success")
					} else {
						fmt.Println("failed.", response.GetHttpContentString())
					}
				}
			}
		}

		if !found {
			fmt.Println(record.RR+"."+config.DomainName, "not found in records.")
			fmt.Println("adding record:", record.RR+"."+config.DomainName, "->", ip)

			addDomainRequest := alidns.CreateAddDomainRecordRequest()
			addDomainRequest.Scheme = "https"

			addDomainRequest.DomainName = config.DomainName
			addDomainRequest.RR = record.RR
			addDomainRequest.Type = record.Type
			addDomainRequest.Value = ip

			response, err := client.AddDomainRecord(addDomainRequest)
			if err != nil {
				fmt.Println(err.Error())
			}
			if response.IsSuccess() {
				fmt.Println("success")
			} else {
				fmt.Println("failed.", response.GetHttpContentString())
			}
		}
	}

	fmt.Println()
}
