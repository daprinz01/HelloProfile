create table email_verification(
    "id" uuid PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    "email" varchar,
    "otp" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
create index on email_verification(otp);

alter table email_verification add constraint uc_email_verification UNIQUE(otp);