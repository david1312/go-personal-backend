package constants

// BNI
var (
	ATMInstructionBNI             []string = []string{"Masukkan Kartu Anda.", "Pilih Bahasa.", "Masukkan PIN ATM Anda.", "Pilih `Menu Lainnya`.", "Pilih `Transfer`.", "Pilih Jenis rekening yang akan Anda gunakan (Contoh; `Dari Rekening Tabungan`).", "Pilih `Virtual Account Billing`.", "Masukkan nomor Virtual Account Anda (contoh: 8241002201150001).", "Tagihan yang harus dibayarkan akan muncul pada layar konfirmasi.", "Konfirmasi, apabila telah sesuai, lanjutkan transaksi.", "Transaksi Anda telah selesai."}
	InternetBankingInstructionBNI []string = []string{"Ketik alamat https://ibank.bni.co.id kemudian klik `Enter`.", "Masukkan User ID dan Password.", "Pilih menu `Transfer`.", "Pilih `Virtual Account Billing`.", "Kemudian masukan nomor Virtual Account Anda (contoh: 8241002201150001) yang hendak dibayarkan. Lalu pilih rekening debet yang akan digunakan. Kemudian tekan `lanjut`.", "Kemudin tagihan yang harus dibayarkan akan muncul pada layar konfirmasi.", "Masukkan Kode Otentikasi Token.", "Pembayaran Anda telah berhasil."}
	MobileBankingInstructionBNI   []string = []string{"Akses BNI Mobile Banking dari handphone kemudian masukkan user ID dan password.", "Pilih menu `Transfer`.", "Pilih menu `Virtual Account Billing` kemudian pilih rekening debet.", "Masukkan nomor Virtual Account Anda (contoh: 8241002201150001) pada menu `input baru`.", "Tagihan yang harus dibayarkan akan muncul pada layar konfirmasi", "Konfirmasi transaksi dan masukkan Password Transaksi.", "Pembayaran Anda Telah Berhasil."}
	CODIntruction                 []string = []string{"Datang ke outlet kami di alamat sesuai dengan ke outlet yang anda pilih", "Ada beberapa cara pembayaran sistem COD yang bisa dilakukan.", "Pertama anda bisa bayar menggunakan uang cash / tunai.", "Kedua Anda juga bisa membayar menggunakan mesin EDC / Scan QRIS yang ada dioutlet kami.", "Terakhir anda dapat membayar dengan scan QRIS yang ada di halaman ini", "Jangan lupa simpan bukti transfer dan tunjukan bukti pemesanan kepada kasir untuk memverifikasi pembayaran anda."}

	ATMInstructionPermata             []string = []string{"Pada menu utama, pilih Transaksi Lainnya", "Pilih menu Pembayaran", "Pilih Pembayaran Lainnya", "Pilih Virtual Account", "Masukkan nomor Virtual Account (misal 7508123456789012)", "Jumlah yang harus dibayar dan nomor rekening akan muncul pada halaman konfirmasi pembayaran. Jika informasi sudah benar, pilih Benar"}
	InternetBankingInstructionPermata []string = []string{"Buka website PermataNet https://www.permatanet.com/", "Masukan User ID & Password", "Pilih Pembayaran Tagihan", " Pilih Virtual Account", "Masukkan 16 digit kode bayar (7508123456789012)", "Masukkan nominal pembayaran", "Muncul konfirmasi pembayaran", "Masukan Mobile PIN", "Transaksi selesai"}
	MobileBankingInstructionPermata   []string = []string{"Buka aplikasi PermataMobile X", "Masukan Password", "Pilih Bayar Tagihan", "Pilih Virtual Account", "Masukkan Nomor Virtual Account (7508123456789012)", "Pilih rekening", "Masukkan nominal pembayaran", "Muncul konfirmasi pembayaran", "Masukan Mobile PIN", "Transaksi selesai"}

	ATMInstructionBRI             []string = []string{"Masukkan Kartu Debit BRI dan PIN Anda", "Pilih menu Transaksi Lain Pembayaran Lainnya > BRIVA", "Masukkan 5 angka kode perusahaan untuk Semesta Ban (80777) dan Nomor HP yang terdaftar di akun Semesta Ban Anda (Contoh 80777085314782388)", "Di halaman konfirmasi, pastikan detail pembayaran sudah sesuai seperti Nomor BRIVA, Nama Pelanggan dan Jumlah Pembayaran", "Ikuti instruksi untuk menyelesaikan transaksi", "Simpan struk transaksi sebagai bukti pembayaran"}
	InternetBankingInstructionBRI []string = []string{"Login aplikasi BRI Mobile", "Pilih menu Mobile Banking BRI > Pembayaran > BRIVA 3. Masukkan nomor virtual account", "Masukan Jumlah Pembayaran", "Masukkan PIN", "Simpan notifikasi SMS sebagai bukti pembayaran"}
	MobileBankingInstructionBRI   []string = []string{"Login pada alamat Internet Banking BRI (https://bri.co.id/internet-banking)", "Pilih menu Pembayaran Tagihan Pembayaran >BRIVA", "Masukan nomor virtual account", "Di halaman konfirmasi, pastikan detail pembayaran sudah sesuai seperti Nomor BRIVA, Nama Pelanggan dan Jumlah Pembayaran", "Masukkan password dan mToken", "Cetak/simpan struk pembayaran BRIVA sebagai bukti pembayaran"}
)
