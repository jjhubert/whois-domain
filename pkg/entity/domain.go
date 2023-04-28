package entity

type ApiLayerQueryResult struct {
	Result ApiLayerWhoisResult `json:"result"`
}

type ApiLayerWhoisResult struct {
	DomainName     string `json:"domain_name"`
	WhoisServer    string `json:"whois_server"`
	UpdatedDate    string `json:"updated_date"`
	CreationDate   string `json:"creation_date"`
	ExpirationDate string `json:"expiration_date"`
}

type WhoisResult struct {
	Domain     string `json:"domain"`
	UpdateDate string `json:"update_date"`
	CreateDate string `json:"create_date"`
	ExpireDate string `json:"expire_date"`
}
