package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
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

	res, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var ipJson IpJson
	json.Unmarshal(data, &ipJson)
	ip := ipJson.Query
	if ipJson.Status != "success" {
		log.Fatal("get ip failed")
	}

	client, err := alidns.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessSecret)

	describeDomainRequest := alidns.CreateDescribeDomainRecordsRequest()
	describeDomainRequest.Scheme = "https"

	describeDomainRequest.DomainName = config.DomainName

	domainRecords, err := client.DescribeDomainRecords(describeDomainRequest)
	if err != nil {
		fmt.Println(err.Error())
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
