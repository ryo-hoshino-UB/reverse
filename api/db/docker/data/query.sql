-- name: CreateGame :exec
INSERT INTO games (started_at) VALUES (NOW());
SELECT LAST_INSERT_ID() AS id, started_at FROM games WHERE id = LAST_INSERT_ID();

-- name: ListGames :many
SELECT id, started_at FROM games;

-- name: CreateTurn :exec
INSERT INTO turns (game_id, turn_count, next_disc, end_at) 
VALUES ($1, $2, $3, NOW());
SELECT id, game_id, turn_count, next_disc, end_at FROM turns WHERE id = LAST_INSERT_ID();


-- name: ListTurns :many
SELECT id, game_id, turn_count, next_disc, end_at FROM turns;

-- name: CreateMove :exec
INSERT INTO moves (turn_id, disc, x, y) 
VALUES ($1, $2, $3, $4);
SELECT id, turn_id, disc, x, y FROM moves WHERE id = LAST_INSERT_ID();

-- name: ListMoves :many
SELECT id, turn_id, disc, x, y FROM moves;

-- name: CreateSquare :exec
INSERT INTO squares (turn_id, x, y, disc) 
VALUES ($1, $2, $3, $4);
SELECT id, turn_id, x, y, disc FROM squares WHERE id = LAST_INSERT_ID();

-- name: ListSquares :many
SELECT id, turn_id, x, y, disc FROM squares;

-- name: CreateGameResult :exec
INSERT INTO game_results (game_id, winner_disc, end_at) 
VALUES ($1, $2, NOW());
SELECT id, game_id, winner_disc, end_at FROM game_results WHERE id = LAST_INSERT_ID();

-- name: ListGameResults :many
SELECT id, game_id, winner_disc, end_at FROM game_results;
