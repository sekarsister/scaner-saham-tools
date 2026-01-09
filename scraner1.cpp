#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <cmath>
#include <iomanip>
#include <ctime>
#include <sstream>
#include <map>
#include <random>

using namespace std;

struct StockData {
    string symbol;
    string name;
    double openPrice;
    double highPrice;
    double lowPrice;
    double closePrice;
    double prevClose;
    long volume;
    double changePercent;
    double avgVolume;
    double rsi;
    double sma20;
    double ema9;
    double volatility;
    double gapPercent;
    double afternoonScore;
    double morningScore;
};

struct SignalResult {
    string symbol;
    string name;
    string signalType;
    double currentPrice;
    double targetPrice;
    double stopLoss;
    double potentialGain;
    double riskReward;
    int strength;
    string reason;
};

class StockDatabase {
private:
    vector<StockData> stocks;
    mt19937 rng;

    void calculateScores(StockData& data) {
        double score = 0;
        
        if (data.rsi < 30) score += 25;
        else if (data.rsi < 40) score += 15;
        else if (data.rsi < 50) score += 5;
        
        if (data.closePrice < data.sma20 * 0.98) score += 20;
        else if (data.closePrice < data.sma20) score += 10;
        
        if (data.volume > data.avgVolume * 1.5) score += 15;
        else if (data.volume > data.avgVolume) score += 8;
        
        if (data.volatility > 2.0 && data.volatility < 4.0) score += 15;
        if (data.changePercent > -3.0 && data.changePercent < -0.5) score += 15;
        if (data.gapPercent > 0.5) score += 10;
        
        data.afternoonScore = min(100.0, score);
        
        score = 0;
        if (data.rsi > 70) score += 25;
        else if (data.rsi > 60) score += 15;
        else if (data.rsi > 50) score += 5;
        
        if (data.closePrice > data.sma20 * 1.02) score += 20;
        else if (data.closePrice > data.sma20) score += 10;
        
        if (data.volume > data.avgVolume * 1.3 && data.changePercent > 0) score += 15;
        if (data.gapPercent > 1.0) score += 20;
        else if (data.gapPercent > 0.5) score += 10;
        if (data.changePercent > 1.0) score += 10;
        
        data.morningScore = min(100.0, score);
    }

    void generateRandomData(StockData& data) {
        uniform_real_distribution<> priceDist(500, 50000);
        uniform_real_distribution<> changeDist(-5.0, 5.0);
        uniform_real_distribution<> volumeDist(1000000, 100000000);
        uniform_real_distribution<> rsiDist(20, 80);
        uniform_real_distribution<> volDist(1.0, 5.0);

        data.prevClose = priceDist(rng);
        data.changePercent = changeDist(rng);
        data.closePrice = data.prevClose * (1 + data.changePercent / 100);
        
        double range = data.closePrice * 0.03;
        data.openPrice = data.closePrice + (changeDist(rng) / 100 * data.closePrice);
        data.highPrice = max(data.openPrice, data.closePrice) + abs(range * 0.5);
        data.lowPrice = min(data.openPrice, data.closePrice) - abs(range * 0.5);
        
        data.volume = static_cast<long>(volumeDist(rng));
        data.avgVolume = data.volume * (0.8 + (changeDist(rng) + 5) / 10 * 0.4);
        
        data.rsi = rsiDist(rng);
        data.sma20 = data.closePrice * (0.95 + (rsiDist(rng) - 50) / 500);
        data.ema9 = data.closePrice * (0.98 + (rsiDist(rng) - 50) / 1000);
        data.volatility = volDist(rng);
        
        uniform_real_distribution<> gapDist(-2.0, 3.0);
        data.gapPercent = gapDist(rng);
        
        calculateScores(data);
    }

public:
    StockDatabase() {
        rng.seed(time(nullptr));
        
        vector<pair<string, string>> stockList = {
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
            {"EXCL", "XL Axiata"}
        };

        for (const auto& stock : stockList) {
            StockData data;
            data.symbol = stock.first;
            data.name = stock.second;
            generateRandomData(data);
            stocks.push_back(data);
        }
    }

