-- name: GetStates :many
select * from states;

-- name: GetState :one
select * from states where name = $1;

-- name: GetStatesByCountry :many
select * from states a where a.country_id = (select b.id from countries b where b.name = $1);

-- name: CreateState :exec
insert into states (
    name,
    b.country_id
)
values(
    $1,
    (select a.Id from countries a where a.name = $2)
);

-- name: UpdateState :exec
update states set name = $1 where name = $2;

-- name: DeleteState :exec
delete from states where name = $1;
