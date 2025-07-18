// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: webhook_workers.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createWebhookWorker = `-- name: CreateWebhookWorker :one
INSERT INTO "WebhookWorker" (
    "id",
    "createdAt",
    "updatedAt",
    "name",
    "secret",
    "url",
    "tenantId",
    "tokenId",
    "tokenValue",
    "deleted"
)
VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1::text,
    $2::text,
    $3::text,
    $4::uuid,
    $5::uuid,
    $6::text,
    coalesce($7::boolean, false)
)
RETURNING id, "createdAt", "updatedAt", name, secret, url, "tokenValue", deleted, "tokenId", "tenantId"
`

type CreateWebhookWorkerParams struct {
	Name       string      `json:"name"`
	Secret     string      `json:"secret"`
	Url        string      `json:"url"`
	Tenantid   pgtype.UUID `json:"tenantid"`
	TokenId    pgtype.UUID `json:"tokenId"`
	TokenValue pgtype.Text `json:"tokenValue"`
	Deleted    pgtype.Bool `json:"deleted"`
}

func (q *Queries) CreateWebhookWorker(ctx context.Context, db DBTX, arg CreateWebhookWorkerParams) (*WebhookWorker, error) {
	row := db.QueryRow(ctx, createWebhookWorker,
		arg.Name,
		arg.Secret,
		arg.Url,
		arg.Tenantid,
		arg.TokenId,
		arg.TokenValue,
		arg.Deleted,
	)
	var i WebhookWorker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Secret,
		&i.Url,
		&i.TokenValue,
		&i.Deleted,
		&i.TokenId,
		&i.TenantId,
	)
	return &i, err
}

const getWebhookWorkerByID = `-- name: GetWebhookWorkerByID :one
SELECT id, "createdAt", "updatedAt", name, secret, url, "tokenValue", deleted, "tokenId", "tenantId"
FROM "WebhookWorker"
WHERE "id" = $1::uuid
`

func (q *Queries) GetWebhookWorkerByID(ctx context.Context, db DBTX, id pgtype.UUID) (*WebhookWorker, error) {
	row := db.QueryRow(ctx, getWebhookWorkerByID, id)
	var i WebhookWorker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Secret,
		&i.Url,
		&i.TokenValue,
		&i.Deleted,
		&i.TokenId,
		&i.TenantId,
	)
	return &i, err
}

const hardDeleteWebhookWorker = `-- name: HardDeleteWebhookWorker :exec
DELETE FROM "WebhookWorker"
WHERE
  "id" = $1::uuid
  AND "tenantId" = $2::uuid
`

type HardDeleteWebhookWorkerParams struct {
	ID       pgtype.UUID `json:"id"`
	Tenantid pgtype.UUID `json:"tenantid"`
}

func (q *Queries) HardDeleteWebhookWorker(ctx context.Context, db DBTX, arg HardDeleteWebhookWorkerParams) error {
	_, err := db.Exec(ctx, hardDeleteWebhookWorker, arg.ID, arg.Tenantid)
	return err
}

const insertWebhookWorkerRequest = `-- name: InsertWebhookWorkerRequest :exec
WITH delete_old AS (
    -- Delete old requests
    DELETE FROM "WebhookWorkerRequest"
    WHERE "webhookWorkerId" = $1::uuid
    AND "createdAt" < NOW() - INTERVAL '15 minutes'
)
INSERT INTO "WebhookWorkerRequest" (
    "id",
    "createdAt",
    "webhookWorkerId",
    "method",
    "statusCode"
) VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    $1::uuid,
    $2::"WebhookWorkerRequestMethod",
    $3::integer
)
`

type InsertWebhookWorkerRequestParams struct {
	Webhookworkerid pgtype.UUID                `json:"webhookworkerid"`
	Method          WebhookWorkerRequestMethod `json:"method"`
	Statuscode      int32                      `json:"statuscode"`
}

func (q *Queries) InsertWebhookWorkerRequest(ctx context.Context, db DBTX, arg InsertWebhookWorkerRequestParams) error {
	_, err := db.Exec(ctx, insertWebhookWorkerRequest, arg.Webhookworkerid, arg.Method, arg.Statuscode)
	return err
}

const listActiveWebhookWorkers = `-- name: ListActiveWebhookWorkers :many
SELECT id, "createdAt", "updatedAt", name, secret, url, "tokenValue", deleted, "tokenId", "tenantId"
FROM "WebhookWorker"
WHERE "tenantId" = $1::uuid AND "deleted" = false
`

func (q *Queries) ListActiveWebhookWorkers(ctx context.Context, db DBTX, tenantid pgtype.UUID) ([]*WebhookWorker, error) {
	rows, err := db.Query(ctx, listActiveWebhookWorkers, tenantid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*WebhookWorker
	for rows.Next() {
		var i WebhookWorker
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Secret,
			&i.Url,
			&i.TokenValue,
			&i.Deleted,
			&i.TokenId,
			&i.TenantId,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listWebhookWorkerRequests = `-- name: ListWebhookWorkerRequests :many