    vector<StockData> getAllStocks() { return stocks; }
    
    void refreshData() {
        for (auto& stock : stocks) {
            generateRandomData(stock);
        }
    }
};

class StockScanner {
private:
    StockDatabase& db;
    vector<SignalResult> buySignals;
    vector<SignalResult> sellSignals;

    SignalResult generateBuySignal(const StockData& stock) {
        SignalResult signal;
        signal.symbol = stock.symbol;
        signal.name = stock.name;
        signal.signalType = "BUY_AFTERNOON";
        signal.currentPrice = stock.closePrice;
        
        double expectedGap = max(0.5, stock.volatility * 0.5);
        signal.targetPrice = stock.closePrice * (1 + expectedGap / 100);
        signal.stopLoss = stock.lowPrice * 0.99;
        
        signal.potentialGain = ((signal.targetPrice - signal.currentPrice) / signal.currentPrice) * 100;
        double risk = ((signal.currentPrice - signal.stopLoss) / signal.currentPrice) * 100;
        signal.riskReward = signal.potentialGain / max(0.1, risk);
        
        signal.strength = static_cast<int>(stock.afternoonScore / 20);
        signal.strength = max(1, min(5, signal.strength));
        
        stringstream ss;
        ss << "RSI=" << fixed << setprecision(1) << stock.rsi;
        if (stock.rsi < 40) ss << " (oversold)";
        ss << ", Vol=" << setprecision(0) << (stock.volume / 1000000.0) << "M";
        if (stock.volume > stock.avgVolume) ss << " (above avg)";
        ss << ", Gap=" << setprecision(1) << stock.gapPercent << "%";
        signal.reason = ss.str();
        
        return signal;
    }

    SignalResult generateSellSignal(const StockData& stock) {
        SignalResult signal;
        signal.symbol = stock.symbol;
        signal.name = stock.name;
        signal.signalType = "SELL_MORNING";
        signal.currentPrice = stock.closePrice;
        
        double expectedGap = stock.gapPercent;
        signal.targetPrice = stock.closePrice * (1 + expectedGap / 100);
        signal.stopLoss = stock.closePrice * 0.98;
        
        signal.potentialGain = expectedGap;
        signal.riskReward = signal.potentialGain / 2.0;
        
        signal.strength = static_cast<int>(stock.morningScore / 20);
        signal.strength = max(1, min(5, signal.strength));
        
        stringstream ss;
        ss << "RSI=" << fixed << setprecision(1) << stock.rsi;
        if (stock.rsi > 60) ss << " (overbought)";
        ss << ", Expected gap=" << setprecision(1) << stock.gapPercent << "%";
        ss << ", R:R=" << setprecision(2) << signal.riskReward;
        signal.reason = ss.str();
        
        return signal;
    }

public:
    StockScanner(StockDatabase& database) : db(database) {}

    void scan() {
        buySignals.clear();
        sellSignals.clear();
        
        vector<StockData> stocks = db.getAllStocks();
        
        for (const auto& stock : stocks) {
            if (stock.afternoonScore >= 50) {
                SignalResult signal = generateBuySignal(stock);
                if (signal.strength >= 2) buySignals.push_back(signal);
            }
            
            if (stock.morningScore >= 50) {
                SignalResult signal = generateSellSignal(stock);
                if (signal.strength >= 2) sellSignals.push_back(signal);
            }
        }
        
        sort(buySignals.begin(), buySignals.end(), 
             [](const SignalResult& a, const SignalResult& b) { return a.strength > b.strength; });
        
        sort(sellSignals.begin(), sellSignals.end(),
             [](const SignalResult& a, const SignalResult& b) { return a.strength > b.strength; });
    }

    vector<SignalResult> getBuySignals() { return buySignals; }
    vector<SignalResult> getSellSignals() { return sellSignals; }
};

string formatPrice(double price) {
    stringstream ss;
    ss << "Rp " << fixed << setprecision(0) << price;
    return ss.str();
}

