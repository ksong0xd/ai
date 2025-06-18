package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
)

// fetchDoc retrieves and parses HTML from a URL.
func fetchDoc(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// movingAverages computes simple moving averages for a stock symbol.
func movingAverages(symbol string, daysShort, daysLong int) (float64, float64, error) {
	end := time.Now()
	start := end.AddDate(0, 0, -daysLong-5)

	params := &chart.Params{
		Symbol:   symbol,
		Start:    datetime.New(&start),
		End:      datetime.New(&end),
		Interval: datetime.OneDay,
	}

	iter := chart.Get(params)
	var closes []float64
	for iter.Next() {
		bar := iter.Bar()
		f, _ := bar.Close.Float64()
		closes = append(closes, f)
	}
	if err := iter.Err(); err != nil {
		return 0, 0, err
	}
	if len(closes) < daysLong {
		return 0, 0, fmt.Errorf("insufficient data returned")
	}
	var shortSum, longSum float64
	for i := len(closes) - daysShort; i < len(closes); i++ {
		shortSum += closes[i]
	}
	for i := len(closes) - daysLong; i < len(closes); i++ {
		longSum += closes[i]
	}
	shortMA := shortSum / float64(daysShort)
	longMA := longSum / float64(daysLong)
	return shortMA, longMA, nil
}

func main() {
	symbol := "TSLA"

	// Scrape Tesla news headlines from CNBC
	cnbcDoc, err := fetchDoc("https://www.cnbc.com/quotes/TSLA?tab=news")
	if err != nil {
		log.Printf("Failed to get CNBC: %v", err)
	} else {
		fmt.Println("Recent news headlines:")
		cnbcDoc.Find("div.Card-titleContainer").Each(func(i int, s *goquery.Selection) {
			title := s.Text()
			fmt.Printf("- %s\n", title)
		})
		fmt.Println()
	}

	shortMA, longMA, err := movingAverages(symbol, 20, 50)
	if err != nil {
		log.Fatalf("Error retrieving price data: %v", err)
	}

	fmt.Printf("%s 20-day MA: %.2f\n", symbol, shortMA)
	fmt.Printf("%s 50-day MA: %.2f\n\n", symbol, longMA)

	if shortMA > longMA {
		fmt.Println("Short-term trend appears bullish relative to the 50-day average.")
	} else {
		fmt.Println("Short-term trend appears bearish relative to the 50-day average.")
	}

	fmt.Println("\n(Remember that this is a simplified example—consider more indicators for a real analysis.)")
	time.Sleep(1 * time.Second)
}
