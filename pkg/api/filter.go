package api

import (
	"regexp"
	"strings"
	"time"
	"whois-domain/pkg/entity"
)

type domainCheck struct {
	WhoisResultList   []entity.WhoisResult
	ExpireDomainList  []entity.WhoisResult
	UnknownDomainList []entity.WhoisResult
}

const (
	DATETARGETFORMAT = "2006-01-02"
	DATEORIGINFORMAT = "02-Jan-2006"
	DATEINFINITE     = "0000-00-00"
	DATEREGEXPATTERN = "[1-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]"
	AFTERMONTHS      = 4
)

var objectDomainCheck *domainCheck

func GetDomainFilterObject() *domainCheck {
	if objectDomainCheck == nil {
		objectDomainCheck = &domainCheck{
			ExpireDomainList:  make([]entity.WhoisResult, 0),
			UnknownDomainList: make([]entity.WhoisResult, 0),
		}
	}
	return objectDomainCheck
}

func (d *domainCheck) filtrateUnknownExpireDateDomain() {
	for _, v := range d.WhoisResultList {
		if v.ExpireDate == "" {
			d.UnknownDomainList = append(d.UnknownDomainList, v)
		}
	}
}

func (d *domainCheck) rewriteDateStr() {
	for i := range d.WhoisResultList {
		ds := d.cutDateStr(d.WhoisResultList[i].ExpireDate)
		if ds == DATEINFINITE {
			d.UnknownDomainList = append(d.UnknownDomainList, d.WhoisResultList[i])
			continue
		}
		t, err := time.ParseInLocation(DATEORIGINFORMAT, ds, time.Local)
		if err != nil {
			d.WhoisResultList[i].ExpireDate = ds
			continue
		}
		d.WhoisResultList[i].ExpireDate = t.Format(DATETARGETFORMAT)
	}
}

func (d *domainCheck) cutDateStr(s string) string {
	strSplit := strings.Split(s, " ")
	if s == "" {
		return s
	}
	if len(strSplit) > 1 {
		return strSplit[0]
	}
	strSplit = strings.Split(s, "T")
	if len(strSplit) > 1 {
		return strSplit[0]
	}
	r, _ := regexp.Compile(DATEREGEXPATTERN)
	if !r.MatchString(s) {
		return DATEINFINITE
	}
	return s
}

func (d *domainCheck) getLatelyExpireDomain(m int) {
	now := time.Now()
	getDate := now.AddDate(0, m, 0)
	nowDate := getDate.Format(DATETARGETFORMAT)
	for _, v := range d.WhoisResultList {
		t1, err1 := time.Parse(DATETARGETFORMAT, v.ExpireDate)
		t2, err2 := time.Parse(DATETARGETFORMAT, nowDate)
		if err1 == nil && err2 == nil && t1.Before(t2) {
			d.ExpireDomainList = append(d.ExpireDomainList, v)
		}
	}
}

func (d *domainCheck) GetExpireDomain(m int) {
	d.filtrateUnknownExpireDateDomain()
	d.rewriteDateStr()
	d.getLatelyExpireDomain(m)
}
