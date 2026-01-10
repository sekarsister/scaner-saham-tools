package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type StockData struct {
	Symbol          string
	Name            string
	ClosePrice      float64
	ChangePercent   float64
	Volume          int64
	ForeignBuy      int64
	ForeignSell     int64
	NetForeignBuy   int64
	NetForeignValue float64
	ForeignPercent  float64
	Accumulation    int
	Score           float64
}

type ScanResult struct {
	Symbol          string
	Name            string
	Price           float64
	Change          float64
	NetForeignBuy   int64
	NetForeignValue float64
	ForeignPercent  float64
	Accumulation    int
	Strength        int
	Signal          string
}

var stockList = []struct {
	symbol string
	name   string
}{
	{"BBCA", "Bank Central Asia"},
	{"BBRI", "Bank Rakyat Indonesia"},
	{"BMRI", "Bank Mandiri"},
	{"TLKM", "Telkom Indonesia"},
	{"ASII", "Astra International"},
	{"UNVR", "Unilever Indonesia"},
	{"ICBP", "Indofood CBP"},
	{"GOTO", "GoTo Gojek Tokopedia"},
	{"BUKA", "Bukalapak"},
	{"ARTO", "Bank Jago"},
	{"EMTK", "Elang Mahkota"},
	{"MDKA", "Merdeka Copper Gold"},
	{"ANTM", "Aneka Tambang"},
	{"INCO", "Vale Indonesia"},
	{"PTBA", "Bukit Asam"},
	{"ADRO", "Adaro Energy"},
	{"ITMG", "Indo Tambangraya"},
	{"PGAS", "Perusahaan Gas Negara"},
	{"JSMR", "Jasa Marga"},
	{"CPIN", "Charoen Pokphand"},
	{"JPFA", "Japfa Comfeed"},
	{"ACES", "Ace Hardware"},
	{"ERAA", "Erajaya Swasembada"},
	{"MAPA", "MAP Aktif Adiperkasa"},
	{"SIDO", "Sido Muncul"},
	{"KLBF", "Kalbe Farma"},
	{"INDF", "Indofood Sukses"},
	{"GGRM", "Gudang Garam"},
	{"HMSP", "HM Sampoerna"},
	{"EXCL", "XL Axiata"},
	{"BRIS", "Bank Syariah Indonesia"},
	{"AMMN", "Amman Mineral"},
	{"BRPT", "Barito Pacific"},
	{"TPIA", "Chandra Asri"},
	{"SMGR", "Semen Indonesia"},
}

func generateStockData() []StockData {
	var stocks []StockData

	for _, s := range stockList {
		price := 500 + rand.Float64()*49500
		change := -5 + rand.Float64()*10
		volume := int64(1000000 + rand.Intn(99000000))

		foreignBuy := int64(float64(volume) * (0.1 + rand.Float64()*0.4))
		foreignSell := int64(float64(volume) * (0.1 + rand.Float64()*0.4))
		netFB := foreignBuy - foreignSell
		netFBValue := float64(netFB) * price

		foreignPct := float64(foreignBuy+foreignSell) / float64(volume) * 100

		accum := rand.Intn(10) - 3

		score := calculateScore(netFB, netFBValue, foreignPct, accum, change)

		stocks = append(stocks, StockData{
			Symbol:          s.symbol,
			Name:            s.name,
			ClosePrice:      price,
			ChangePercent:   change,
			Volume:          volume,
			ForeignBuy:      foreignBuy,
			ForeignSell:     foreignSell,
			NetForeignBuy:   netFB,
			NetForeignValue: netFBValue,
			ForeignPercent:  foreignPct,
			Accumulation:    accum,
			Score:           score,
		})
	}

	return stocks
}

