# Point of Sales

## Proses secara garis besar (Belum final)

1. Manajemen Inventori
2. Laporan Penjualan
3. Manajemen Staf
4. Proses Pembayaran

## Functional Requirement (Belum final)

-   Super Admin (admin toko utama)

    -   Login
    -   Manage Admin cabang (cabang 1,2,3, ...)
        -   Create & Delete account for admin cabang
        -   Suspend & Unsuspend account for admin cabang
        -   Update profile & Reset password for admin cabang
    -   Melihat laporan penjualan (sales reports) dari setiap cabang (tabel, grafik, pdf / excell)
        -   Penjualan per hari
        -   Penjualan per minggu
        -   Penjualan per bulan

-   Admin Cabang

    -   login
    -   Manage cashier berdasarkan cabang masing-masing
        -   Create & Delete account for cashier cabang
        -   Suspend & Unsuspend account for cashier cabang
        -   Update profile & Reset password for cashier cabang
    -   Melihat laporan penjualan (sales reports) dari cabang yang bersangkutan (Tabel & grafik)
        -   Penjualan per hari
        -   Penjualan per minggu
        -   Penjualan per bulan
    -   Manage inventory
        -   CRUD barang

-   Cashier per cabang
    -   Login
    -   CRUD order, transaksi

# Schema Database (Belum final)

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE branch (
   branch_id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
   branch_name VARCHAR(255) NOT NULL,
   branch_address TEXT NOT NULL,
   branch_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   branch_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE "user" (
   user_id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
   user_name VARCHAR(255) NOT NULL,
   user_email VARCHAR(255) NOT NULL UNIQUE,
   user_password VARCHAR(255) NOT NULL,
   user_role ENUM ('super_admin', 'branch_admin', 'cashier') NOT NULL,
   user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   user_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE branch_user (
   branch_user_id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
   branch_id uuid NOT NULL,
   user_id uuid NOT NULL,
   branch_user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   branch_user_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   FOREIGN KEY (branch_id) REFERENCES branch (branch_id),
   FOREIGN KEY (user_id) REFERENCES "user" (user_id)
);

CREATE TABLE inventory (
   item_id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
   item_name VARCHAR(255) NOT NULL,
   item_category VARCHAR(255) NOT NULL,
   item_description TEXT,
   item_price DECIMAL(10,2) NOT NULL,
   item_unit VARCHAR(10) NOT NULL,
   item_reorder_level INT NOT NULL,
   item_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   item_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE branch_inventory (
   branch_inventory_id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
   branch_id uuid NOT NULL,
   item_id uuid NOT NULL,
   item_quantity INT NOT NULL,
   branch_inventory_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   branch_inventory_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   FOREIGN KEY (branch_id) REFERENCES branch (branch_id),
   FOREIGN KEY (item_id) REFERENCES inventory (item_id)
);

CREATE TABLE transaction (
   transaction_id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
   transaction_code VARCHAR(255) NOT NULL,
   branch_id uuid NOT NULL,
   user_id uuid NOT NULL,
   transaction_total DECIMAL(10,2) NOT NULL,
   transaction_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (branch_id) REFERENCES branch (branch_id),
   FOREIGN KEY (user_id) REFERENCES "user" (user_id)
);

CREATE TABLE orders (
   order_id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
   order_code VARCHAR(255) NOT NULL,
   branch_id uuid NOT NULL,
   user_id uuid NOT NULL,
   order_total DECIMAL(10,2) NOT NULL,
   order_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (branch_id) REFERENCES branch (branch_id),
   FOREIGN KEY (user_id) REFERENCES user (user_id)
);

CREATE TABLE order_item (
   order_item_id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
   order_id uuid NOT NULL,
   item_id uuid NOT NULL,
   item_quantity INT NOT NULL,
   item_price DECIMAL(10,2) NOT NULL,
   order_item_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (order_id) REFERENCES orders (order_id),
   FOREIGN KEY (item_id) REFERENCES inventory (item_id)
);
```

questions:

-   inventory nya cukup 1 tabel aja arya? atau perlu dibuat untuk setiap cabang?
    nanti dari gudang di toko utama, baru akan dikirim ke cabang-cabang?
    gimana ya bagusnya arya xixii
