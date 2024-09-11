package shm_test

import (
	"encoding/json"
	"testing"

	"github.com/lianyun0502/shm"
	// "github.com/stretchr/testify/assert"
)


type TestOrderBook struct{
	Topic  string            `json:"Top"`
	Time   int64             `json:"T"`
	Symbol string            `json:"S"`
	Bids   map[string]string `json:"Bids"`
	Asks   map[string]string `json:"Asks"`
}

var testOrderBook = TestOrderBook{
	Topic:  "orderbook",
	Time:   1629782400,
	Symbol: "BTCUSDT",
	Bids:   map[string]string{
		"1": "100",
		"2": "200",
		"3": "300",
		"4": "400",
		"5": "500",
		"6": "600",
		"7": "700",
		"8": "800",
		"9": "900",
		"10": "1000",
		"11": "1100",
		"12": "1200",
		"13": "1300",
		"14": "1400",
		"15": "1500",
		"16": "1600",
		"17": "1700",
		"18": "1800",
		"19": "1900",
		"20": "2000",
		"21": "2100",
		"22": "2200",
		"23": "2300",
		"24": "2400",
		"25": "2500",
		"26": "2600",
		"27": "2700",
		"28": "2800",
		"29": "2900",
		"30": "3000",
	},
	Asks:   map[string]string{
		"1": "100",
		"2": "200",
		"3": "300",
		"4": "400",
		"5": "500",
		"6": "600",
		"7": "700",
		"8": "800",
		"9": "900",
		"10": "1000",
		"11": "1100",
		"12": "1200",
		"13": "1300",
		"14": "1400",
		"15": "1500",
		"16": "1600",
		"17": "1700",
		"18": "1800",
		"19": "1900",
		"20": "2000",
		"21": "2100",
		"22": "2200",
		"23": "2300",
		"24": "2400",
		"25": "2500",
		"26": "2600",
		"27": "2700",
		"28": "2800",
		"29": "2900",
		"30": "3000",
	},
}


func BenchmarkGetBinary(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bin, _ := shm.GetBinary(testOrderBook)
		_, _ = shm.GetStruct[TestOrderBook](bin)
	}
}

func BenchmarkJsonBinary(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bin, _ := json.Marshal(testOrderBook)
		_ = json.Unmarshal(bin, &TestOrderBook{})
	}
}