func calculateScore(netFB int64, netFBValue, foreignPct float64, accum int, change float64) float64 {
	score := 0.0

	if netFB > 0 {
		score += 20
		if netFB > 5000000 {
			score += 15
		} else if netFB > 1000000 {
			score += 10
		}
	}

	if netFBValue > 50000000000 {
		score += 20
	} else if netFBValue > 10000000000 {
		score += 15
	} else if netFBValue > 1000000000 {
		score += 10
	}

	if foreignPct > 40 {
		score += 15
	} else if foreignPct > 25 {
		score += 10
	}

	if accum > 3 {
		score += 15
	} else if accum > 0 {
		score += 8
	}

	if change > 0 && change < 3 {
		score += 10
	}

	return math.Min(100, score)
}

func scanNetForeignBuy(stocks []StockData) []ScanResult {
	var results []ScanResult

	for _, stock := range stocks {
		if stock.NetForeignBuy > 0 && stock.Score >= 40 {
			strength := int(stock.Score / 20)
			if strength < 1 {
				strength = 1
			}
			if strength > 5 {
				strength = 5
			}

			signal := "ACCUMULATE"
			if stock.Score >= 70 {
				signal = "STRONG BUY"
			} else if stock.Score >= 55 {
				signal = "BUY"
			}

			results = append(results, ScanResult{
				Symbol:          stock.Symbol,
				Name:            stock.Name,
				Price:           stock.ClosePrice,
				Change:          stock.ChangePercent,
				NetForeignBuy:   stock.NetForeignBuy,
				NetForeignValue: stock.NetForeignValue,
				ForeignPercent:  stock.ForeignPercent,
				Accumulation:    stock.Accumulation,
				Strength:        strength,
				Signal:          signal,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].NetForeignValue > results[j].NetForeignValue
	})

	return results
}

func scanNetForeignSell(stocks []StockData) []ScanResult {
	var results []ScanResult

	for _, stock := range stocks {
		if stock.NetForeignBuy < -500000 {
			strength := 1
			netSell := -stock.NetForeignBuy
			if netSell > 10000000 {
				strength = 5
			} else if netSell > 5000000 {
				strength = 4
			} else if netSell > 2000000 {
				strength = 3
			} else if netSell > 1000000 {
				strength = 2
			}

			signal := "DISTRIBUTE"
			if strength >= 4 {
				signal = "STRONG SELL"
			} else if strength >= 3 {
				signal = "SELL"
			}

			results = append(results, ScanResult{
				Symbol:          stock.Symbol,
				Name:            stock.Name,
				Price:           stock.ClosePrice,
				Change:          stock.ChangePercent,
				NetForeignBuy:   stock.NetForeignBuy,
				NetForeignValue: stock.NetForeignValue,
				ForeignPercent:  stock.ForeignPercent,
				Accumulation:    stock.Accumulation,
				Strength:        strength,
				Signal:          signal,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].NetForeignValue < results[j].NetForeignValue
	})

	return results
}

func formatMoney(val float64) string {
	absVal := math.Abs(val)
	sign := ""
	if val < 0 {
		sign = "-"
	}

	if absVal >= 1000000000000 {
		return fmt.Sprintf("%sRp %.1fT", sign, absVal/1000000000000)
	} else if absVal >= 1000000000 {
		return fmt.Sprintf("%sRp %.1fB", sign, absVal/1000000000)
	} else if absVal >= 1000000 {
		return fmt.Sprintf("%sRp %.1fM", sign, absVal/1000000)
	}
	return fmt.Sprintf("%sRp %.0f", sign, absVal)
}

func formatVolume(vol int64) string {
	absVol := vol
	sign := ""
	if vol < 0 {
		absVol = -vol
		sign = "-"
	}

	if absVol >= 1000000 {
		return fmt.Sprintf("%s%.1fM", sign, float64(absVol)/1000000)
	} else if absVol >= 1000 {
		return fmt.Sprintf("%s%.1fK", sign, float64(absVol)/1000)
	}
	return fmt.Sprintf("%s%d", sign, absVol)
}

func getStars(n int) string {
	return strings.Repeat("*", n) + strings.Repeat(" ", 5-n)
}

func printHeader() {
	fmt.Print("\033[2J\033[H")
	fmt.Println(strings.Repeat("=", 95))
	fmt.Println("                      NET FOREIGN BUY SCANNER")
	fmt.Println("                    Indonesia Stock Exchange (IDX)")
	fmt.Println(strings.Repeat("=", 95))
	fmt.Printf(" Waktu: %s\n", time.Now().Format("02 Jan 2006 15:04:05"))
	fmt.Println(strings.Repeat("-", 95))
	fmt.Println()
}

func printNetForeignBuy(results []ScanResult) {
	fmt.Println("\033[1;32m                        NET FOREIGN BUY (Akumulasi Asing)\033[0m")
	fmt.Println(" Saham yang sedang diakumulasi oleh investor asing")
	fmt.Println(strings.Repeat("-", 95))

	if len(results) == 0 {
		fmt.Println(" Tidak ada saham dengan net foreign buy signifikan.")
		fmt.Println()
		return
	}

	fmt.Printf(" %-7s %-20s %-10s %-7s %-10s %-12s %-6s %-4s %-6s %-10s\n",
		"KODE", "NAMA", "HARGA", "CHG%", "NET FB", "VALUE", "F%", "ACC", "RATE", "SIGNAL")
	fmt.Println(strings.Repeat("-", 95))

	count := 0
	for _, r := range results {
		if count >= 12 {
			break
		}

		chgColor := "\033[32m"
		if r.Change < 0 {
			chgColor = "\033[31m"
		}

		name := r.Name
		if len(name) > 18 {
			name = name[:18]
		}

		fmt.Printf(" %-7s %-20s Rp%-8.0f %s%-6.1f%%\033[0m %-10s %-12s %-5.0f%% %-4d %-6s %-10s\n",
			r.Symbol,
			name,
			r.Price,
			chgColor,
			r.Change,
			formatVolume(r.NetForeignBuy),
			formatMoney(r.NetForeignValue),
			r.ForeignPercent,
			r.Accumulation,
			getStars(r.Strength),
			r.Signal)
		count++
	}
	fmt.Println()
}

func printNetForeignSell(results []ScanResult) {
	fmt.Println("\033[1;31m                       NET FOREIGN SELL (Distribusi Asing)\033[0m")
	fmt.Println(" Saham yang sedang dijual oleh investor asing")
	fmt.Println(strings.Repeat("-", 95))

	if len(results) == 0 {
		fmt.Println(" Tidak ada saham dengan net foreign sell signifikan.")
		fmt.Println()
		return
	}

	fmt.Printf(" %-7s %-20s %-10s %-7s %-10s %-12s %-6s %-4s %-6s %-10s\n",
		"KODE", "NAMA", "HARGA", "CHG%", "NET FS", "VALUE", "F%", "ACC", "RATE", "SIGNAL")
	fmt.Println(strings.Repeat("-", 95))

	count := 0
	for _, r := range results {
		if count >= 8 {
			break
		}

		chgColor := "\033[32m"
		if r.Change < 0 {
			chgColor = "\033[31m"
		}

		name := r.Name
		if len(name) > 18 {
			name = name[:18]
		}

		fmt.Printf(" %-7s %-20s Rp%-8.0f %s%-6.1f%%\033[0m %-10s %-12s %-5.0f%% %-4d %-6s %-10s\n",
			r.Symbol,
			name,
			r.Price,
			chgColor,
			r.Change,
			formatVolume(r.NetForeignBuy),
			formatMoney(r.NetForeignValue),
			r.ForeignPercent,
			r.Accumulation,
			getStars(r.Strength),
			r.Signal)
		count++
	}
	fmt.Println()
}

func printAllStocks(stocks []StockData) {
	printHeader()
	fmt.Println("                           DATA SEMUA SAHAM")
	fmt.Println(strings.Repeat("-", 95))

	fmt.Printf(" %-7s %-18s %-10s %-7s %-10s %-10s %-12s %-6s\n",
		"KODE", "NAMA", "HARGA", "CHG%", "FB", "FS", "NET VALUE", "SCORE")
	fmt.Println(strings.Repeat("-", 95))

	for _, s := range stocks {
		chgColor := "\033[32m"
		if s.ChangePercent < 0 {
			chgColor = "\033[31m"
		}

		name := s.Name
		if len(name) > 16 {
			name = name[:16]
		}

		fmt.Printf(" %-7s %-18s Rp%-8.0f %s%-6.1f%%\033[0m %-10s %-10s %-12s %-6.0f\n",
			s.Symbol,
			name,
			s.ClosePrice,
			chgColor,
			s.ChangePercent,
			formatVolume(s.ForeignBuy),
			formatVolume(s.ForeignSell),
			formatMoney(s.NetForeignValue),
			s.Score)
	}
	fmt.Println()
}

func printGuide() {
	printHeader()
	fmt.Println("                         PANDUAN NET FOREIGN BUY")
	fmt.Println(strings.Repeat("-", 95))
	fmt.Println()
	fmt.Println(" APA ITU NET FOREIGN BUY?")
	fmt.Println(" Net Foreign Buy adalah selisih antara pembelian dan penjualan")
	fmt.Println(" oleh investor asing. Jika positif, berarti asing lebih banyak")
	fmt.Println(" membeli (akumulasi). Jika negatif, berarti asing lebih banyak")
	fmt.Println(" menjual (distribusi).")
	fmt.Println()
	fmt.Println(" MENGAPA PENTING?")
	fmt.Println(" - Investor asing memiliki riset dan analisa mendalam")
	fmt.Println(" - Akumulasi asing sering mendahului kenaikan harga")
	fmt.Println(" - Distribusi asing bisa menjadi sinyal peringatan")
	fmt.Println()
	fmt.Println(" KRITERIA SCREENING:")
	fmt.Println(" - Net FB > 0 dengan nilai signifikan (> 1 Miliar)")
	fmt.Println(" - Foreign Percentage > 25% dari total volume")
	fmt.Println(" - Akumulasi berhari-hari (Accumulation Days > 3)")
	fmt.Println(" - Harga masih dalam tren naik moderat")
	fmt.Println()
	fmt.Println(" SIGNAL:")
	fmt.Println(" STRONG BUY  = Score >= 70, akumulasi sangat kuat")
	fmt.Println(" BUY         = Score >= 55, akumulasi cukup kuat")
	fmt.Println(" ACCUMULATE  = Score >= 40, mulai ada akumulasi")
	fmt.Println()
	fmt.Println(" PERINGATAN:")
	fmt.Println(" - Selalu kombinasikan dengan analisa teknikal")
	fmt.Println(" - Net foreign buy bukan jaminan harga naik")
	fmt.Println(" - Perhatikan juga kondisi market secara keseluruhan")
	fmt.Println()
}

func printMenu() {
	fmt.Println()
	fmt.Println(" MENU:")
	fmt.Println(" [1] Scan Ulang")
	fmt.Println(" [2] Lihat Semua Saham")
	fmt.Println(" [3] Panduan")
	fmt.Println(" [4] Keluar")
	fmt.Println()
	fmt.Print(" Pilihan: ")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		stocks := generateStockData()

		printHeader()

		buyResults := scanNetForeignBuy(stocks)
		sellResults := scanNetForeignSell(stocks)

		printNetForeignBuy(buyResults)
		printNetForeignSell(sellResults)

		printMenu()

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			continue
		case "2":
			printAllStocks(stocks)
			fmt.Println(" Tekan Enter...")
			fmt.Scanln()
		case "3":
			printGuide()
			fmt.Println(" Tekan Enter...")
			fmt.Scanln()
		case "4", "q", "Q":
			fmt.Println()
			fmt.Println(" Terima kasih!")
			fmt.Println(" Selamat berinvestasi.")
			fmt.Println()
			return
		}
	}
}
