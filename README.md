# IDX Stock Scanner

Scanner saham Indonesia (IDX) dengan dua strategi trading berbeda.

## Program yang Tersedia

### 1. Stock Scanner - Beli Sore Jual Pagi (C++)

Scanner untuk strategi swing trading overnight dengan memanfaatkan gap pembukaan.

**File:** `stock_scanner.cpp` / `stock_scanner.exe`

**Strategi:**
- Beli saham menjelang closing (14:30-15:00 WIB)
- Jual saat pembukaan (09:00-09:30 WIB)
- Memanfaatkan gap overnight karena sentimen global

**Kriteria Beli Sore:**
- RSI < 40 (oversold)
- Harga di bawah SMA20
- Volume di atas rata-rata
- Penurunan -0.5% sampai -3%

**Kriteria Jual Pagi:**
- RSI > 60 (overbought)
- Gap up terjadi
- Take profit 0.5% - 2%
- Stop loss -2%

---

### 2. Net Foreign Buy Scanner (Golang)

Scanner untuk tracking akumulasi dan distribusi oleh investor asing.

**File:** `net_foreign_scanner.go` / `net_foreign_scanner.exe`

**Konsep:**
Net Foreign Buy adalah selisih pembelian dan penjualan oleh investor asing. Positif berarti asing lebih banyak membeli (akumulasi), negatif berarti asing lebih banyak menjual (distribusi).

**Kolom Data:**
| Kolom | Deskripsi |
|-------|-----------|
| NET FB | Net Foreign Buy dalam lot |
| VALUE | Nilai rupiah net foreign |
| F% | Persentase transaksi asing dari total volume |
| ACC | Hari akumulasi berturut-turut |
| SIGNAL | Sinyal berdasarkan analisa |

**Signal:**
- STRONG BUY - Score >= 70
- BUY - Score >= 55
- ACCUMULATE - Score >= 40
- DISTRIBUTE - Asing mulai jual
- STRONG SELL - Asing jual besar-besaran

---

## Kompilasi

### C++ Scanner
```bash
g++ -o stock_scanner.exe stock_scanner.cpp -std=c++11
```

### Golang Scanner
```bash
go build -o net_foreign_scanner.exe net_foreign_scanner.go
```

---

## Menjalankan Program

```bash
# C++ Scanner
.\stock_scanner.exe

# Golang Scanner
.\net_foreign_scanner.exe
```

---

## Menu Program

Kedua program memiliki menu interaktif:
1. Scan Ulang - Refresh data saham
2. Lihat Semua Saham - Tampilkan seluruh data
3. Panduan - Penjelasan strategi
4. Keluar - Tutup program

---

## Daftar Saham

Scanner mencakup 30+ saham blue chip Indonesia:

| Sektor | Saham |
|--------|-------|
| Perbankan | BBCA, BBRI, BMRI, ARTO, BRIS |
| Telekomunikasi | TLKM, EXCL |
| Konsumer | UNVR, ICBP, INDF, GGRM, HMSP |
| Teknologi | GOTO, BUKA, EMTK |
| Tambang | ANTM, INCO, PTBA, ADRO, ITMG, MDKA, AMMN |
| Infrastruktur | JSMR, PGAS |
| Lainnya | ASII, CPIN, JPFA, ACES, ERAA, MAPA, SIDO, KLBF |

---

## Catatan Penting

1. Data dalam program ini adalah simulasi untuk keperluan edukasi
2. Bukan merupakan rekomendasi investasi
3. Selalu lakukan riset mandiri sebelum trading
4. Gunakan manajemen risiko yang baik
5. Konsultasikan dengan penasihat keuangan profesional

---

## Requirements

- **C++ Scanner:** MinGW/GCC dengan C++11 support
- **Golang Scanner:** Go 1.16 atau lebih baru
- **Terminal:** Support ANSI color codes (Windows Terminal, PowerShell, atau CMD dengan ANSI enabled)

---

## Author

I Gede Satria Adi Pratama

---

## License

MIT License - Bebas digunakan untuk edukasi dan pengembangan pribadi.
