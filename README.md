# Fitur Aplikasi 
- Tambah produk
- Lihat produk
- Hapus Produk
- List Produk dengan pagination
- Upload foto suatu produk maximal 10 foto
- Hapus foto produk


# Endpoint-endpoint
- User
  - /user/login _[POST]_
  - /user/register _[POST]_
  - /user/detail _[GET, PUT]_
- Product
  - /product _[POST, GET]_
  - /product/{product_id:[0-9]+} _[PUT, GET, DELETE]_
  - /product/{product_id:[0-9]+}/photos _[POST, GET]_
  - /product/{product_id:[0-9]+}/photos/{photo_id:[0-9]+} _[GET, DELETE]_

# Installasi
## Konfigurasi
untuk menggunakan service ini, anda perlu membuat konfigurasi .yaml kurang lebih seperti ini
```yaml
port: ":8888"
static_folder: "static/"
secret_key: "222222222223333332323232323@20200218Gara2BCA"
jwt_expired_time: 2
primary_db:
  host: "172.19.0.2"
  username: "root"
  password: "a"
  db_name: "CodeInterview"
  port: "3306"
```
Keterangan :
- anda hanya perlu membuat database kemudian program otomatis akan membuat schema table.
- `jwt_expired_time` hitungan dalam jam dan nilai default adalah 1 yang berarti 1 jam, jadi jika anda tidak memasukkan value 
kedalam field ini maka nilai default akan terpakai yaitu 1 jam.
- value dari field `static_folder` harus di akhiri dengan `/`
- value dari field `static_folder` harus ada

## Penggunaan
Setelah konfigurasi selesai kita bisa menggunakan service untuk :
- unit testing
  `make test`
- run
  `make run`

# Detail tentang endpoint
- User
  - /user/login _[POST]_
  - /user/register _[POST]_
  - /user/detail _[GET, PUT]_
- Product
  - /product _[POST, GET]_
  - /product/{product_id:[0-9]+} _[PUT, GET, DELETE]_
  - /product/{product_id:[0-9]+}/photos _[POST, GET]_
  - /product/{product_id:[0-9]+}/photos/{photo_id:[0-9]+} _[GET, DELETE]_
