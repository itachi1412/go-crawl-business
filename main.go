package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	processXml "github.com/pistolbz/processxml"
	"github.com/pistolbz/src/utilities"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	client := &http.Client{}

	siteIndex := 1
	fmt.Printf("[+] Crawler sitemap-%v\n", siteIndex)
	start := time.Now()
	companies := utilities.NewCompanies()
	urlSet := processXml.ReadSiteMap(siteIndex)

	// Chạy Goroutines sử dụng semaphore để quản lý goroutines
	sem := semaphore.NewWeighted(8 * int64(runtime.NumCPU()))
	group, ctx := errgroup.WithContext(context.Background())
	for _, i := range urlSet.Urls {
		err := sem.Acquire(ctx, 1)
		if err != nil {
			fmt.Printf("Acquire err = %+v\n", err)
			continue
		}
		a := strings.Replace(i.Loc, " ", "", -1)
		if strings.Contains(a, "https://infodoanhnghiep.com/thong-tin") {
			group.Go(func() error {
				defer sem.Release(1)
				// work
				err := companies.ExtractInfomation(a, client)
				checkError(err)
				return nil
			})
		}

	}

	if err := group.Wait(); err != nil {
		fmt.Printf("g.Wait() err = %+v\n", err)
	}

	companiesJson, err := json.Marshal(companies)
	checkError(err)

	err = ioutil.WriteFile("result/companies-"+strconv.Itoa(siteIndex)+".json", companiesJson, 0644)
	checkError(err)
	fmt.Println("Crawler done! The result are saved in the file companies-" + strconv.Itoa(siteIndex) + ".json")
	fmt.Printf("Time running: %+v\n\n", time.Since(start))
}
