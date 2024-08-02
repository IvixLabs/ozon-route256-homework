package inmemory

import (
	_ "embed"
	"encoding/json"
	"route256/loms/internal/model"
)

//go:embed stock-data.json
var stockData []byte

func InitStockData(storage *Storage) {
	stocks := make(map[model.Sku]model.Stock)

	rawStocks := make(map[model.Sku]model.Count)
	err := json.Unmarshal(stockData, &rawStocks)
	if err != nil {
		panic(err)
	}

	for k, v := range rawStocks {
		stockObj := model.NewStock(k, v)
		stocks[k] = *stockObj
	}

	storage.Stocks = stocks

}
