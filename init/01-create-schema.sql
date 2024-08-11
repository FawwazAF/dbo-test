CREATE TABLE IF NOT EXISTS dbo_mst_product (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price REAL NOT NULL DEFAULT 0,
    stock INT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS dbo_trx_customer(
    id bigserial not null primary key,
    username varchar(20) not null,
    password varchar not null,
    name varchar(100) not null default '',
    email varchar(50) not null default '',
    phone_number varchar(20) not null default '',
    date_of_birth date not null default CURRENT_DATE,
    address varchar(100) not null default '',
    status smallint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

CREATE TABLE IF NOT EXISTS dbo_trx_order(
    id bigserial not null primary key,
    invoice varchar not null,
    customer_id int8 not null references dbo_trx_customer(id),
    status smallint not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

CREATE TABLE IF NOT EXISTS dbo_dtl_order_detail(
    id bigserial not null primary key,
    order_id int8 not null references dbo_trx_order(id),
    product_id int8 not null references dbo_mst_product(id),
    quantity int4 not null default 0,
    total_price real not null default 0,
    status int not null default 1,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint dbo_dtl_order_detail_order_product_id_un unique (order_id, product_id)
);