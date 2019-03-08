package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadAPIKey() (key, secret string) {
	homedir, err := os.UserHomeDir()
	check(err)
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/luno-key", homedir))
	check(err)

	lines := strings.Split(string(dat), "\n")

	return lines[0], lines[1]
}

func main() {

	ctx := context.Background()
	key, secret := ReadAPIKey()

	lunoClient := luno.NewClient()
	if err := lunoClient.SetAuth(key, secret); err != nil {
		check(err)
	}

	req := luno.GetTickersRequest{}

	resp, err := lunoClient.GetTickers(ctx, &req)
	check(err)

	for _, v := range resp.Tickers {
		time.Sleep(1 * time.Second)
		fmt.Println(fmt.Sprintf("%+v", v))

		PrintOrderBookForTicket(ctx, lunoClient, v.Pair)
	}

	time.Sleep(1 * time.Second)

	ExerciseQuoteAndReportBalance(ctx, lunoClient)
}

func ExerciseQuoteAndReportBalance(ctx context.Context, lunoClient *luno.Client) {
	createQuoteReq := luno.CreateQuoteRequest{
		Pair:       "EURXBT",
		Type:       "BUY",
		BaseAmount: decimal.NewFromInt64(1),
	}
	quoteResp, err := lunoClient.CreateQuote(ctx, &createQuoteReq)
	check(err)

	exerciseReq := luno.ExerciseQuoteRequest{Id: quoteResp.Id}
	_, err = lunoClient.ExerciseQuote(ctx, &exerciseReq)
	check(err)

	time.Sleep(1 * time.Second)
	balanceReq := luno.GetBalancesRequest{}
	balanceResp, err := lunoClient.GetBalances(ctx, &balanceReq)
	check(err)

	for _, v := range balanceResp.Balance {
		fmt.Println(fmt.Sprintf("%+v", v))
	}
}

func PrintOrderBookForTicket(ctx context.Context, lunoClient *luno.Client, pair string) {
	req := luno.GetOrderBookRequest{Pair: pair}

	resp, err := lunoClient.GetOrderBook(ctx, &req)
	check(err)

	fmt.Println("Asks:")
	for _, v := range resp.Asks {
		fmt.Println(fmt.Sprintf("%+v", v))
	}

	fmt.Println("\nBids:")
	for _, v := range resp.Bids {
		fmt.Println(fmt.Sprintf("%+v", v))
	}
}
