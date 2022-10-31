package restocks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/kompocikdot/labelbuddy/src/utils"
	"github.com/kompocikdot/labelbuddy/src/pdf"
)

type item struct {
	ShipTo string
	SellPrice string
	ItemName string
	ItemSize string
	Id string
}

type soldItem struct {
	Id string
	SoldDateString string
	SellPrice string
}

func login(email, password string, client *http.Client) error {
	token := getToken(client)
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
	defer res.Body.Close()
	bodyBytes, _ := io.ReadAll(res.Body)
	if utils.LoginErrRe.Match(bodyBytes) {
		return fmt.Errorf("can't login with given credentenials")
	}
	return nil
}

func RetrieveSalesLinks(email, password string) ([]item, []string, string, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	err := login(email, password, client)
	if err != nil {
		return nil, []string{}, "", err
	}

	var endOfPages bool
	var pageIndex uint8 = 1
	var items []item
	var itemIds []string

	for !endOfPages {
		url := fmt.Sprintf("https://restocks.net/en/account/listings/consignment?page=%d&search=", pageIndex)
		req, err := http.NewRequest(http.MethodGet, url, nil)
	
		if err != nil {
			log.Panic(err.Error())
		}
	
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")
		res, err := client.Do(req)
		if err != nil {
			log.Panic(err.Error())
		}
	
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)
		if endOfPages = utils.EndOfPagesRe.Match(resBody); endOfPages {
			break
		}
		
		strBody := string(resBody)

		allMatches := utils.LinksRe.FindAllStringSubmatch(strBody, -1)
		allShipBefore := utils.DatesRe.FindAllStringSubmatch(strBody, -1)
		allPrices := utils.ConsignPriceRe.FindAllStringSubmatch(strBody, -1)
		allShoeNames := utils.ItemNameRe.FindAllStringSubmatch(strBody, -1)
		allShoeSizes := utils.SizeRe.FindAllStringSubmatch(strBody, -1)
		ids := utils.ItemIdRe.FindAllStringSubmatch(strBody, -1)
		dates := utils.RegexStrDatesToDates(allShipBefore)


		
		for _, arr := range allMatches {
			itemIds = append(itemIds, arr[1])
		}
	
	

		for index := range allMatches {
			items = append(items, item{
				ShipTo: dates[index],
				SellPrice: allPrices[index][2],
				ItemName: allShoeNames[index][1],
				ItemSize: utils.ReplaceSizeUnicodesToString(allShoeSizes[index][1]),
				Id: ids[index][1],
			})
		}
		pageIndex++
	}

	downloadLabels(itemIds, client)
	emailDomain := strings.Split(email, "@")[0]
	pdfPathName := fmt.Sprintf("generated/%s.pdf", emailDomain)

	var fileNames []string
	for index, file := range itemIds {
		itemIds[index] = fmt.Sprintf("generated/%s.gif", file) 
		fileNames = append(fileNames, fmt.Sprintf("%s.gif", file))
	}
	
	pdf.GeneratePDF(itemIds, pdfPathName)
	fileNames = append(fileNames, emailDomain + ".pdf")
	return items, fileNames, (emailDomain + ".pdf"), nil
}

func getToken(client *http.Client) string {
	req, err := http.NewRequest(http.MethodGet, "https://restocks.net/en/login", nil)
	if err != nil {
		log.Panic(err.Error())
	}

	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")
	
	res, err := client.Do(req)
	if err != nil {
		log.Panic(err.Error())
	}

	defer res.Body.Close()
	byteBody, _ := io.ReadAll(res.Body)
	token := utils.TokenRe.FindStringSubmatch(string(byteBody))[1]
	return token
}

func RetrieveItemsPayments(email, password string) ([][]string, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	err := login(email, password, client)
	if err != nil {
		return nil, err
	}

	var endOfPages bool = false
	var pageIndex uint8 = 1

	var soldItems [][]string

	for !endOfPages {
		reqUrl := fmt.Sprintf("https://restocks.net/en/account/sales/history?page=%d&search=", pageIndex)
		req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		if err != nil {
			log.Panic(err.Error())
		}

		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")

		res, err := client.Do(req)
		if err != nil {
			log.Panic(err.Error())
		}

		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)
		endOfPages = utils.EndOfPagesRe.Match(resBody)
		
		if !endOfPages {
			stringified := string(resBody)
			dates := utils.SoldItemsDatesRe.FindAllStringSubmatch(stringified, -1)
			ids := utils.ItemIdRe.FindAllStringSubmatch(stringified, -1)
			prices := utils.PriceRe.FindAllStringSubmatch(stringified, -1)

			for index, datesSubArr := range(dates) {
				strDate := strings.Replace(datesSubArr[2], "\\/", "-", -1)
				time, _ := time.Parse("02-01-06", strDate)
				timeStr := time.Format("02.01.2006")
				
				item := []string{
					ids[index][1],
					timeStr,
					prices[index][2],
				}
				soldItems = append(soldItems, item)
			}
		}
		pageIndex++
	}

	newItem := [][]string{}
	newItem = append(newItem, []string{"id", "date", "price"})
	soldItems = append(newItem, soldItems...)

	return soldItems, nil
}