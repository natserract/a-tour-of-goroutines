CREATE TABLE "categories" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NULL,
    "deletedAt" timestamptz NULL,
    CONSTRAINT product_category_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX product_category ON "categories" USING btree ("name") WHERE ("deletedAt" IS NOT NULL);
