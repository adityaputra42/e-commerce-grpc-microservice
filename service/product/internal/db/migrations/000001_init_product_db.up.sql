-- Tabel kategori
CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(100) NOT NULL,
  "icon" varchar(255) NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz DEFAULT NULL
);

-- Tabel produk
CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "category_id" bigint NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "images" text[] NOT NULL,
  "rating" real NOT NULL DEFAULT 0,
  "price" double precision NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz DEFAULT NULL
);

CREATE INDEX ON "products" ("category_id");
COMMENT ON COLUMN "products"."price" IS 'must be positive';
ALTER TABLE "products"
  ADD CONSTRAINT "fk_products_category"
  FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

-- Tabel varian warna
CREATE TABLE "color_varians" (
  "id" bigserial PRIMARY KEY,
  "product_id" bigint NOT NULL,
  "name" varchar(100) NOT NULL,
  "color" varchar(100) NOT NULL,
  "images" text[] NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz DEFAULT NULL
);

CREATE INDEX ON "color_varians" ("product_id");
ALTER TABLE "color_varians"
  ADD CONSTRAINT "fk_color_varians_product"
  FOREIGN KEY ("product_id") REFERENCES "products" ("id");

-- Tabel varian ukuran
CREATE TABLE "size_varians" (
  "id" bigserial PRIMARY KEY,
  "color_varian_id" bigint NOT NULL,
  "size" varchar(50) NOT NULL,
  "stock" bigint NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz DEFAULT NULL
);

CREATE INDEX ON "size_varians" ("color_varian_id");
ALTER TABLE "size_varians"
  ADD CONSTRAINT "fk_size_varians_color_varian"
  FOREIGN KEY ("color_varian_id") REFERENCES "color_varians" ("id");
