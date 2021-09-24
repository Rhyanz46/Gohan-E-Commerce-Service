# Fitur Aplikasi 
- Tambah produk
- Lihat produk
- Hapus Produk
- List Produk dengan pagination
- Upload foto suatu produk maximal 10 foto
- Hapus foto produk


# Endpoints
- User
  - /user/login _[POST]_
  - /user/register _[POST]_
  - /user/detail _[GET, PUT]_
- Product
  - /product _[POST, GET]_
  - /product/{product_id:[0-9]+} _[PUT, GET, DELETE]_
  - /product/{product_id:[0-9]+}/photos _[POST, GET]_
  - /product/{product_id:[0-9]+}/photos/{photo_id:[0-9]+} _[GET, DELETE]_

