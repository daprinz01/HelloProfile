-- name: GetCallToActions :many
select * from call_to_action;

-- name: GetCallToAction :one
select * from call_to_action where id=$1 limit 1;
