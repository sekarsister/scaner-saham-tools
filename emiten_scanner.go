package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Emiten struct {
	Symbol        string
	Name          string
	Sector        string
	Price         float64
	Open          float64
	High          float64
	Low           float64
	PrevClose     float64
	Change        float64
	Volume        int64
	AvgVolume     int64
	RSI           float64
	MACD          float64
	GapPercent    float64
	Volatility    float64
	MorningMoment float64
	AfternoonDip  float64
	ScoreBSJP     float64
	ScoreBPJS     float64
}

type ScanResult struct {
	Symbol   string
	Name     string
	Sector   string
	Price    float64
	Change   float64
	Target   float64
	StopLoss float64
	Score    float64
	Signal   string
	Reason   string
}

var emitenList = []struct {
	symbol string
	name   string
	sector string
}{
	{"BBCA", "Bank Central Asia", "Banking"},
	{"BBRI", "Bank Rakyat Indonesia", "Banking"},
	{"BMRI", "Bank Mandiri", "Banking"},
	{"BBNI", "Bank Negara Indonesia", "Banking"},
	{"BRIS", "Bank Syariah Indonesia", "Banking"},
	{"ARTO", "Bank Jago", "Banking"},
	{"BTPS", "Bank BTPN Syariah", "Banking"},
	{"MEGA", "Bank Mega", "Banking"},
	{"NISP", "Bank OCBC NISP", "Banking"},
	{"BNGA", "Bank CIMB Niaga", "Banking"},
	{"TLKM", "Telkom Indonesia", "Telecom"},
	{"EXCL", "XL Axiata", "Telecom"},
	{"ISAT", "Indosat Ooredoo", "Telecom"},
	{"FREN", "Smartfren Telecom", "Telecom"},
	{"ASII", "Astra International", "Automotive"},
	{"AUTO", "Astra Otoparts", "Automotive"},
	{"SMSM", "Selamat Sempurna", "Automotive"},
	{"UNVR", "Unilever Indonesia", "Consumer"},
	{"ICBP", "Indofood CBP", "Consumer"},
	{"INDF", "Indofood Sukses", "Consumer"},
	{"MYOR", "Mayora Indah", "Consumer"},
	{"KLBF", "Kalbe Farma", "Healthcare"},
	{"SIDO", "Sido Muncul", "Healthcare"},
	{"DVLA", "Darya Varia", "Healthcare"},
	{"KAEF", "Kimia Farma", "Healthcare"},
	{"PYFA", "Pyridam Farma", "Healthcare"},
	{"GGRM", "Gudang Garam", "Tobacco"},
	{"HMSP", "HM Sampoerna", "Tobacco"},
	{"GOTO", "GoTo Gojek Tokopedia", "Technology"},
	{"BUKA", "Bukalapak", "Technology"},
	{"EMTK", "Elang Mahkota", "Technology"},
	{"MTDL", "Metrodata Electronics", "Technology"},
	{"ANTM", "Aneka Tambang", "Mining"},
	{"INCO", "Vale Indonesia", "Mining"},
	{"PTBA", "Bukit Asam", "Mining"},
	{"ADRO", "Adaro Energy", "Mining"},
	{"ITMG", "Indo Tambangraya", "Mining"},
	{"MDKA", "Merdeka Copper Gold", "Mining"},
	{"AMMN", "Amman Mineral", "Mining"},
	{"TINS", "Timah", "Mining"},
	{"MEDC", "Medco Energi", "Energy"},
	{"PGAS", "Perusahaan Gas Negara", "Energy"},
	{"AKRA", "AKR Corporindo", "Energy"},
	{"JSMR", "Jasa Marga", "Infrastructure"},
	{"WIKA", "Wijaya Karya", "Infrastructure"},
	{"PTPP", "PP Persero", "Infrastructure"},
	{"WSKT", "Waskita Karya", "Infrastructure"},
	{"CPIN", "Charoen Pokphand", "Poultry"},
	{"JPFA", "Japfa Comfeed", "Poultry"},
	{"MAIN", "Malindo Feedmill", "Poultry"},
	{"ACES", "Ace Hardware", "Retail"},
	{"ERAA", "Erajaya Swasembada", "Retail"},
	{"MAPI", "Mitra Adiperkasa", "Retail"},
	{"LPPF", "Matahari Dept Store", "Retail"},
	{"RALS", "Ramayana Lestari", "Retail"},
	{"SMGR", "Semen Indonesia", "Cement"},
	{"INTP", "Indocement", "Cement"},
	{"SMCB", "Solusi Bangun Indonesia", "Cement"},
	{"BRPT", "Barito Pacific", "Chemical"},
	{"TPIA", "Chandra Asri", "Chemical"},
	{"INKP", "Indah Kiat Pulp", "Paper"},
	{"TKIM", "Pabrik Kertas Tjiwi", "Paper"},
}

