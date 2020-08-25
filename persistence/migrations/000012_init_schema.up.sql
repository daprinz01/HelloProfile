create table refresh_token(
    "id" bigserial primary key UNIQUE,
    "user_id" bigserial not null UNIQUE FOREIGN KEY REFERENCES (users),
    "token" varchar not null UNIQUE,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
