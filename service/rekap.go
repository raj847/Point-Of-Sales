package service

func DoRekapEachMonth() {
	// wib, _ := time.LoadLocation("Asia/Jakarta")
	// s := gocron.NewScheduler(wib)

	// s.Every(1).Month().At("23:00").Do(Rekap)
	// s.StartAsync()
}

func Rekap() {
	// ambil data transaksi dari database berdasarkan bulan sekarang (based on created_at)

	/*
		product_id
		product_name
		total_terjual
		sisa_stock
		total_pendapatan
		bulan
	*/

	// get current month
	// now := time.Now()
	// wib, _ := time.LoadLocation("Asia/Jakarta")
	// startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, wib)
	// endOfMonth := startOfMonth.AddDate(0, 1, -1)

	// sample query using gorm
	// result, err := db.Table("transactions").Select("sum(amount) as total").Where("created_at >= ? AND created_at <= ?", startOfMonth, endOfMonth).Rows()

	// ambil data

	// hitung

	// save

}
