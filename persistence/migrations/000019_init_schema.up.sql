create table otp(
    id bigserial primary key,
    user_id bigserial,
    otp_token varchar null,
    created_at TIMESTAMPtz not NULL DEFAULT (now()),
    is_sms_preferred boolean not null default FALSE,
    is_email_preferred boolean not null default TRUE,
    purpose varchar null,

    FOREIGN KEY (user_id) REFERENCES users(id),
);


Create INDEX on otp (user_id, otp_token);