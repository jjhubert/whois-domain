package main

import (
	"fmt"
	"os"
	"whois-domain/pkg/api"
	"whois-domain/pkg/entity"
)

func main() {
	c := entity.GetConfigObject()
	queryClient := api.GetQueryDomainClientObject()
	queryClient.Config = c

	// 请求whois服务器获取域名过期时间
	r, err := queryClient.GetResult()
	if err != nil {
		os.Exit(1)
	}
	filterObject := api.GetDomainFilterObject()
	filterObject.WhoisResultList = r

	// 过滤出即将过期的域名
	filterObject.GetExpireDomain(c.Global.AfterMonths)

	// 发送信息
	n := api.GetNotificationObject()
	n.Config = c
	if len(filterObject.ExpireDomainList) > 0 {
		n.SendExpireMsg(filterObject.ExpireDomainList)
	}
	if len(filterObject.UnknownDomainList) > 0 {
		n.SendUnknownMsg(filterObject.UnknownDomainList)
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(fmt.Sprintf("Recovered in %s", r))
		}
	}()
}
