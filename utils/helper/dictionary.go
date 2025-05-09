package helper

var ValidCategories = map[string][]string{
	"Asset": {
		"Kas & Bank",
		"Piutang",
		"Perlengkapan",
		"Aset Tetap",
	},
	"Liability": {
		"Utang Dagang",
		"Pinjaman",
		"Kewajiban Lancar",
		"Pajak",
		"Pembayaran Bagi Hasil",
		"Hutang Usaha",
	},
	"Equity": {
		"Modal",
		"Laba Ditahan",
	},
	"Revenue": {
		"Penjualan",
		"Jasa",
		"Pendapatan Dari Bunga",
	},

	"Expense": {
		"Operasional",
		"Utilitas",
		"Non Operasional",
	},
}