func generateEmitenData() []Emiten {
	var emitens []Emiten

	for _, e := range emitenList {
		price := 100 + rand.Float64()*49900
		prevClose := price * (0.95 + rand.Float64()*0.1)
		change := (price - prevClose) / prevClose * 100

		open := prevClose * (0.99 + rand.Float64()*0.02)
		high := math.Max(open, price) * (1 + rand.Float64()*0.02)
		low := math.Min(open, price) * (1 - rand.Float64()*0.02)

		volume := int64(500000 + rand.Intn(99500000))
		avgVolume := int64(float64(volume) * (0.7 + rand.Float64()*0.6))

		rsi := 20 + rand.Float64()*60
		macd := -2 + rand.Float64()*4
		gap := -2 + rand.Float64()*4
		volatility := 1 + rand.Float64()*4

		morningMom := rand.Float64() * 100
		afternoonDip := rand.Float64() * 100

		scoreBSJP := calculateBSJP(rsi, change, volatility, gap, afternoonDip, volume, avgVolume)
		scoreBPJS := calculateBPJS(rsi, change, volatility, morningMom, volume, avgVolume)

		emitens = append(emitens, Emiten{
			Symbol:        e.symbol,
			Name:          e.name,
			Sector:        e.sector,
			Price:         price,
			Open:          open,
			High:          high,
			Low:           low,
			PrevClose:     prevClose,
			Change:        change,
			Volume:        volume,
			AvgVolume:     avgVolume,
			RSI:           rsi,
			MACD:          macd,
			GapPercent:    gap,
			Volatility:    volatility,
			MorningMoment: morningMom,
			AfternoonDip:  afternoonDip,
			ScoreBSJP:     scoreBSJP,
			ScoreBPJS:     scoreBPJS,
		})
	}

	return emitens
}

func calculateBSJP(rsi, change, volatility, gap, afternoonDip float64, vol, avgVol int64) float64 {
	score := 0.0

	if rsi < 35 {
		score += 20
	} else if rsi < 45 {
		score += 12
	}

	if change > -3 && change < -0.5 {
		score += 18
	} else if change > -5 && change < 0 {
		score += 10
	}

	if volatility > 2 && volatility < 4 {
		score += 15
	}

	if gap > 0.5 {
		score += 15
	} else if gap > 0 {
		score += 8
	}

	if afternoonDip > 60 {
		score += 17
	} else if afternoonDip > 40 {
		score += 10
	}

	if vol > avgVol {
		score += 15
	}

	return math.Min(100, score)
}

func calculateBPJS(rsi, change, volatility, morningMom float64, vol, avgVol int64) float64 {
	score := 0.0

	if rsi > 55 && rsi < 70 {
		score += 18
	} else if rsi > 45 {
		score += 10
	}

	if change > 0.5 && change < 3 {
		score += 20
	} else if change > 0 {
		score += 12
	}

	if volatility > 1.5 && volatility < 3 {
		score += 15
	}

	if morningMom > 60 {
		score += 20
	} else if morningMom > 40 {
		score += 12
	}

	if vol > avgVol*12/10 {
		score += 15
	}

	return math.Min(100, score)
}

