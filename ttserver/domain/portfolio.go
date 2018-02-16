package domain

import (
	"fmt"
	"log"
	"sync"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

type portfolio struct {
	sync.RWMutex
	name     string
	isLive   bool
	balances []*BalanceAs
}

const DEFAULT_SYMBOL = SymbolType("BTC")

// initPortfolios - fetches latest balances from exchange
func (s *server) initPortfolios() error {

	s.log("Initialising portfolios")
	s.livePortfolio = &portfolio{
		name:   "LIVE",
		isLive: true,
	}
	if err := s.livePortfolio.refreshBalances(DEFAULT_SYMBOL); err != nil {
		return err
	}

	s.simPorts = make([]*portfolio, 0)
	return nil
}

// updatePortfolios - fetches latest balances and reprices
func (s *server) updatePortfolios() error {
	if err := s.livePortfolio.refreshBalances(DEFAULT_SYMBOL); err != nil {
		return err
	}
	return nil
}

// updateMetrics - sends metrics about portfolios to Influx
func (s *server) saveMetrics() error {

	s.RLock()
	defer s.RUnlock()

	// live metrics
	if err := s.livePortfolio.saveMetrics(); err != nil {
		return err
	}

	// simulated portfolio metrics
	for _, portfolio := range s.simPorts {
		if err := portfolio.saveMetrics(); err != nil {
			return err
		}

	}

	return nil
}

func (p *portfolio) saveMetrics() error {

	log.Printf("Sending portfolio balance data to influxdb")
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  DefaultMetrics.dbName,
		Precision: "ns",
	})

	if err != nil {
		return fmt.Errorf("failed to create batch points: %s", err)
	}

	portType := "live"
	if !p.isLive {
		portType = "simulated"
	}

	for _, balance := range p.balances {

		// Create a point and add to batch
		tags := map[string]string{"symbol": string(balance.Symbol), "name": p.name,
			"live": portType}
		fields := make(map[string]interface{}, 0)

		toSymbols := []SymbolType{SymbolType(BTC), SymbolType(ETH), SymbolType(USDT)}

		for _, toSym := range toSymbols {
			if symPrice, err := DefaultArchive.GetLatestPriceAs(SymbolType(balance.Symbol), toSym); err != nil {
				log.Printf("No %s price for %s symbol - %s", toSym, balance.Symbol, err)
			} else {
				fields[fmt.Sprintf("price.%s", toSym)] = symPrice.Price
				// calc current value = total * price
				value := symPrice.Price * balance.Total
				fields[fmt.Sprintf("value.%s", toSym)] = value
			}
		}

		if len(fields) > 0 {
			fields["exchange"] = balance.Exchange
			// add coin totals
			fields["total"] = balance.Total
			fields["locked"] = balance.Locked
			fields["free"] = balance.Free

			// only add fields with points
			pt, err := client.NewPoint("portfolio_balance", tags, fields, balance.At)
			if err != nil {
				fmt.Println("Error: ", err.Error())
			}

			bp.AddPoint(pt)
		}
	}
	// Write the batch
	return DefaultMetrics.Write(bp)

}

// refreshBalances - fetch latest balances from exchange
func (p *portfolio) refreshBalances(as SymbolType) error {
	fmt.Printf("Refreshing balances\n")
	p.Lock()
	defer p.Unlock()

	exchBalances, err := DefaultClient.GetBalances()
	if err != nil {
		return fmt.Errorf("failed to get balances from exchange: %s", err)
	}

	p.balances = make([]*BalanceAs, 0)

	// convert exchange balances to trada balances
	for _, exchBalance := range exchBalances {

		b := &BalanceAs{
			Balance: exchBalance,
			Total:   exchBalance.Free + exchBalance.Locked,
			At:      time.Now().UTC(),
			As:      as,
		}

		if err := b.reprice(); err != nil {
			return fmt.Errorf("failed reprice balance: %#v - %s", b, err)
		}

		p.balances = append(p.balances, b)
	}

	return nil
}

func (b *BalanceAs) reprice() error {
	// find latest price for trading pair
	priceAs, err := DefaultArchive.GetLatestPriceAs(SymbolType(b.Symbol), b.As)
	if err != nil {
		return fmt.Errorf("failed to get latest price for: %s as %s - %s", b.Symbol, b.As, err)
	}

	// reprice balance
	b.Price = priceAs.Price
	b.Value = priceAs.Price * b.Total

	// get 24h price

	daySummary, err := DefaultArchive.GetDaySummaryAs(SymbolType(b.Symbol), b.As)
	if err != nil {
		// no daily price info, but lets carry on
		//return fmt.Errorf("failed to get day summary for: %s as %s - %s", b.Symbol, b.As, err)
	} else {
		b.Price24H = daySummary.ClosePrice
		b.Value24H = b.Price24H * b.Total
		b.Change24H = daySummary.ChangePrice
		b.ChangePct24H = daySummary.ChangePercent
	}

	return nil
}
