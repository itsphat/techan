package techan

import (
	"fmt"
	"strings"

	"github.com/sdcoffey/big"
)

// Candle represents basic market information for a security over a given time period
type Candle struct {
	Period     TimePeriod
	OpenPrice  big.Decimal
	ClosePrice big.Decimal
	MaxPrice   big.Decimal
	MinPrice   big.Decimal
	Volume     big.Decimal
	TradeCount uint
}

// NewCandle returns a new *Candle for a given time period
func NewCandle(period TimePeriod) (c *Candle) {
	return &Candle{
		Period:     period,
		OpenPrice:  big.ZERO,
		ClosePrice: big.ZERO,
		MaxPrice:   big.ZERO,
		MinPrice:   big.ZERO,
		Volume:     big.ZERO,
	}
}

// AddTrade adds a trade to this candle. It will determine if the current price is higher or lower than the min or max
// price and increment the tradecount.
func (c *Candle) AddTrade(tradeAmount, tradePrice big.Decimal) {
	if c.OpenPrice.Zero() {
		c.OpenPrice = tradePrice
	}
	c.ClosePrice = tradePrice

	if c.MaxPrice.Zero() {
		c.MaxPrice = tradePrice
	} else if tradePrice.GT(c.MaxPrice) {
		c.MaxPrice = tradePrice
	}

	if c.MinPrice.Zero() {
		c.MinPrice = tradePrice
	} else if tradePrice.LT(c.MinPrice) {
		c.MinPrice = tradePrice
	}

	if c.Volume.Zero() {
		c.Volume = tradeAmount
	} else {
		c.Volume = c.Volume.Add(tradeAmount)
	}

	c.TradeCount++
}

// UpdateCandle aggregates an existing candle with a new one, it's useful to sync a candle from multiple
// shorter time-period candles. For example, a 5-minute candle can be aggreated from 5 1-minute candles.
func (c *Candle) UpdateCandle(newCandle *Candle) {
	if newCandle == nil {
		return
	}
	if !(c.Period.Start.Before(newCandle.Period.Start) &&
		c.Period.End.After(newCandle.Period.Start)) {
		return
	}
	if c.MaxPrice.Zero() {
		c.MaxPrice = newCandle.MaxPrice
	} else if newCandle.MaxPrice.GT(c.MaxPrice) {
		c.MaxPrice = newCandle.MaxPrice
	}
	if c.MinPrice.Zero() {
		c.MinPrice = newCandle.MinPrice
	} else if newCandle.MinPrice.LT(c.MinPrice) {
		c.MinPrice = newCandle.MinPrice
	}
	c.ClosePrice = newCandle.ClosePrice
	c.Volume = c.Volume.Add(newCandle.Volume)
	c.TradeCount += newCandle.TradeCount
}

func (c *Candle) String() string {
	return strings.TrimSpace(fmt.Sprintf(
		`
Time:	%s
Open:	%s
Close:	%s
High:	%s
Low:	%s
Volume:	%s
	`,
		c.Period,
		c.OpenPrice.FormattedString(2),
		c.ClosePrice.FormattedString(2),
		c.MaxPrice.FormattedString(2),
		c.MinPrice.FormattedString(2),
		c.Volume.FormattedString(2),
	))
}