func scanBSJP(emitens []Emiten) []ScanResult {
	var results []ScanResult

	for _, e := range emitens {
		if e.ScoreBSJP >= 45 {
			target := e.Price * (1 + e.Volatility*0.3/100)
			stopLoss := e.Low * 0.99

			signal := "WATCH"
			if e.ScoreBSJP >= 75 {
				signal = "STRONG BUY"
			} else if e.ScoreBSJP >= 60 {
				signal = "BUY"
			}

			reason := fmt.Sprintf("RSI=%.0f, Gap=%.1f%%, Vol=%s", e.RSI, e.GapPercent, formatVol(e.Volume))

			results = append(results, ScanResult{
				Symbol:   e.Symbol,
				Name:     e.Name,
				Sector:   e.Sector,
				Price:    e.Price,
				Change:   e.Change,
				Target:   target,
				StopLoss: stopLoss,
				Score:    e.ScoreBSJP,
				Signal:   signal,
				Reason:   reason,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

func scanBPJS(emitens []Emiten) []ScanResult {
	var results []ScanResult

	for _, e := range emitens {
		if e.ScoreBPJS >= 45 {
			target := e.Price * (1 + e.Volatility*0.4/100)
			stopLoss := e.Price * 0.985

			signal := "WATCH"
			if e.ScoreBPJS >= 75 {
				signal = "STRONG BUY"
			} else if e.ScoreBPJS >= 60 {
				signal = "BUY"
			}

			reason := fmt.Sprintf("RSI=%.0f, Mom=%.0f%%, Vol=%s", e.RSI, e.MorningMoment, formatVol(e.Volume))

			results = append(results, ScanResult{
				Symbol:   e.Symbol,
				Name:     e.Name,
				Sector:   e.Sector,
				Price:    e.Price,
				Change:   e.Change,
				Target:   target,
				StopLoss: stopLoss,
				Score:    e.ScoreBPJS,
				Signal:   signal,
				Reason:   reason,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

func formatPrice(p float64) string {
	if p >= 1000 {
		return fmt.Sprintf("Rp%.0f", p)
	}
	return fmt.Sprintf("Rp%.0f", p)
}

func formatVol(v int64) string {
	if v >= 1000000000 {
		return fmt.Sprintf("%.1fB", float64(v)/1000000000)
	} else if v >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(v)/1000000)
	} else if v >= 1000 {
		return fmt.Sprintf("%.1fK", float64(v)/1000)
	}
	return fmt.Sprintf("%d", v)
}

func getStars(score float64) string {
	n := int(score / 20)
	if n > 5 {
		n = 5
	}
	if n < 1 {
		n = 1
	}
	return strings.Repeat("*", n) + strings.Repeat(" ", 5-n)
}

func printHeader() {
	fmt.Print("\033[2J\033[H")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("                          EMITEN SCANNER BSJP & BPJS")
	fmt.Println("                        Indonesia Stock Exchange (IDX)")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf(" Waktu: %s\n", time.Now().Format("02 Jan 2006 15:04:05 WIB"))
	fmt.Println(strings.Repeat("-", 100))
}

func printBSJP(results []ScanResult) {
	fmt.Println()
	fmt.Println("\033[1;33m                    BSJP - BELI SORE JUAL PAGI\033[0m")
	fmt.Println(" Strategi: Beli 14:30-15:00 WIB, Jual 09:00-09:30 WIB besok")
	fmt.Println(strings.Repeat("-", 100))

	if len(results) == 0 {
		fmt.Println(" Tidak ada emiten yang memenuhi kriteria BSJP saat ini.")
		return
	}

	fmt.Printf(" %-7s %-22s %-12s %-10s %-7s %-10s %-10s %-6s %-10s\n",
		"KODE", "NAMA", "SEKTOR", "HARGA", "CHG%", "TARGET", "SL", "SCORE", "SIGNAL")
	fmt.Println(strings.Repeat("-", 100))

	count := 0
	for _, r := range results {
		if count >= 15 {
			break
		}

		chgClr := "\033[32m"
		if r.Change < 0 {
			chgClr = "\033[31m"
		}

		name := r.Name
		if len(name) > 20 {
			name = name[:20]
		}
		sector := r.Sector
		if len(sector) > 10 {
			sector = sector[:10]
		}

		fmt.Printf(" %-7s %-22s %-12s %-10s %s%-6.1f%%\033[0m %-10s %-10s %-6.0f %-10s\n",
			r.Symbol, name, sector, formatPrice(r.Price),
			chgClr, r.Change, formatPrice(r.Target), formatPrice(r.StopLoss),
			r.Score, r.Signal)
		count++
	}

	fmt.Printf("\n Total emiten BSJP: %d\n", len(results))
}

func printBPJS(results []ScanResult) {
	fmt.Println()
	fmt.Println("\033[1;36m                    BPJS - BELI PAGI JUAL SORE\033[0m")
	fmt.Println(" Strategi: Beli 09:00-09:30 WIB, Jual 14:30-15:00 WIB hari yang sama")
	fmt.Println(strings.Repeat("-", 100))

	if len(results) == 0 {
		fmt.Println(" Tidak ada emiten yang memenuhi kriteria BPJS saat ini.")
		return
	}

	fmt.Printf(" %-7s %-22s %-12s %-10s %-7s %-10s %-10s %-6s %-10s\n",
		"KODE", "NAMA", "SEKTOR", "HARGA", "CHG%", "TARGET", "SL", "SCORE", "SIGNAL")
	fmt.Println(strings.Repeat("-", 100))

	count := 0
	for _, r := range results {
		if count >= 15 {
			break
		}

		chgClr := "\033[32m"
		if r.Change < 0 {
			chgClr = "\033[31m"
		}

		name := r.Name
		if len(name) > 20 {
			name = name[:20]
		}
		sector := r.Sector
		if len(sector) > 10 {
			sector = sector[:10]
		}

		fmt.Printf(" %-7s %-22s %-12s %-10s %s%-6.1f%%\033[0m %-10s %-10s %-6.0f %-10s\n",
			r.Symbol, name, sector, formatPrice(r.Price),
			chgClr, r.Change, formatPrice(r.Target), formatPrice(r.StopLoss),
			r.Score, r.Signal)
		count++
	}

	fmt.Printf("\n Total emiten BPJS: %d\n", len(results))
}

func printAllEmiten(emitens []Emiten) {
	printHeader()
	fmt.Println()
	fmt.Println("                           DAFTAR SEMUA EMITEN")
	fmt.Println(strings.Repeat("-", 100))

	fmt.Printf(" %-7s %-20s %-10s %-10s %-7s %-6s %-8s %-8s %-8s\n",
		"KODE", "NAMA", "SEKTOR", "HARGA", "CHG%", "RSI", "VOL", "BSJP", "BPJS")
	fmt.Println(strings.Repeat("-", 100))

	for _, e := range emitens {
		chgClr := "\033[32m"
		if e.Change < 0 {
			chgClr = "\033[31m"
		}

		name := e.Name
		if len(name) > 18 {
			name = name[:18]
		}
		sector := e.Sector
		if len(sector) > 8 {
			sector = sector[:8]
		}

		fmt.Printf(" %-7s %-20s %-10s %-10s %s%-6.1f%%\033[0m %-6.0f %-8s %-8.0f %-8.0f\n",
			e.Symbol, name, sector, formatPrice(e.Price),
			chgClr, e.Change, e.RSI, formatVol(e.Volume),
			e.ScoreBSJP, e.ScoreBPJS)
	}

	fmt.Printf("\n Total emiten: %d\n", len(emitens))
}

func printBySector(emitens []Emiten) {
	printHeader()
	fmt.Println()
	fmt.Println("                         EMITEN PER SEKTOR")
	fmt.Println(strings.Repeat("-", 100))

	sectors := make(map[string][]Emiten)
	for _, e := range emitens {
		sectors[e.Sector] = append(sectors[e.Sector], e)
	}

	var sectorNames []string
	for s := range sectors {
		sectorNames = append(sectorNames, s)
	}
	sort.Strings(sectorNames)

	for _, sector := range sectorNames {
		list := sectors[sector]
		fmt.Printf("\n \033[1;34m%s (%d emiten)\033[0m\n", sector, len(list))
		fmt.Println(strings.Repeat("-", 60))

		for _, e := range list {
			chgClr := "\033[32m"
			if e.Change < 0 {
				chgClr = "\033[31m"
			}

			fmt.Printf(" %-7s %-25s %-10s %s%-6.1f%%\033[0m BSJP=%-3.0f BPJS=%-3.0f\n",
				e.Symbol, e.Name, formatPrice(e.Price),
				chgClr, e.Change, e.ScoreBSJP, e.ScoreBPJS)
		}
	}
}

func printGuide() {
	printHeader()
	fmt.Println()
	fmt.Println("                    PANDUAN STRATEGI BSJP & BPJS")
	fmt.Println(strings.Repeat("-", 100))
	fmt.Println()
	fmt.Println(" BSJP - BELI SORE JUAL PAGI")
	fmt.Println(" ==========================")
	fmt.Println(" Waktu Beli  : 14:30 - 15:00 WIB (menjelang closing)")
	fmt.Println(" Waktu Jual  : 09:00 - 09:30 WIB (keesokan hari)")
	fmt.Println(" Konsep      : Memanfaatkan gap overnight dan momentum pembukaan")
	fmt.Println()
	fmt.Println(" Kriteria:")
	fmt.Println(" - RSI < 45 (kondisi oversold)")
	fmt.Println(" - Harga turun -0.5% s/d -3% (koreksi sehat)")
	fmt.Println(" - Volatilitas 2-4% (potensi gap)")
	fmt.Println(" - Volume di atas rata-rata")
	fmt.Println(" - Ada pattern afternoon dip")
	fmt.Println()
	fmt.Println(" BPJS - BELI PAGI JUAL SORE")
	fmt.Println(" ==========================")
	fmt.Println(" Waktu Beli  : 09:00 - 09:30 WIB (setelah opening)")
	fmt.Println(" Waktu Jual  : 14:30 - 15:00 WIB (hari yang sama)")
	fmt.Println(" Konsep      : Intraday momentum trading")
	fmt.Println()
	fmt.Println(" Kriteria:")
	fmt.Println(" - RSI 55-70 (momentum naik)")
	fmt.Println(" - Harga naik +0.5% s/d +3% saat opening")
	fmt.Println(" - Morning momentum kuat")
	fmt.Println(" - Volume tinggi di awal sesi")
	fmt.Println()
	fmt.Println(" MANAJEMEN RISIKO")
	fmt.Println(" ================")
	fmt.Println(" - Maksimal 3-5 saham per hari")
	fmt.Println(" - Stop loss -1.5% untuk BPJS, -2% untuk BSJP")
	fmt.Println(" - Take profit +1% untuk BPJS, +1.5% untuk BSJP")
	fmt.Println(" - Jangan all-in, diversifikasi")
	fmt.Println()
}

func printStatistics(emitens []Emiten, bsjpResults, bpjsResults []ScanResult) {
	printHeader()
	fmt.Println()
	fmt.Println("                           STATISTIK SCAN")
	fmt.Println(strings.Repeat("-", 100))
	fmt.Println()
	fmt.Printf(" Total Emiten Terscan    : %d\n", len(emitens))
	fmt.Printf(" Emiten Lolos BSJP       : %d\n", len(bsjpResults))
	fmt.Printf(" Emiten Lolos BPJS       : %d\n", len(bpjsResults))
	fmt.Println()

	strongBuyBSJP := 0
	buyBSJP := 0
	for _, r := range bsjpResults {
		if r.Signal == "STRONG BUY" {
			strongBuyBSJP++
		} else if r.Signal == "BUY" {
			buyBSJP++
		}
	}

	strongBuyBPJS := 0
	buyBPJS := 0
	for _, r := range bpjsResults {
		if r.Signal == "STRONG BUY" {
			strongBuyBPJS++
		} else if r.Signal == "BUY" {
			buyBPJS++
		}
	}

	fmt.Println(" BSJP Signals:")
	fmt.Printf("   STRONG BUY : %d\n", strongBuyBSJP)
	fmt.Printf("   BUY        : %d\n", buyBSJP)
	fmt.Printf("   WATCH      : %d\n", len(bsjpResults)-strongBuyBSJP-buyBSJP)
	fmt.Println()
	fmt.Println(" BPJS Signals:")
	fmt.Printf("   STRONG BUY : %d\n", strongBuyBPJS)
	fmt.Printf("   BUY        : %d\n", buyBPJS)
	fmt.Printf("   WATCH      : %d\n", len(bpjsResults)-strongBuyBPJS-buyBPJS)
	fmt.Println()
}

func printMenu() {
	fmt.Println()
	fmt.Println(strings.Repeat("-", 100))
	fmt.Println(" MENU:")
	fmt.Println(" [1] Scan Ulang")
	fmt.Println(" [2] Lihat Semua Emiten")
	fmt.Println(" [3] Lihat Per Sektor")
	fmt.Println(" [4] Statistik")
	fmt.Println(" [5] Panduan Strategi")
	fmt.Println(" [6] Keluar")
	fmt.Println()
	fmt.Print(" Pilihan: ")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		emitens := generateEmitenData()

		printHeader()

		bsjpResults := scanBSJP(emitens)
		bpjsResults := scanBPJS(emitens)

		printBSJP(bsjpResults)
		printBPJS(bpjsResults)

		printMenu()

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			continue
		case "2":
			printAllEmiten(emitens)
			fmt.Println("\n Tekan Enter...")
			fmt.Scanln()
		case "3":
			printBySector(emitens)
			fmt.Println("\n Tekan Enter...")
			fmt.Scanln()
		case "4":
			printStatistics(emitens, bsjpResults, bpjsResults)
			fmt.Println("\n Tekan Enter...")
			fmt.Scanln()
		case "5":
			printGuide()
			fmt.Println("\n Tekan Enter...")
			fmt.Scanln()
		case "6", "q", "Q":
			fmt.Println()
			fmt.Println(" Terima kasih!")
			fmt.Println(" Selamat berinvestasi dengan bijak.")
			fmt.Println()
			return
		}
	}
}