string formatVolume(long vol) {
    stringstream ss;
    if (vol >= 1000000000) ss << fixed << setprecision(1) << (vol / 1000000000.0) << "B";
    else if (vol >= 1000000) ss << fixed << setprecision(1) << (vol / 1000000.0) << "M";
    else if (vol >= 1000) ss << fixed << setprecision(1) << (vol / 1000.0) << "K";
    else ss << vol;
    return ss.str();
}

string getStars(int count) {
    return string(count, '*') + string(5 - count, ' ');
}

void printLine(int width = 90) {
    cout << string(width, '=') << endl;
}

void printHeader() {
    cout << "\033[2J\033[H";
    printLine();
    cout << "                    STOCK SCANNER - BELI SORE JUAL PAGI" << endl;
    cout << "                      Indonesia Stock Exchange (IDX)" << endl;
    printLine();
    time_t now = time(nullptr);
    cout << " Waktu: " << ctime(&now);
    cout << string(90, '-') << endl << endl;
}

void printBuySignals(vector<SignalResult>& signals) {
    cout << "\033[1;32m";
    cout << "                         SINYAL BELI SORE" << endl;
    cout << "\033[0m";
    cout << " Beli menjelang closing (14:30-15:00 WIB)" << endl;
    cout << string(90, '-') << endl;
    
    if (signals.empty()) {
        cout << " Tidak ada sinyal beli saat ini." << endl;
        return;
    }
    
    cout << left << setw(8) << " KODE" << setw(22) << "NAMA" << setw(12) << "HARGA"
         << setw(12) << "TARGET" << setw(12) << "SL" << setw(7) << "GAIN"
         << setw(6) << "R:R" << setw(6) << "RATE" << endl;
    cout << string(90, '-') << endl;
    
    int count = 0;
    for (auto& sig : signals) {
        if (count >= 8) break;
        cout << " " << left << setw(7) << sig.symbol
             << setw(22) << sig.name.substr(0, 19)
             << setw(12) << formatPrice(sig.currentPrice)
             << setw(12) << formatPrice(sig.targetPrice)
             << setw(12) << formatPrice(sig.stopLoss)
             << setw(7) << fixed << setprecision(1) << sig.potentialGain << "%"
             << setw(6) << fixed << setprecision(2) << sig.riskReward
             << setw(6) << getStars(sig.strength) << endl;
        cout << "         \033[2m" << sig.reason << "\033[0m" << endl;
        count++;
    }
    cout << endl;
}

void printSellSignals(vector<SignalResult>& signals) {
    cout << "\033[1;31m";
    cout << "                         SINYAL JUAL PAGI" << endl;
    cout << "\033[0m";
    cout << " Jual saat opening (09:00-09:30 WIB)" << endl;
    cout << string(90, '-') << endl;
    
    if (signals.empty()) {
        cout << " Tidak ada sinyal jual saat ini." << endl;
        return;
    }
    
    cout << left << setw(8) << " KODE" << setw(22) << "NAMA" << setw(12) << "HARGA"
         << setw(12) << "EXP.OPEN" << setw(9) << "GAP"
         << setw(6) << "R:R" << setw(6) << "RATE" << endl;
    cout << string(90, '-') << endl;
    
    int count = 0;
    for (auto& sig : signals) {
        if (count >= 8) break;
        cout << " " << left << setw(7) << sig.symbol
             << setw(22) << sig.name.substr(0, 19)
             << setw(12) << formatPrice(sig.currentPrice)
             << setw(12) << formatPrice(sig.targetPrice)
             << setw(9) << fixed << setprecision(1) << sig.potentialGain << "%"
             << setw(6) << fixed << setprecision(2) << sig.riskReward
             << setw(6) << getStars(sig.strength) << endl;
        cout << "         \033[2m" << sig.reason << "\033[0m" << endl;
        count++;
    }
    cout << endl;
}

