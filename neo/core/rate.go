package core

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Rate rate info
type Rate struct {
	Neo *Neo
	Btc *Btc
}

// Neo Neo info
type Neo struct {
	Btc  float32   `bson:"btc"`
	Usd  float32   `bson:"usd"`
	Date time.Time `bson:"date"`
}

// Btc Btc info
type Btc struct {
	Buy  float32   `json:"buy"`
	Sell float32   `json:"sell"`
	Date time.Time `bson:"date"`
}

// RateAccessor accessor
type RateAccessor struct {
}

// Description description of accessor
func (accessor *RateAccessor) Description() string {
	return "Neo and Btc info accessor"
}

// GetLatestRatesInfo list latest 2 rates info for comparison
func (accessor *RateAccessor) GetLatestRatesInfo(db interface{}) ([]*Rate, error) {
	var neo []*Neo
	var btc []*Btc

	collection := db.(*mgo.Database).C("rates")
	err := collection.Find(nil).Sort("-date").Limit(2).All(&neo)
	if err != nil {
		return nil, err
	}

	collection = db.(*mgo.Database).C("btcrates")
	err = collection.Find(nil).Sort("-date").Limit(2).All(&btc)
	if err != nil {
		return nil, err
	}

	return []*Rate{&Rate{
		Neo: neo[0],
		Btc: btc[0],
	}, &Rate{
		Neo: neo[1],
		Btc: btc[1],
	}}, nil
}

// GetDailyRateSummary get daily summary
func (accessor *RateAccessor) GetDailyRateSummary(db interface{}) (neotousdmostexpensive, neotousdcheapest, neotobtcmostexpensive, neotobtccheapest *Neo, btcbuymostexpensive, btcbuycheapest, btcsellmostexpensive, btcsellcheapest *Btc, err error) {
	from := time.Now().AddDate(0, 0, -1)
	to := time.Now()
	// neo
	collection := db.(*mgo.Database).C("rates")
	query := collection.Find(bson.M{
		"date": bson.M{
			"$gt": from,
			"$lt": to,
		},
	})
	if err = query.Sort("-usd").One(&neotousdmostexpensive); err != nil {
		return
	}
	if err = query.Sort("usd").One(&neotousdcheapest); err != nil {
		return
	}
	if err = query.Sort("-btc").One(&neotobtcmostexpensive); err != nil {
		return
	}
	if err = query.Sort("btc").One(&neotobtccheapest); err != nil {
		return
	}

	// btc
	collection = db.(*mgo.Database).C("btcrates")
	query = collection.Find(bson.M{
		"date": bson.M{
			"$gt": from,
			"$lt": to,
		},
	})
	if err = query.Sort("-buy").One(&btcbuymostexpensive); err != nil {
		return
	}
	if err = query.Sort("buy").One(&btcbuycheapest); err != nil {
		return
	}
	if err = query.Sort("-sell").One(&btcsellmostexpensive); err != nil {
		return
	}
	if err = query.Sort("sell").One(&btcsellcheapest); err != nil {
		return
	}

	return
}
