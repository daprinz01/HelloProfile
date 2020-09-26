create table refresh_token(
    "id" bigserial primary key UNIQUE,
    "user_id" bigserial not null,
    "token" varchar not null UNIQUE,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

alter table "refresh_token" add foreign key ("user_id") references "users" ("id");
