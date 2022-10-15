package restocks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kompocikdot/labelbuddy/src/utils"
)

var client = utils.HttpClient

type item struct {
	Link string
	ShipTo string
}

func login(email, password string) error {
	token := getToken()
	values := url.Values{
		"_token": {token},
		"email": {email},
		"password": {password},
	}
	
	req, err := http.NewRequest(http.MethodPost, "https://restocks.net/en/login", strings.NewReader(values.Encode()))
	if err != nil {
		log.Panic(err.Error())
	}


	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	bodyBytes, _ := io.ReadAll(res.Body)
	if utils.LoginErrRe.Match(bodyBytes) {
		return fmt.Errorf("can't login with given credentenials")
	}
	return nil
}

func RetrieveSalesLinks(email, password string) ([]item, error) {
	err := login(email, password)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, "https://restocks.net/en/account/listings/consignment", nil)
	if err != nil {
		log.Panic(err.Error())
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	resBody, _ := io.ReadAll(res.Body)
	allMatches := utils.LinksRe.FindAllStringSubmatch(string(resBody), -1)
	allShipBefore := utils.DatesRe.FindAllStringSubmatch(string(resBody), -1)
	
	var dates []string
	for _, match := range allShipBefore {
		strDate := strings.Replace(match[2], "\\/", "-", -1)
		time, _ := time.Parse("02-01-06", strDate)
		timeStr := time.Format("02.01.2006")
		dates = append(dates, timeStr)
	}

	var items []item
	for index, match := range allMatches {
		link := fmt.Sprintf("https://restocks.net/en/account/sales/send-label/%s", match[1])
		items = append(items, item{Link: link, ShipTo: dates[index]})
	}

	return items, nil
}

func getToken() string {
	req, err := http.NewRequest(http.MethodGet, "https://restocks.net/en/login", nil)
	if err != nil {
		log.Panic(err.Error())
	}

	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")
	
	res, err := client.Do(req)
	if err != nil {
		log.Panic(err.Error())
	}

	byteBody, _ := io.ReadAll(res.Body)
	token := utils.TokenRe.FindStringSubmatch(string(byteBody))[1]
	return token


}