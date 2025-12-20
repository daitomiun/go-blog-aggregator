-- name: CreatePost :exec
INSERT INTO POSTS (ID, CREATED_AT, UPDATED_AT, TITLE, URL, DESCRIPTION, PUBLISHED_AT, FEED_ID)
VALUES ($1, $2,$3,$4,$5,$6, $7, $8);

-- name: GetPostsByUserId :many
select 
    posts.ID,
    posts.created_at,
    posts.updated_at,
    posts.title,
    posts.url,
    posts.description,
    posts.published_at,
    posts.feed_id
from posts 
inner join feed_follows on feed_follows.feed_id = posts.feed_id
where feed_follows.user_id = $1 
order by posts.published_at desc 
limit $2;


-- name: GetPostByUrl :one
SELECT
    ID,
    URL
FROM POSTS WHERE URL = $1 LIMIT 1;
