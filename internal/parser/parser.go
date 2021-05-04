package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/dusansimic/yaas"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

func Parse(c *gin.Context, re yaas.RawEvent) (yaas.Event, error) {
	var e yaas.Event

	e.Time = time.Now()
	e.Code = re.Code

	_, path, err := parseUrl(re.URL)
	if err != nil {
		return yaas.Event{}, err
	}

	e.Path = path
	e.URL = re.URL

	dom, _, err := parseUrl((re.Referrer))
	if err != nil {
		return yaas.Event{}, err
	}
	e.Referrer = dom

	// cc, err := locateIP(c.ClientIP())
	cc, err := locateIP("188.2.20.138")
	if err != nil {
		return yaas.Event{}, err
	}

	e.CountryCode = cc

	ua := user_agent.New(c.Request.UserAgent())
	browser, _ := ua.Browser()

	e.Browser = yaas.ParseBrowser(browser)
	e.OS = yaas.GetOS(ua.OSInfo().Name)
	e.Device = yaas.ParseDevice(re.Width)

	return e, nil
}

func parseUrl(u string) (string, string, error) {
	re, err := regexp.Compile(`(http|https)://(.*?)($|\?.*|/.*)`)
	if err != nil {
		return "", "", err
	}

	// Get path from the url
	m := re.FindStringSubmatch(u)
	if len(m) < 3 {
		return "", "", nil
	} else if len(m) == 3 {
		return m[2], "/", nil
	} else {
		return m[2], m[3], nil
	}
}

type geoIPReps struct {
	CountryCode string `json:"country_code"`
}

func locateIP(ip string) (cc string, err error) {
	r, err := http.Get(fmt.Sprintf("https://geoip.fedoraproject.org/city?ip=%s", ip))
	if err != nil {
		return
	}
	defer func() {
		err = r.Body.Close()
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	var g geoIPReps
	if err = json.Unmarshal(b, &g); err != nil {
		return
	}

	return g.CountryCode, nil
}
