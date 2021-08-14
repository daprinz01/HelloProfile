create table refresh_token(
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4() UNIQUE,
    "user_id" uuid not null,
    "token" varchar not null UNIQUE,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

alter table "refresh_token" add foreign key ("user_id") references "users" ("id");
