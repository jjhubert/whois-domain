package api

import (
	"context"
	"encoding/json"
	"github.com/whois-api-llc/whois-api-go"
	"io"
	"log"
	"net/http"
	"sync"
	"whois-domain/pkg/entity"
)

type queryDomainClient struct {
	Config *entity.Config
	sync.RWMutex
}

var objectQueryDomainClient *queryDomainClient

func (q *queryDomainClient) GetResult() ([]entity.WhoisResult, error) {
	resList := make([]entity.WhoisResult, 0)
	c := make(chan entity.WhoisResult)
	s := make(chan bool)
	for _, v := range q.Config.Domain {
		runFlag := false
		go func(v string, f bool) {
			r := &entity.WhoisResult{}
			if q.Config.Global.Ip2WhoisApiKey != "" && f == false {
				tmpRes, err := q.getResultByIp2Whois(v)
				if err != nil {
					log.Println(err)
					return
				}
				r = tmpRes
				f = true
			}

			if q.Config.Global.WhoisXmlApiKey != "" && f == false {
				tmpRes, err := q.getResultByWhoisXmlApi(v)
				if err != nil {
					log.Println(err)
					return
				}
				r = tmpRes
				f = true
			}

			if q.Config.Global.ApiLayerApiKey != "" {
				tmpRes, err := q.getResultByAPILayerWhoisApi(v)
				if err != nil {
					log.Println(err)
					return
				}
				r = tmpRes
				f = true
			}
			c <- *r
		}(v, runFlag)

	}
	go func() {
		i := 0
		for {
			if i == len(q.Config.Domain) {
				close(c)
				s <- true
				break
			}
			resList = append(resList, <-c)
			i++
		}
	}()
	<-s
	close(s)
	if q.Config.Manifest.DomainExpireList != nil {
		for k, v := range q.Config.Manifest.DomainExpireList {
			r := entity.WhoisResult{
				Domain:     k,
				ExpireDate: v,
			}
			resList = append(resList, r)
		}
	}
	return resList, nil
}

func (q *queryDomainClient) getResultByWhoisXmlApi(d string) (*entity.WhoisResult, error) {
	client := whoisapi.NewBasicClient(q.Config.Global.WhoisXmlApiKey)
	resp, _, err := client.WhoisService.Data(context.Background(), d)
	if err != nil {
		return nil, err
	}
	r := entity.WhoisResult{}
	r.Domain = d
	r.ExpireDate = resp.RegistryData.ExpiresDate
	r.CreateDate = resp.RegistryData.CreatedDate
	return &r, nil
}

func (q *queryDomainClient) getResultByIp2Whois(d string) (*entity.WhoisResult, error) {
	request, err := http.NewRequest("GET", q.Config.Global.Ip2WhoisUrl, nil)
	if err != nil {
		return nil, err
	}
	query := request.URL.Query()
	query.Set("key", q.Config.Global.Ip2WhoisApiKey)
	query.Set("format", "json")
	query.Set("domain", d)
	request.URL.RawQuery = query.Encode()
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		r := entity.WhoisResult{Domain: d, ExpireDate: ""}
		return &r, nil
	}
	r := entity.WhoisResult{}
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawData, &r)
	if err != nil {
		return nil, err
	}
	r.Domain = d
	return &r, nil
}

func (q *queryDomainClient) getResultByAPILayerWhoisApi(d string) (*entity.WhoisResult, error) {
	request, err := http.NewRequest("GET", q.Config.Global.ApiLayerWhoisUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("apikey", q.Config.Global.ApiLayerApiKey)
	query := request.URL.Query()
	query.Set("domain", d)
	request.URL.RawQuery = query.Encode()
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		r := entity.WhoisResult{Domain: d, ExpireDate: ""}
		return &r, nil
	}
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	qr := entity.ApiLayerQueryResult{}
	err = json.Unmarshal(rawData, &qr)
	if err != nil {
		return nil, err
	}
	r := entity.WhoisResult{}
	r.Domain = d
	r.ExpireDate = qr.Result.ExpirationDate
	r.CreateDate = qr.Result.CreationDate
	r.UpdateDate = qr.Result.UpdatedDate
	return &r, nil
}

func GetQueryDomainClientObject() *queryDomainClient {
	if objectQueryDomainClient == nil {
		objectQueryDomainClient = &queryDomainClient{}
	}
	return objectQueryDomainClient
}
