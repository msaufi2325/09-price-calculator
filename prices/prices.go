package prices

import (
	"fmt"

	"example.com/price-calculator/conversion"
	"example.com/price-calculator/iomanager"
)

func (job *TaxIncludedPriceJob) LoadData() error {

	lines, err := job.IOManager.ReadLines()
	if err != nil {
		return err
	}

	prices, err := conversion.StringsToFloat(lines)

	if err != nil {
		return err
	}

	job.InputPrices = prices
	return nil
}

type TaxIncludedPriceJob struct {
	IOManager         iomanager.IOManager `json:"-"`
	TaxRate           float64             `json:"tax_rate"`
	InputPrices       []float64           `json:"inputprices"`
	TaxIncludedPrices map[string]string   `json:"tax_included_prices"`
}

func (job *TaxIncludedPriceJob) Process(doneChan chan bool, errorChan chan error) {
	err := job.LoadData()

	if err != nil {
		fmt.Println(err)
		// return err
		errorChan <- err
		return
	}

	result := make(map[string]string)
	for _, price := range job.InputPrices {
		taxIncludedPrice := price * (1 + job.TaxRate)
		result[fmt.Sprintf("%.2f", price)] = fmt.Sprintf("%.2f", taxIncludedPrice)
	}

	job.TaxIncludedPrices = result
	job.IOManager.WriteResult(job)
	doneChan <- true
}

func NewTaxIncludedPriceJob(iom iomanager.IOManager, taxRate float64) *TaxIncludedPriceJob {
	return &TaxIncludedPriceJob{
		IOManager:   iom,
		InputPrices: []float64{10, 20, 30},
		TaxRate:     taxRate,
	}
}