void printAllStocks(vector<StockData>& stocks) {
    printHeader();
    cout << "                           DATA SEMUA SAHAM" << endl;
    cout << string(90, '-') << endl;
    
    cout << left << setw(8) << " KODE" << setw(18) << "NAMA" << setw(10) << "CLOSE"
         << setw(8) << "CHG%" << setw(7) << "VOL" << setw(6) << "RSI"
         << setw(8) << "BUY" << setw(8) << "SELL" << endl;
    cout << string(90, '-') << endl;
    
    for (auto& stock : stocks) {
        string clr = stock.changePercent >= 0 ? "\033[32m" : "\033[31m";
        cout << " " << left << setw(7) << stock.symbol
             << setw(18) << stock.name.substr(0, 15)
             << setw(10) << formatPrice(stock.closePrice)
             << clr << setw(8) << fixed << setprecision(1) << stock.changePercent << "%" << "\033[0m"
             << setw(7) << formatVolume(stock.volume)
             << setw(6) << fixed << setprecision(0) << stock.rsi
             << setw(8) << fixed << setprecision(0) << stock.afternoonScore
             << setw(8) << fixed << setprecision(0) << stock.morningScore << endl;
    }
}

void printStrategy() {
    printHeader();
    cout << "                         PANDUAN STRATEGI" << endl;
    cout << string(90, '-') << endl;
    cout << endl;
    cout << " KONSEP:" << endl;
    cout << " Strategi ini memanfaatkan gap overnight dimana harga saham" << endl;
    cout << " bergerak saat market tutup karena berita atau sentimen global." << endl;
    cout << endl;
    cout << " WAKTU:" << endl;
    cout << " BELI  -> 14:30 - 15:00 WIB (sebelum closing)" << endl;
    cout << " JUAL  -> 09:00 - 09:30 WIB (setelah opening)" << endl;
    cout << endl;
    cout << " KRITERIA BELI:" << endl;
    cout << " - RSI < 40 (oversold)" << endl;
    cout << " - Harga di bawah SMA20" << endl;
    cout << " - Volume di atas rata-rata" << endl;
    cout << " - Penurunan -0.5% sampai -3%" << endl;
    cout << endl;
    cout << " KRITERIA JUAL:" << endl;
    cout << " - RSI > 60 (overbought)" << endl;
    cout << " - Gap up terjadi" << endl;
    cout << " - Take profit 0.5% - 2%" << endl;
    cout << " - Stop loss -2%" << endl;
    cout << endl;
    cout << " RISIKO:" << endl;
    cout << " - Max 3-5 saham per hari" << endl;
    cout << " - Alokasi 10-20% per saham" << endl;
    cout << " - Selalu pasang stop loss" << endl;
    cout << endl;
}

void printMenu() {
    cout << endl;
    cout << " MENU:" << endl;
    cout << " [1] Scan Ulang" << endl;
    cout << " [2] Lihat Semua Saham" << endl;
    cout << " [3] Panduan Strategi" << endl;
    cout << " [4] Keluar" << endl;
    cout << endl;
    cout << " Pilihan: ";
}

int main() {
    StockDatabase db;
    StockScanner scanner(db);
    
    char choice;
    
    while (true) {
        printHeader();
        scanner.scan();
        
        vector<SignalResult> buys = scanner.getBuySignals();
        vector<SignalResult> sells = scanner.getSellSignals();
        
        printBuySignals(buys);
        printSellSignals(sells);
        printMenu();
        
        cin >> choice;
        
        if (choice == '1') {
            db.refreshData();
        }
        else if (choice == '2') {
            vector<StockData> allStocks = db.getAllStocks();
            printAllStocks(allStocks);
            cout << endl << " Tekan Enter...";
            cin.ignore();
            cin.get();
        }
        else if (choice == '3') {
            printStrategy();
            cout << " Tekan Enter...";
            cin.ignore();
            cin.get();
        }
        else if (choice == '4' || choice == 'q' || choice == 'Q') {
            cout << endl;
            cout << " Terima kasih!" << endl;
            cout << " Selamat berinvestasi." << endl << endl;
            break;
        }
    }
    
    return 0;
}
