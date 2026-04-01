# Stability Team Technical Test - Submission

## 1) Issues Found

1. Belum ada input validation saat membuat task baru.
2. Data task disimpan di variabel in-memory, sehingga tidak persisten.
3. Ada risiko race condition saat banyak request mengakses data secara bersamaan.
4. Format response API belum konsisten dan kurang jelas.
5. Belum ada logic auto-increment ID saat menambahkan task baru.

## 2) How I Fixed Them

1. Menambahkan validation untuk request create/update task menggunakan validator pada field yang dibutuhkan.
2. Memindahkan penyimpanan task dari variabel hardcoded ke file `data/tasks.json` agar data persisten.
3. Menambahkan `sync.RWMutex` untuk mengamankan operasi read/write ke data task.
4. Menstandarkan format response sukses dan error menggunakan type response agar konsisten di semua endpoint.
5. Menambahkan logic auto-generate ID saat create task dengan mengambil ID terbesar yang ada lalu menambahkannya 1.

## 3) Improvements Added

1. Menambahkan endpoint edit task (update title dan done) dengan pengecekan `no changes`.
2. Menambahkan auto-generate ID saat create task.
3. Menambahkan rate limiter pada aplikasi untuk membatasi request berlebihan.
