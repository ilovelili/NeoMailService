package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"neo/config"
	"neo/core"
	"neo/dal"
	"neo/util"
	"time"
)

var (
	conf    *config.Config
	db      interface{}
	err     error
	pattern = flag.String("pattern", "neo", "pattern of mail service")
)

func init() {
	conf = config.GetConfig()
	db, err = dal.Open(conf)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	accessor := &core.RateAccessor{}

	switch *pattern {
	case "neo":
		err = neo(accessor)

	case "agentsmith":
		err = agentsmith(accessor)

	case "morpheus":
		err = morpheus(accessor)
	}

	if err != nil {
		panic(err)
	}
}

// neo normal mail notify (interval is 6 hours)
func neo(accessor *core.RateAccessor) (err error) {
	rates, err := accessor.GetLatestRatesInfo(db)
	if err != nil {
		return
	}

	firstrate := rates[0]
	secondrate := rates[1]

	body, err := json.Marshal(firstrate)
	if err != nil {
		return
	}

	var usdrate string
	var btcrate string
	if firstrate.Neo.Usd >= secondrate.Neo.Usd {
		usdrate = fmt.Sprintf("+%.4g", (firstrate.Neo.Usd/secondrate.Neo.Usd-1)*100)
	} else {
		usdrate = fmt.Sprintf("-%.4g", (secondrate.Neo.Usd/firstrate.Neo.Usd-1)*100)
	}

	if firstrate.Neo.Btc >= secondrate.Neo.Btc {
		btcrate = fmt.Sprintf("+%.4g", (firstrate.Neo.Btc/secondrate.Neo.Btc-1)*100)
	} else {
		btcrate = fmt.Sprintf("-%.4g", (secondrate.Neo.Btc/firstrate.Neo.Btc-1)*100)
	}

	subject := fmt.Sprintf("[Neo]: Neo price - $%g (%s%%) | %g BTC (%s%%)", firstrate.Neo.Usd, usdrate, firstrate.Neo.Btc, btcrate)
	return util.SendMail(conf, subject, "current exchange rate:"+string(body))
}

// agentsmith price changing notify (interval is 30 minutes)
func agentsmith(accessor *core.RateAccessor) (err error) {
	thresholdhit := false
	subject := "[AgentSmith]: "
	rates, err := accessor.GetLatestRatesInfo(db)
	if err != nil {
		return
	}

	firstrate := rates[0]
	secondrate := rates[1]

	if firstrate.Neo.Usd*conf.AgentSmith.Threshold <= secondrate.Neo.Usd {
		thresholdhit = true
		subject = subject + fmt.Sprintf("Neo to USD has decreased %.4g%%!", (secondrate.Neo.Usd/firstrate.Neo.Usd-1)*100)
	} else if firstrate.Neo.Usd >= secondrate.Neo.Usd*conf.AgentSmith.Threshold {
		thresholdhit = true
		subject = subject + fmt.Sprintf("Neo to USD has increased %.4g%%!", (firstrate.Neo.Usd/secondrate.Neo.Usd-1)*100)
	} else if firstrate.Neo.Btc*conf.AgentSmith.Threshold <= secondrate.Neo.Btc {
		thresholdhit = true
		subject = subject + fmt.Sprintf("Neo to BTC has decreased %.4g%%!", (secondrate.Neo.Btc/firstrate.Neo.Btc-1)*100)
	} else if firstrate.Neo.Btc >= secondrate.Neo.Btc*conf.AgentSmith.Threshold {
		thresholdhit = true
		subject = subject + fmt.Sprintf("Neo to BTC has increased %.4g%%!", (firstrate.Neo.Btc/secondrate.Neo.Btc-1)*100)
	} else if firstrate.Btc.Buy*conf.AgentSmith.Threshold <= secondrate.Btc.Buy {
		thresholdhit = true
		subject = subject + fmt.Sprintf("BTC buy price decreased %.4g%%!", (secondrate.Btc.Buy/firstrate.Btc.Buy-1)*100)
	} else if firstrate.Btc.Buy >= secondrate.Btc.Buy*conf.AgentSmith.Threshold {
		thresholdhit = true
		subject = subject + fmt.Sprintf("BTC buy price increased %.4g%%!", (firstrate.Btc.Buy/secondrate.Btc.Buy-1)*100)
	}

	if thresholdhit {
		body, err := json.Marshal(firstrate)
		if err != nil {
			return err
		}
		return util.SendMail(conf, subject, "current exchange rate:"+string(body))
	}

	return nil
}

// morpheus daily summary
func morpheus(accessor *core.RateAccessor) (err error) {
	neotousdmostexpensive, neotousdcheapest, neotobtcmostexpensive, neotobtccheapest, btcbuymostexpensive, btcbuycheapest, btcsellmostexpensive, btcsellcheapest, err := accessor.GetDailyRateSummary(db)
	if err != nil {
		return
	}

	subject := fmt.Sprintf("[Morpheus]: Daily Summary - %s", time.Now().Format("2006-01-02"))

	// of course we can use template instead of fmt
	body := fmt.Sprintln("======================================================================") +
		fmt.Sprintf("Neo to USD most expensive: %g at %s\n", neotousdmostexpensive.Usd, neotousdmostexpensive.Date.Format("2006-01-02 15:04")) +
		fmt.Sprintf("Neo to USD cheapest: %g at %s\n", neotousdcheapest.Usd, neotousdcheapest.Date.Format("2006-01-02 15:04")) +
		fmt.Sprintf("Neo to BTC most expensive: %g at %s\n", neotobtcmostexpensive.Btc, neotobtcmostexpensive.Date.Format("2006-01-02 15:04")) +
		fmt.Sprintf("Neo to BTC cheapest: %g at %s\n", neotobtccheapest.Btc, neotobtccheapest.Date.Format("2006-01-02 15:04")) +
		fmt.Sprintln("======================================================================") +
		fmt.Sprintf("BTC buy most expensive: %g at %s\n", btcbuymostexpensive.Buy, btcbuymostexpensive.Date.Format("2006-01-02 15:04")) +
		fmt.Sprintf("BTC buy cheapest: %g at %s\n", btcbuycheapest.Buy, btcbuycheapest.Date.Format("2006-01-02 15:04")) +
		fmt.Sprintf("BTC sell most expensive: %g at %s\n", btcsellmostexpensive.Sell, btcsellmostexpensive.Date.Format("2006-01-02 15:04")) +
		fmt.Sprintf("BTC sell cheapest: %g at %s\n", btcsellcheapest.Sell, btcsellcheapest.Date.Format("2006-01-02 15:04"))

	return util.SendMail(conf, subject, body)
}
