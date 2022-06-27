CREATE TABLE "products" (
  "id" uuid PRIMARY KEY,
  "title" varchar(255) NOT NULL,
  "image" varchar(255) NOT NULL,
  "price" bigint NOT NULL,
  "sale" integer NOT NULL,
  "sale_old_price" bigint NOT NULL,
  "category_id" uuid NOT NULL,
  "type" varchar(255) NOT NULL,
  "subtype" varchar(255) NOT NULL,
  "description" varchar(255),
  "created_at" timestamp NOT NULL
);

CREATE TABLE "categories" (
  "id" uuid PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "surname" varchar(255) NOT NULL,
  "email" varchar(255) NOT NULL,
  "phone" varchar(255) NOT NULL,
  "role" varchar(255) NOT NULL,
  "password_hash" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL
);

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

COMMENT ON COLUMN "products"."price" IS 'must be positive';

COMMENT ON COLUMN "products"."sale_old_price" IS 'must be positive';

INSERT INTO "public"."categories" ("id", "name", "created_at") VALUES
('333241c8-3664-48a7-a1d6-749b5279e4f8', 'Мужчинам', '2022-01-12 13:08:21.32963+00');
INSERT INTO "public"."categories" ("id", "name", "created_at") VALUES
('9d8d6381-940c-4c21-990e-96d1fc7f54e5', 'Женщинам', '2022-01-12 13:08:21.32963+00');

