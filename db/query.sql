-- Active: 1692091993022@@127.0.0.1@3306@ks-laboratory-backend

-- name: GetUserWithPosts :many

select
    u.id as user_id,
    u.name as user_name,
    p.id as post_id,
    p.title as post_title,
    p.body as post_body,
    p.user_id as post_user_id
from users as u
    left join posts as p on u.id = p.user_id
where u.id = sqlc.arg(user_id);
