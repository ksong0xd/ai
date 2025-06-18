# Stock Analysis Example

This repository contains a sample Go program that scrapes recent Tesla news from CNBC and computes simple moving averages for the Tesla (TSLA) stock price using the `finance-go` package.

The program outputs recent news headlines and indicates whether the short-term moving average is above or below the long-term average. It's intended as a minimal example and not as trading advice.

## Requirements

- Go 1.20 or higher
- Internet access to fetch news and price data

## Running

```
go run ./...
```

This will fetch the required modules, pull news from CNBC, and print a basic short- vs long-term outlook for TSLA.