SELECT id, "createdAt", "webhookWorkerId", method, "statusCode"
FROM "WebhookWorkerRequest"
WHERE "webhookWorkerId" = $1::uuid
ORDER BY "createdAt" DESC
LIMIT 50
`

func (q *Queries) ListWebhookWorkerRequests(ctx context.Context, db DBTX, webhookworkerid pgtype.UUID) ([]*WebhookWorkerRequest, error) {
	rows, err := db.Query(ctx, listWebhookWorkerRequests, webhookworkerid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*WebhookWorkerRequest
	for rows.Next() {
		var i WebhookWorkerRequest
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.WebhookWorkerId,
			&i.Method,
			&i.StatusCode,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listWebhookWorkersByPartitionId = `-- name: ListWebhookWorkersByPartitionId :many
WITH tenants AS (
    SELECT
        "id"
    FROM
        "Tenant"
    WHERE
        "workerPartitionId" = $1::text OR
        "workerPartitionId" IS NULL
), update_partition AS (
    UPDATE
        "TenantWorkerPartition"
    SET
        "lastHeartbeat" = NOW()
    WHERE
        "id" = $1::text
)
SELECT id, "createdAt", "updatedAt", name, secret, url, "tokenValue", deleted, "tokenId", "tenantId"
FROM "WebhookWorker"
WHERE "tenantId" IN (SELECT "id" FROM tenants)
`

func (q *Queries) ListWebhookWorkersByPartitionId(ctx context.Context, db DBTX, workerpartitionid string) ([]*WebhookWorker, error) {
	rows, err := db.Query(ctx, listWebhookWorkersByPartitionId, workerpartitionid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*WebhookWorker
	for rows.Next() {
		var i WebhookWorker
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Secret,
			&i.Url,
			&i.TokenValue,
			&i.Deleted,
			&i.TokenId,
			&i.TenantId,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const softDeleteWebhookWorker = `-- name: SoftDeleteWebhookWorker :exec
UPDATE "WebhookWorker"
SET
  "deleted" = true,
  "updatedAt" = CURRENT_TIMESTAMP
WHERE
  "id" = $1::uuid
  AND "tenantId" = $2::uuid
`

type SoftDeleteWebhookWorkerParams struct {
	ID       pgtype.UUID `json:"id"`
	Tenantid pgtype.UUID `json:"tenantid"`
}

func (q *Queries) SoftDeleteWebhookWorker(ctx context.Context, db DBTX, arg SoftDeleteWebhookWorkerParams) error {
	_, err := db.Exec(ctx, softDeleteWebhookWorker, arg.ID, arg.Tenantid)
	return err
}

const updateWebhookWorkerToken = `-- name: UpdateWebhookWorkerToken :one
UPDATE "WebhookWorker"
SET
    "updatedAt" = CURRENT_TIMESTAMP,
    "tokenValue" = COALESCE($1::text, "tokenValue"),
    "tokenId" = COALESCE($2::uuid, "tokenId")
WHERE
    "id" = $3::uuid
    AND "tenantId" = $4::uuid
RETURNING id, "createdAt", "updatedAt", name, secret, url, "tokenValue", deleted, "tokenId", "tenantId"
`

type UpdateWebhookWorkerTokenParams struct {
	TokenValue pgtype.Text `json:"tokenValue"`
	TokenId    pgtype.UUID `json:"tokenId"`
	ID         pgtype.UUID `json:"id"`
	Tenantid   pgtype.UUID `json:"tenantid"`
}

func (q *Queries) UpdateWebhookWorkerToken(ctx context.Context, db DBTX, arg UpdateWebhookWorkerTokenParams) (*WebhookWorker, error) {
	row := db.QueryRow(ctx, updateWebhookWorkerToken,
		arg.TokenValue,
		arg.TokenId,
		arg.ID,
		arg.Tenantid,
	)
	var i WebhookWorker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Secret,
		&i.Url,
		&i.TokenValue,
		&i.Deleted,
		&i.TokenId,
		&i.TenantId,
	)
	return &i, err
}
