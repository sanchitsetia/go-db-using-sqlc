-- name: GetUser :one
select id,username,email from users where id = $1;

-- name: ListUsers :many
select id,username,email from users order by username;

-- name: CreateUser :one
insert into users (username,email,password) values($1,$2,$3) RETURNING id,username,email;

-- name: DeleteUser :exec
delete from users where id = $1;