INSERT INTO "public"."products" ("id", "title", "image", "price", "sale", "sale_old_price", "category_id", "type", "subtype", "description", "created_at") VALUES
('453b4f0f-1f56-4c57-b43d-7b79792450a7', 'Твидовый кардиган из хлопка', 'w1.webp', 749000, 0, 0, '9d8d6381-940c-4c21-990e-96d1fc7f54e5', 'Одежда', 'Старые-коллекции', '', '2022-01-12 13:08:21.32963+00');
INSERT INTO "public"."products" ("id", "title", "image", "price", "sale", "sale_old_price", "category_id", "type", "subtype", "description", "created_at") VALUES
('b07221f8-4133-4688-b2d6-d677f41f5b74', 'Объемный водоотталкивающий тренч', 'w2.webp', 499000, 50, 999000, '9d8d6381-940c-4c21-990e-96d1fc7f54e5', 'Одежда', 'Старые-коллекции', '', '2022-01-12 13:10:18.882593+00');
INSERT INTO "public"."products" ("id", "title", "image", "price", "sale", "sale_old_price", "category_id", "type", "subtype", "description", "created_at") VALUES
('7c21f349-5e20-453b-83ca-c3279296f98a', 'Жилет из трикотажа в рубчик', 'w3.webp', 329000, 0, 0, '9d8d6381-940c-4c21-990e-96d1fc7f54e5', 'Одежда', 'Старые-коллекции', '', '2022-01-12 13:11:02.622235+00');
INSERT INTO "public"."products" ("id", "title", "image", "price", "sale", "sale_old_price", "category_id", "type", "subtype", "description", "created_at") VALUES
('96a7193a-403d-4e01-94e6-c02c5bcb61f1', 'Хлопковая рубашка в полоску', 'w4.webp', 359000, 0, 0, '9d8d6381-940c-4c21-990e-96d1fc7f54e5', 'Одежда', 'Вышевка', '', '2022-01-12 13:16:55.558842+00'),
('e942a6a2-2aed-40e1-a305-350f4d7e2813', 'Объемная рубашку в клетку', 'w5.webp', 359000, 0, 0, '9d8d6381-940c-4c21-990e-96d1fc7f54e5', 'Одежда', 'Вышевка', '', '2022-01-12 13:17:21.866024+00'),
('72c07fab-dd14-47fa-b478-d59255dcf8dd', 'Джинсы с завышеной талией Wideleg', 'w6.webp', 169000, 0, 0, '9d8d6381-940c-4c21-990e-96d1fc7f54e5', 'Одежда', 'Средиземноморье', '', '2022-01-12 13:17:58.430896+00'),
('ee8f7092-8bf7-4ae3-bf29-484cba7eb7f9', 'Джинсы с завышеной талией Wideleg', 'w7.webp', 359000, 0, 0, '9d8d6381-940c-4c21-990e-96d1fc7f54e5', 'Одежда', 'Средиземноморье', '', '2022-01-12 13:18:34.087917+00'),
('940f894a-37bc-4352-9782-dc68d77b7cf2', 'Кожаные сапоги', 'w8.webp', 499000, 0, 0, '9d8d6381-940c-4c21-990e-96d1fc7f54e5', 'Обувь', 'Средиземноморье', '', '2022-01-12 13:19:10.755827+00'),
('e1b5ef67-a82a-427c-a88c-a507e02e9aea', 'Хлопковые носки в рубчик', 'w9.webp', 99000, 0, 0, '9d8d6381-940c-4c21-990e-96d1fc7f54e5', 'Обувь', 'Средиземноморье', '', '2022-01-12 13:19:43.175537+00'),
('528415ca-8244-4bb4-b190-0c0fa82fc5d0', 'Джинсовая куртка из хлопка', 'w10.webp', 359000, 0, 0, '9d8d6381-940c-4c21-990e-96d1fc7f54e5', 'Одежда', 'Старые-коллекции', '', '2022-01-12 13:20:17.121218+00'),
('0e362be1-fa04-4a5a-a021-aaeb7e2c0b6a', 'Фактурная трикотажная футболка из хлопка', 'm1.webp', 249000, 0, 0, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Одежда', 'Старые-коллекции', '', '2022-01-12 13:20:55.219142+00'),
('bcdf2e5b-1ccc-4429-81c3-4bbb32190b94', 'Футболка тай-дай relaxed fit', 'm2.webp', 229000, 0, 0, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Одежда', 'Средиземноморье', '', '2022-01-12 13:21:59.959962+00'),
('a0b7df33-1f26-49b3-9bda-916af6a3684a', 'Хлопковая футболка стрейч', 'm3.webp', 199000, 0, 0, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Одежда', 'Спорт', '', '2022-01-12 13:22:30.073053+00'),
('3486287f-10f8-498c-bad3-92dab022c616', 'Рубашка из фланели в клетку', 'm4.webp', 499000, 0, 0, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Одежда', 'Старые-коллекции', '', '2022-01-12 13:24:10.046077+00'),
('7b8fb5dd-3abf-47d9-a278-1667c9ad349d', 'Хлопковая рубашка с терморегуляцией', 'm5.webp', 249000, 31, 359000, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Одежда', 'Старые-коллекции', '', '2022-01-12 13:24:58.352971+00'),
('9a1bf731-9dad-4665-b396-fdc671ba24d7', 'Рубашка relaxed-fit из хлопка', 'm6.webp', 499000, 0, 0, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Одежда', 'Спорт', '', '2022-01-12 13:25:36.75912+00'),
('088a257d-8fb8-456d-8582-cf3bd06f90b0', 'Хлопковые брюки в грузовом стиле', 'm7.webp', 499000, 0, 0, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Одежда', 'Старые-коллекции', '', '2022-01-12 13:26:12.28865+00'),
('8aceffc7-d49f-4057-9e15-bf496fa1eb4f', 'Непромокаемый тренч из нейлона', 'm8.webp', 1079000, 10, 1199000, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Одежда', 'Старые-коллекции', '', '2022-01-12 13:27:00.113489+00'),
('2153537c-6c62-489a-a5c4-55bdb4aaa714', 'Однотонные кроссовки из кожи', 'm9.webp', 359000, 28, 499000, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Обувь', 'Спорт', '', '2022-01-12 13:27:52.429254+00'),
('131a95ae-5e4d-4c64-a0e4-d72d9e502268', 'Укороченные брюки-карго из хлопка', 'm10.webp', 449000, 0, 0, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Одежда', 'Вышевка', '', '2022-01-12 13:28:33.449473+00'),
('ef9aecdf-0860-49b0-b924-e4a418e19a00', 'Брюки-джоггеры из хлопка', 'm11.webp', 399000, 0, 0, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Одежда', 'Старые-коллекции', '', '2022-01-12 13:29:11.668989+00'),
('244f4012-20c4-417f-892e-48166ec4610d', 'Стеганая непромокаемая парка с капюшоном', 'm12.webp', 1299000, 0, 0, '333241c8-3664-48a7-a1d6-749b5279e4f8', 'Одежда', 'Старые-коллекции', '', '2022-01-12 13:29:43.317067+00');