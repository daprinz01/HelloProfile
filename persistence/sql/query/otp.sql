-- name: GetAllOtp :many
select * from otp;

-- name: GetOtp :one
select * from otp where user_id = (select a.Id from users a where a.username = $1 or email = $1) and otp_token = $2;

-- name: CreateOtp :exec
insert into otp (
    user_id,
    otp_token,
    is_sms_preferred,
    is_email_preferred,
    purpose
)
values ($1, $2, $3, $4, $5);

-- name: DeleteOtp :exec
delete from otp where user_id = $1 and otp_token = $2;