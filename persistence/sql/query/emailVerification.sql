-- name: GetEmailVerifications :many
select * from email_verification;

-- name: GetEmailVerification :one
select * from email_verification where otp=$1;

-- name: CreateEmailVerification :exec
insert into email_verification(
    email, otp
) values ($1, $2);

-- name: DeleteEmailVerification :exec
delete from email_verification where otp=$1;