package restocks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func downloadAndSaveLabel(wg *sync.WaitGroup, idString string, client *http.Client) {
	defer wg.Done()
	link := fmt.Sprintf("https://restocks.net/en/account/sales/send-label/%s", idString)
	
	req, err := http.NewRequest(http.MethodGet, link, nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")

	if err != nil {
		log.Panic(err)
	}

	res, err :=	client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()
	
	
	out, err := os.Create(fmt.Sprintf("./generated/%s.gif", idString))
	if err != nil  {
		log.Panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil  {
		log.Panic(err)
	}

}

func downloadLabels(ids []string, client *http.Client) {
	var wg sync.WaitGroup

	for _, idString := range ids {
		wg.Add(1)	
		go downloadAndSaveLabel(&wg, idString, client)
	}
	wg.Wait()
}

func ClearFiles(names []string) {
	for _, name := range names {
		os.Remove(fmt.Sprintf("./generated/%s", name))
	}
}