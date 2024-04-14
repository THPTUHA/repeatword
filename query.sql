-- name: GetVobsCollection :many
SELECT v.id, v.word FROM
collections c, collection_words cw, vobs v
WHERE c.id = cw.collection_id 
AND cw.vob_id = v.id
AND c.id = ? ORDER BY RAND();

-- name: GetMeans :many
SELECT * FROM means m
WHERE m.part_id IN (sqlc.slice('part_ids'));

-- name: GetExamples :many
SELECT * FROM examples e
WHERE e.mean_id IN (sqlc.slice('mean_ids'));

-- name: SetWord :exec
CALL SetWord(?,?);

-- name: GetVobsRandom :one
SELECT GetVobsRandom(?, ?);