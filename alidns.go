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

type Config struct {
	RegionId     string `json:"regionId"`
	AccessKeyId  string `json:"accessKeyId"`
	AccessSecret string `json:"accessSecret"`
	DomainName   string `json:"domainName"`
	RR           string `json:"RR"`
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
	fmt.Println("config file path: ", *configPath)

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

	var i int64 = 0
	var found bool = false
	// find subdomain
	for i = 0; i < domainRecords.TotalCount; i++ {
		if config.RR == domainRecords.DomainRecords.Record[i].RR {
			found = true                                           //found subdomain record
			if ip == domainRecords.DomainRecords.Record[i].Value { // ip not changed
				fmt.Println("ip", ip, "not changed!")
			} else {
				fmt.Printf("ip changed from %s to %s\n", domainRecords.DomainRecords.Record[i].Value, ip)
				fmt.Printf("---> update %s record to %s ...\n", config.RR+"."+config.DomainName, ip)

				updateDomainRecordRequest := alidns.CreateUpdateDomainRecordRequest()
				updateDomainRecordRequest.Scheme = "https"
				updateDomainRecordRequest.RecordId = domainRecords.DomainRecords.Record[i].RecordId
				updateDomainRecordRequest.RR = config.RR
				updateDomainRecordRequest.Type = "A"
				updateDomainRecordRequest.Value = ip
				response, err := client.UpdateDomainRecord(updateDomainRecordRequest)
				if err != nil {
					fmt.Println(err.Error())
				}
				if response.IsSuccess() {
					fmt.Println("success")
				} else {
					fmt.Printf("failed.\n%s\n", response.GetHttpContentString())
				}
			}
		}

	}
	if !found {
		fmt.Printf("ip is %s, but no domain record found.\n", ip)
		fmt.Printf("---> add %s record to %s ...\n", config.RR+"."+config.DomainName, ip)

		addDomainRequest := alidns.CreateAddDomainRecordRequest()
		addDomainRequest.Scheme = "https"

		addDomainRequest.DomainName = config.DomainName
		addDomainRequest.RR = config.RR
		addDomainRequest.Type = "A"
		addDomainRequest.Value = ip

		response, err := client.AddDomainRecord(addDomainRequest)
		if err != nil {
			fmt.Println(err.Error())
		}
		if response.IsSuccess() {
			fmt.Println("success")
		} else {
			fmt.Printf("failed.\n%s\n", response.GetHttpContentString())
		}
	}
	fmt.Println()
}
