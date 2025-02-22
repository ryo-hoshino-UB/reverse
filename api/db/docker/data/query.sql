-- name: CreateGame :execresult
INSERT INTO games (started_at) VALUES (NOW());
SELECT LAST_INSERT_ID() AS id, started_at FROM games WHERE id = LAST_INSERT_ID();

-- name: ListGames :many
SELECT id, started_at FROM games;

-- name: GetGameByID :one
SELECT id, started_at FROM games WHERE id = ?;

-- name: GetLatestGame :one
SELECT id, started_at FROM games order by id desc limit 1;

-- name: CreateTurn :execresult
INSERT INTO turns (game_id, turn_count, next_disc, end_at) 
VALUES (?, ?, ?, NOW());
SELECT id, game_id, turn_count, next_disc, end_at FROM turns WHERE id = LAST_INSERT_ID();

-- name: ListTurns :many
SELECT id, game_id, turn_count, next_disc, end_at FROM turns;

-- name: GetTurnByID :one
SELECT id, game_id, turn_count, next_disc, end_at FROM turns WHERE id = ?;

-- name: GetTurnByGameIDAndTurnCount :one
SELECT id, game_id, turn_count, next_disc, end_at FROM turns WHERE game_id = ? and turn_count = ?;

-- name: CreateMove :execresult
INSERT INTO moves (turn_id, disc, x, y) 
VALUES (?, ?, ?, ?);
SELECT id, turn_id, disc, x, y FROM moves WHERE id = LAST_INSERT_ID();

-- name: ListMoves :many
SELECT id, turn_id, disc, x, y FROM moves;

-- name: GetMoveByID :one
SELECT id, turn_id, disc, x, y FROM moves WHERE id = ?;

-- name: GetMoveByTurnID :one
SELECT id, turn_id, disc, x, y FROM moves WHERE turn_id = ?;

-- name: CreateSquare :execresult
INSERT INTO squares (turn_id, x, y, disc) 
VALUES (?, ?, ?, ?);
SELECT id, turn_id, x, y, disc FROM squares WHERE id = LAST_INSERT_ID();

-- name: ListSquares :many
SELECT id, turn_id, x, y, disc FROM squares;

-- name: GetSquareByID :one
SELECT id, turn_id, x, y, disc FROM squares WHERE id = ?;

-- name: GetSquaresByTurnID :many
SELECT id, turn_id, x, y, disc FROM squares WHERE turn_id = ?;

-- name: CreateGameResult :execresult
INSERT INTO game_results (game_id, winner_disc, end_at) 
VALUES (?, ?, NOW());
SELECT id, game_id, winner_disc, end_at FROM game_results WHERE id = LAST_INSERT_ID();

-- name: ListGameResults :many
SELECT id, game_id, winner_disc, end_at FROM game_results;

-- name: GetGameResultByID :one
SELECT id, game_id, winner_disc, end_at FROM game_results WHERE id = ?;

-- name: GetGameResultByGameID :one
SELECT id, game_id, winner_disc, end_at FROM game_results WHERE game_id = ?;