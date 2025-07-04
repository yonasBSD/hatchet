// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: workers.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createWorker = `-- name: CreateWorker :one
INSERT INTO "Worker" (
    "id",
    "createdAt",
    "updatedAt",
    "tenantId",
    "name",
    "dispatcherId",
    "maxRuns",
    "webhookId",
    "type",
    "sdkVersion",
    "language",
    "languageVersion",
    "os",
    "runtimeExtra"
) VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1::uuid,
    $2::text,
    $3::uuid,
    $4::int,
    $5::uuid,
    $6::"WorkerType",
    $7::text,
    $8::"WorkerSDKS",
    $9::text,
    $10::text,
    $11::text
) RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused", type, "webhookId", language, "languageVersion", os, "runtimeExtra", "sdkVersion"
`

type CreateWorkerParams struct {
	Tenantid        pgtype.UUID    `json:"tenantid"`
	Name            string         `json:"name"`
	Dispatcherid    pgtype.UUID    `json:"dispatcherid"`
	MaxRuns         pgtype.Int4    `json:"maxRuns"`
	WebhookId       pgtype.UUID    `json:"webhookId"`
	Type            NullWorkerType `json:"type"`
	SdkVersion      pgtype.Text    `json:"sdkVersion"`
	Language        NullWorkerSDKS `json:"language"`
	LanguageVersion pgtype.Text    `json:"languageVersion"`
	Os              pgtype.Text    `json:"os"`
	RuntimeExtra    pgtype.Text    `json:"runtimeExtra"`
}

func (q *Queries) CreateWorker(ctx context.Context, db DBTX, arg CreateWorkerParams) (*Worker, error) {
	row := db.QueryRow(ctx, createWorker,
		arg.Tenantid,
		arg.Name,
		arg.Dispatcherid,
		arg.MaxRuns,
		arg.WebhookId,
		arg.Type,
		arg.SdkVersion,
		arg.Language,
		arg.LanguageVersion,
		arg.Os,
		arg.RuntimeExtra,
	)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.LastHeartbeatAt,
		&i.Name,
		&i.DispatcherId,
		&i.MaxRuns,
		&i.IsActive,
		&i.LastListenerEstablished,
		&i.IsPaused,
		&i.Type,
		&i.WebhookId,
		&i.Language,
		&i.LanguageVersion,
		&i.Os,
		&i.RuntimeExtra,
		&i.SdkVersion,
	)
	return &i, err
}

const deleteOldWorkerAssignEvents = `-- name: DeleteOldWorkerAssignEvents :one
WITH for_delete AS (
    SELECT
        "id"
    FROM "WorkerAssignEvent" wae
    WHERE
        wae."workerId" = $1::uuid
    ORDER BY wae."id" DESC
    OFFSET $2::int
    LIMIT $3::int + 1
), has_more AS (
    SELECT
        CASE
            WHEN COUNT(*) > $3 THEN TRUE
            ELSE FALSE
        END as has_more
    FROM for_delete
)
DELETE FROM "WorkerAssignEvent" wae
WHERE wae."id" IN (SELECT "id" FROM for_delete)
RETURNING
    (SELECT has_more FROM has_more) as has_more
`

type DeleteOldWorkerAssignEventsParams struct {
	Workerid pgtype.UUID `json:"workerid"`
	MaxRuns  int32       `json:"maxRuns"`
	Limit    int32       `json:"limit"`
}

// delete worker assign events outside of the first <maxRuns> events for a worker
func (q *Queries) DeleteOldWorkerAssignEvents(ctx context.Context, db DBTX, arg DeleteOldWorkerAssignEventsParams) (bool, error) {
	row := db.QueryRow(ctx, deleteOldWorkerAssignEvents, arg.Workerid, arg.MaxRuns, arg.Limit)
	var has_more bool
	err := row.Scan(&has_more)
	return has_more, err
}

const deleteOldWorkers = `-- name: DeleteOldWorkers :one
WITH for_delete AS (
    SELECT
        "id"
    FROM "Worker" w
    WHERE
        w."tenantId" = $1::uuid AND
        w."lastHeartbeatAt" < $2::timestamp
    LIMIT $3 + 1
), expired_with_limit AS (
    SELECT
        for_delete."id" as "id"
    FROM for_delete
    LIMIT $3
), has_more AS (
    SELECT
        CASE
            WHEN COUNT(*) > $3 THEN TRUE
            ELSE FALSE
        END as has_more
    FROM for_delete
), delete_events AS (
    DELETE FROM "WorkerAssignEvent" wae
    WHERE wae."workerId" IN (SELECT "id" FROM expired_with_limit)
    RETURNING wae."id"
)
DELETE FROM "Worker" w
WHERE w."id" IN (SELECT "id" FROM expired_with_limit)
RETURNING
    (SELECT has_more FROM has_more) as has_more
`

type DeleteOldWorkersParams struct {
	Tenantid            pgtype.UUID      `json:"tenantid"`
	Lastheartbeatbefore pgtype.Timestamp `json:"lastheartbeatbefore"`
	Limit               interface{}      `json:"limit"`
}

func (q *Queries) DeleteOldWorkers(ctx context.Context, db DBTX, arg DeleteOldWorkersParams) (bool, error) {
	row := db.QueryRow(ctx, deleteOldWorkers, arg.Tenantid, arg.Lastheartbeatbefore, arg.Limit)
	var has_more bool
	err := row.Scan(&has_more)
	return has_more, err
}

const deleteWorker = `-- name: DeleteWorker :one
DELETE FROM
  "Worker"
WHERE
  "id" = $1::uuid
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused", type, "webhookId", language, "languageVersion", os, "runtimeExtra", "sdkVersion"
`

func (q *Queries) DeleteWorker(ctx context.Context, db DBTX, id pgtype.UUID) (*Worker, error) {
	row := db.QueryRow(ctx, deleteWorker, id)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.LastHeartbeatAt,
		&i.Name,
		&i.DispatcherId,
		&i.MaxRuns,
		&i.IsActive,
		&i.LastListenerEstablished,
		&i.IsPaused,
		&i.Type,
		&i.WebhookId,
		&i.Language,
		&i.LanguageVersion,
		&i.Os,
		&i.RuntimeExtra,
		&i.SdkVersion,
	)
	return &i, err
}

const getWorkerActionsByWorkerId = `-- name: GetWorkerActionsByWorkerId :many
WITH inputs AS (
    SELECT UNNEST($2::UUID[]) AS "workerId"
)

SELECT
    w."id" AS "workerId",
    a."actionId" AS actionId
FROM "Worker" w
JOIN inputs i ON w."id" = i."workerId"
LEFT JOIN "_ActionToWorker" aw ON w.id = aw."B"
LEFT JOIN "Action" a ON aw."A" = a.id
WHERE
    a."tenantId" = $1::UUID
`

type GetWorkerActionsByWorkerIdParams struct {
	Tenantid  pgtype.UUID   `json:"tenantid"`
	Workerids []pgtype.UUID `json:"workerids"`
}

type GetWorkerActionsByWorkerIdRow struct {
	WorkerId pgtype.UUID `json:"workerId"`
	Actionid pgtype.Text `json:"actionid"`
}

func (q *Queries) GetWorkerActionsByWorkerId(ctx context.Context, db DBTX, arg GetWorkerActionsByWorkerIdParams) ([]*GetWorkerActionsByWorkerIdRow, error) {
	rows, err := db.Query(ctx, getWorkerActionsByWorkerId, arg.Tenantid, arg.Workerids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetWorkerActionsByWorkerIdRow
	for rows.Next() {
		var i GetWorkerActionsByWorkerIdRow
		if err := rows.Scan(&i.WorkerId, &i.Actionid); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkerById = `-- name: GetWorkerById :one
SELECT
    w.id, w."createdAt", w."updatedAt", w."deletedAt", w."tenantId", w."lastHeartbeatAt", w.name, w."dispatcherId", w."maxRuns", w."isActive", w."lastListenerEstablished", w."isPaused", w.type, w."webhookId", w.language, w."languageVersion", w.os, w."runtimeExtra", w."sdkVersion",
    ww."url" AS "webhookUrl",
    w."maxRuns" - (
        SELECT COUNT(*)
        FROM "SemaphoreQueueItem" sqi
        WHERE
            sqi."tenantId" = w."tenantId" AND
            sqi."workerId" = w."id"
    ) AS "remainingSlots"
FROM
    "Worker" w
LEFT JOIN
    "WebhookWorker" ww ON w."webhookId" = ww."id"
WHERE
    w."id" = $1::uuid
`

type GetWorkerByIdRow struct {
	Worker         Worker      `json:"worker"`
	WebhookUrl     pgtype.Text `json:"webhookUrl"`
	RemainingSlots int32       `json:"remainingSlots"`
}

func (q *Queries) GetWorkerById(ctx context.Context, db DBTX, id pgtype.UUID) (*GetWorkerByIdRow, error) {
	row := db.QueryRow(ctx, getWorkerById, id)
	var i GetWorkerByIdRow
	err := row.Scan(
		&i.Worker.ID,
		&i.Worker.CreatedAt,
		&i.Worker.UpdatedAt,
		&i.Worker.DeletedAt,
		&i.Worker.TenantId,
		&i.Worker.LastHeartbeatAt,
		&i.Worker.Name,
		&i.Worker.DispatcherId,
		&i.Worker.MaxRuns,
		&i.Worker.IsActive,
		&i.Worker.LastListenerEstablished,
		&i.Worker.IsPaused,
		&i.Worker.Type,
		&i.Worker.WebhookId,
		&i.Worker.Language,
		&i.Worker.LanguageVersion,
		&i.Worker.Os,
		&i.Worker.RuntimeExtra,
		&i.Worker.SdkVersion,
		&i.WebhookUrl,
		&i.RemainingSlots,
	)
	return &i, err
}

const getWorkerByWebhookId = `-- name: GetWorkerByWebhookId :one
SELECT
    id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused", type, "webhookId", language, "languageVersion", os, "runtimeExtra", "sdkVersion"
FROM
    "Worker"
WHERE
    "webhookId" = $1::uuid
    AND "tenantId" = $2::uuid
`

type GetWorkerByWebhookIdParams struct {
	Webhookid pgtype.UUID `json:"webhookid"`
	Tenantid  pgtype.UUID `json:"tenantid"`
}

func (q *Queries) GetWorkerByWebhookId(ctx context.Context, db DBTX, arg GetWorkerByWebhookIdParams) (*Worker, error) {
	row := db.QueryRow(ctx, getWorkerByWebhookId, arg.Webhookid, arg.Tenantid)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.LastHeartbeatAt,
		&i.Name,
		&i.DispatcherId,
		&i.MaxRuns,
		&i.IsActive,
		&i.LastListenerEstablished,
		&i.IsPaused,
		&i.Type,
		&i.WebhookId,
		&i.Language,
		&i.LanguageVersion,
		&i.Os,
		&i.RuntimeExtra,
		&i.SdkVersion,
	)
	return &i, err
}

const getWorkerForEngine = `-- name: GetWorkerForEngine :one
SELECT
    w."id" AS "id",
    w."tenantId" AS "tenantId",
    w."dispatcherId" AS "dispatcherId",
    d."lastHeartbeatAt" AS "dispatcherLastHeartbeatAt",
    w."isActive" AS "isActive",
    w."lastListenerEstablished" AS "lastListenerEstablished"
FROM
    "Worker" w
LEFT JOIN
    "Dispatcher" d ON w."dispatcherId" = d."id"
WHERE
    w."tenantId" = $1
    AND w."id" = $2
`

type GetWorkerForEngineParams struct {
	Tenantid pgtype.UUID `json:"tenantid"`
	ID       pgtype.UUID `json:"id"`
}

type GetWorkerForEngineRow struct {
	ID                        pgtype.UUID      `json:"id"`
	TenantId                  pgtype.UUID      `json:"tenantId"`
	DispatcherId              pgtype.UUID      `json:"dispatcherId"`
	DispatcherLastHeartbeatAt pgtype.Timestamp `json:"dispatcherLastHeartbeatAt"`
	IsActive                  bool             `json:"isActive"`
	LastListenerEstablished   pgtype.Timestamp `json:"lastListenerEstablished"`
}

func (q *Queries) GetWorkerForEngine(ctx context.Context, db DBTX, arg GetWorkerForEngineParams) (*GetWorkerForEngineRow, error) {
	row := db.QueryRow(ctx, getWorkerForEngine, arg.Tenantid, arg.ID)
	var i GetWorkerForEngineRow
	err := row.Scan(
		&i.ID,
		&i.TenantId,
		&i.DispatcherId,
		&i.DispatcherLastHeartbeatAt,
		&i.IsActive,
		&i.LastListenerEstablished,
	)
	return &i, err
}

const linkActionsToWorker = `-- name: LinkActionsToWorker :exec
INSERT INTO "_ActionToWorker" (
    "A",
    "B"
) SELECT
    unnest($1::uuid[]),
    $2::uuid
ON CONFLICT DO NOTHING
`

type LinkActionsToWorkerParams struct {
	Actionids []pgtype.UUID `json:"actionids"`
	Workerid  pgtype.UUID   `json:"workerid"`
}

func (q *Queries) LinkActionsToWorker(ctx context.Context, db DBTX, arg LinkActionsToWorkerParams) error {
	_, err := db.Exec(ctx, linkActionsToWorker, arg.Actionids, arg.Workerid)
	return err
}

const linkServicesToWorker = `-- name: LinkServicesToWorker :exec
INSERT INTO "_ServiceToWorker" (
    "A",
    "B"
)
VALUES (
    unnest($1::uuid[]),
    $2::uuid
)
ON CONFLICT DO NOTHING
`

type LinkServicesToWorkerParams struct {
	Services []pgtype.UUID `json:"services"`
	Workerid pgtype.UUID   `json:"workerid"`
}

func (q *Queries) LinkServicesToWorker(ctx context.Context, db DBTX, arg LinkServicesToWorkerParams) error {
	_, err := db.Exec(ctx, linkServicesToWorker, arg.Services, arg.Workerid)
	return err
}

const listDispatcherIdsForWorkers = `-- name: ListDispatcherIdsForWorkers :many
SELECT
    "id" as "workerId",
    "dispatcherId"
FROM
    "Worker"
WHERE
    "tenantId" = $1::uuid
    AND "id" = ANY($2::uuid[])
`

type ListDispatcherIdsForWorkersParams struct {
	Tenantid  pgtype.UUID   `json:"tenantid"`
	Workerids []pgtype.UUID `json:"workerids"`
}

type ListDispatcherIdsForWorkersRow struct {
	WorkerId     pgtype.UUID `json:"workerId"`
	DispatcherId pgtype.UUID `json:"dispatcherId"`
}

func (q *Queries) ListDispatcherIdsForWorkers(ctx context.Context, db DBTX, arg ListDispatcherIdsForWorkersParams) ([]*ListDispatcherIdsForWorkersRow, error) {
	rows, err := db.Query(ctx, listDispatcherIdsForWorkers, arg.Tenantid, arg.Workerids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListDispatcherIdsForWorkersRow
	for rows.Next() {
		var i ListDispatcherIdsForWorkersRow
		if err := rows.Scan(&i.WorkerId, &i.DispatcherId); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listManyWorkerLabels = `-- name: ListManyWorkerLabels :many
SELECT
    "id",
    "key",
    "intValue",
    "strValue",
    "createdAt",
    "updatedAt",
    "workerId"
FROM "WorkerLabel" wl
WHERE wl."workerId" = ANY($1::uuid[])
`

type ListManyWorkerLabelsRow struct {
	ID        int64            `json:"id"`
	Key       string           `json:"key"`
	IntValue  pgtype.Int4      `json:"intValue"`
	StrValue  pgtype.Text      `json:"strValue"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
	UpdatedAt pgtype.Timestamp `json:"updatedAt"`
	WorkerId  pgtype.UUID      `json:"workerId"`
}

func (q *Queries) ListManyWorkerLabels(ctx context.Context, db DBTX, workerids []pgtype.UUID) ([]*ListManyWorkerLabelsRow, error) {
	rows, err := db.Query(ctx, listManyWorkerLabels, workerids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListManyWorkerLabelsRow
	for rows.Next() {
		var i ListManyWorkerLabelsRow
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.IntValue,
			&i.StrValue,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.WorkerId,
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

const listRecentAssignedEventsForWorker = `-- name: ListRecentAssignedEventsForWorker :many
SELECT
    "workerId",
    "assignedStepRuns"
FROM
    "WorkerAssignEvent"
WHERE
    "workerId" = $1::uuid
ORDER BY "id" DESC
LIMIT
    COALESCE($2::int, 100)
`

type ListRecentAssignedEventsForWorkerParams struct {
	Workerid pgtype.UUID `json:"workerid"`
	Limit    pgtype.Int4 `json:"limit"`
}

type ListRecentAssignedEventsForWorkerRow struct {
	WorkerId         pgtype.UUID `json:"workerId"`
	AssignedStepRuns []byte      `json:"assignedStepRuns"`
}

func (q *Queries) ListRecentAssignedEventsForWorker(ctx context.Context, db DBTX, arg ListRecentAssignedEventsForWorkerParams) ([]*ListRecentAssignedEventsForWorkerRow, error) {
	rows, err := db.Query(ctx, listRecentAssignedEventsForWorker, arg.Workerid, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListRecentAssignedEventsForWorkerRow
	for rows.Next() {
		var i ListRecentAssignedEventsForWorkerRow
		if err := rows.Scan(&i.WorkerId, &i.AssignedStepRuns); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSemaphoreSlotsWithStateForWorker = `-- name: ListSemaphoreSlotsWithStateForWorker :many
SELECT
    sr."id" AS "stepRunId",
    sr."status" AS "status",
    s."actionId",
    sr."timeoutAt" AS "timeoutAt",
    sr."startedAt" AS "startedAt",
    jr."workflowRunId" AS "workflowRunId"
FROM
    "SemaphoreQueueItem" sqi
JOIN
    "StepRun" sr ON sr."id" = sqi."stepRunId"
JOIN
    "JobRun" jr ON sr."jobRunId" = jr."id"
JOIN
    "Step" s ON sr."stepId" = s."id"
WHERE
    sqi."tenantId" = $1::uuid
    AND sqi."workerId" = $2::uuid
ORDER BY
    sr."createdAt" DESC
LIMIT
    COALESCE($3::int, 100)
`

type ListSemaphoreSlotsWithStateForWorkerParams struct {
	Tenantid pgtype.UUID `json:"tenantid"`
	Workerid pgtype.UUID `json:"workerid"`
	Limit    pgtype.Int4 `json:"limit"`
}

type ListSemaphoreSlotsWithStateForWorkerRow struct {
	StepRunId     pgtype.UUID      `json:"stepRunId"`
	Status        StepRunStatus    `json:"status"`
	ActionId      string           `json:"actionId"`
	TimeoutAt     pgtype.Timestamp `json:"timeoutAt"`
	StartedAt     pgtype.Timestamp `json:"startedAt"`
	WorkflowRunId pgtype.UUID      `json:"workflowRunId"`
}

func (q *Queries) ListSemaphoreSlotsWithStateForWorker(ctx context.Context, db DBTX, arg ListSemaphoreSlotsWithStateForWorkerParams) ([]*ListSemaphoreSlotsWithStateForWorkerRow, error) {
	rows, err := db.Query(ctx, listSemaphoreSlotsWithStateForWorker, arg.Tenantid, arg.Workerid, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListSemaphoreSlotsWithStateForWorkerRow
	for rows.Next() {
		var i ListSemaphoreSlotsWithStateForWorkerRow
		if err := rows.Scan(
			&i.StepRunId,
			&i.Status,
			&i.ActionId,
			&i.TimeoutAt,
			&i.StartedAt,
			&i.WorkflowRunId,
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

const listWorkerLabels = `-- name: ListWorkerLabels :many
SELECT
    "id",
    "key",
    "intValue",
    "strValue",
    "createdAt",
    "updatedAt"
FROM "WorkerLabel" wl
WHERE wl."workerId" = $1::uuid
`

type ListWorkerLabelsRow struct {
	ID        int64            `json:"id"`
	Key       string           `json:"key"`
	IntValue  pgtype.Int4      `json:"intValue"`
	StrValue  pgtype.Text      `json:"strValue"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
	UpdatedAt pgtype.Timestamp `json:"updatedAt"`
}

func (q *Queries) ListWorkerLabels(ctx context.Context, db DBTX, workerid pgtype.UUID) ([]*ListWorkerLabelsRow, error) {
	rows, err := db.Query(ctx, listWorkerLabels, workerid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListWorkerLabelsRow
	for rows.Next() {
		var i ListWorkerLabelsRow
		if err := rows.Scan(
			&i.ID,
			&i.Key,
			&i.IntValue,
			&i.StrValue,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const listWorkersWithSlotCount = `-- name: ListWorkersWithSlotCount :many
SELECT
    workers.id, workers."createdAt", workers."updatedAt", workers."deletedAt", workers."tenantId", workers."lastHeartbeatAt", workers.name, workers."dispatcherId", workers."maxRuns", workers."isActive", workers."lastListenerEstablished", workers."isPaused", workers.type, workers."webhookId", workers.language, workers."languageVersion", workers.os, workers."runtimeExtra", workers."sdkVersion",
    ww."url" AS "webhookUrl",
    ww."id" AS "webhookId",
    workers."maxRuns" - (
        SELECT COUNT(*)
        FROM "SemaphoreQueueItem" sqi
        WHERE
            sqi."tenantId" = workers."tenantId" AND
            sqi."workerId" = workers."id"
    ) AS "remainingSlots"
FROM
    "Worker" workers
LEFT JOIN
    "WebhookWorker" ww ON workers."webhookId" = ww."id"
WHERE
    workers."tenantId" = $1
    AND (
        $2::text IS NULL OR
        workers."id" IN (
            SELECT "_ActionToWorker"."B"
            FROM "_ActionToWorker"
            INNER JOIN "Action" ON "Action"."id" = "_ActionToWorker"."A"
            WHERE "Action"."tenantId" = $1 AND "Action"."actionId" = $2::text
        )
    )
    AND (
        $3::timestamp IS NULL OR
        workers."lastHeartbeatAt" > $3::timestamp
    )
    AND (
        $4::boolean IS NULL OR
        workers."maxRuns" IS NULL OR
        ($4::boolean AND workers."maxRuns" > (
            SELECT COUNT(*)
            FROM "StepRun" srs
            WHERE srs."workerId" = workers."id" AND srs."status" = 'RUNNING'
        ))
    )
GROUP BY
    workers."id", ww."url", ww."id"
`

type ListWorkersWithSlotCountParams struct {
	Tenantid           pgtype.UUID      `json:"tenantid"`
	ActionId           pgtype.Text      `json:"actionId"`
	LastHeartbeatAfter pgtype.Timestamp `json:"lastHeartbeatAfter"`
	Assignable         pgtype.Bool      `json:"assignable"`
}

type ListWorkersWithSlotCountRow struct {
	Worker         Worker      `json:"worker"`
	WebhookUrl     pgtype.Text `json:"webhookUrl"`
	WebhookId      pgtype.UUID `json:"webhookId"`
	RemainingSlots int32       `json:"remainingSlots"`
}

func (q *Queries) ListWorkersWithSlotCount(ctx context.Context, db DBTX, arg ListWorkersWithSlotCountParams) ([]*ListWorkersWithSlotCountRow, error) {
	rows, err := db.Query(ctx, listWorkersWithSlotCount,
		arg.Tenantid,
		arg.ActionId,
		arg.LastHeartbeatAfter,
		arg.Assignable,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListWorkersWithSlotCountRow
	for rows.Next() {
		var i ListWorkersWithSlotCountRow
		if err := rows.Scan(
			&i.Worker.ID,
			&i.Worker.CreatedAt,
			&i.Worker.UpdatedAt,
			&i.Worker.DeletedAt,
			&i.Worker.TenantId,
			&i.Worker.LastHeartbeatAt,
			&i.Worker.Name,
			&i.Worker.DispatcherId,
			&i.Worker.MaxRuns,
			&i.Worker.IsActive,
			&i.Worker.LastListenerEstablished,
			&i.Worker.IsPaused,
			&i.Worker.Type,
			&i.Worker.WebhookId,
			&i.Worker.Language,
			&i.Worker.LanguageVersion,
			&i.Worker.Os,
			&i.Worker.RuntimeExtra,
			&i.Worker.SdkVersion,
			&i.WebhookUrl,
			&i.WebhookId,
			&i.RemainingSlots,
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

const updateWorker = `-- name: UpdateWorker :one
UPDATE
    "Worker"
SET
    "updatedAt" = CURRENT_TIMESTAMP,
    "dispatcherId" = coalesce($1::uuid, "dispatcherId"),
    "maxRuns" = coalesce($2::int, "maxRuns"),
    "lastHeartbeatAt" = coalesce($3::timestamp, "lastHeartbeatAt"),
    "isActive" = coalesce($4::boolean, "isActive"),
    "isPaused" = coalesce($5::boolean, "isPaused")
WHERE
    "id" = $6::uuid
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused", type, "webhookId", language, "languageVersion", os, "runtimeExtra", "sdkVersion"
`

type UpdateWorkerParams struct {
	DispatcherId    pgtype.UUID      `json:"dispatcherId"`
	MaxRuns         pgtype.Int4      `json:"maxRuns"`
	LastHeartbeatAt pgtype.Timestamp `json:"lastHeartbeatAt"`
	IsActive        pgtype.Bool      `json:"isActive"`
	IsPaused        pgtype.Bool      `json:"isPaused"`
	ID              pgtype.UUID      `json:"id"`
}

func (q *Queries) UpdateWorker(ctx context.Context, db DBTX, arg UpdateWorkerParams) (*Worker, error) {
	row := db.QueryRow(ctx, updateWorker,
		arg.DispatcherId,
		arg.MaxRuns,
		arg.LastHeartbeatAt,
		arg.IsActive,
		arg.IsPaused,
		arg.ID,
	)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.LastHeartbeatAt,
		&i.Name,
		&i.DispatcherId,
		&i.MaxRuns,
		&i.IsActive,
		&i.LastListenerEstablished,
		&i.IsPaused,
		&i.Type,
		&i.WebhookId,
		&i.Language,
		&i.LanguageVersion,
		&i.Os,
		&i.RuntimeExtra,
		&i.SdkVersion,
	)
	return &i, err
}

const updateWorkerActiveStatus = `-- name: UpdateWorkerActiveStatus :one
UPDATE "Worker"
SET
    "isActive" = $1::boolean,
    "lastListenerEstablished" = $2::timestamp
WHERE
    "id" = $3::uuid
    AND (
        "lastListenerEstablished" IS NULL
        OR "lastListenerEstablished" <= $2::timestamp
        )
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused", type, "webhookId", language, "languageVersion", os, "runtimeExtra", "sdkVersion"
`

type UpdateWorkerActiveStatusParams struct {
	Isactive                bool             `json:"isactive"`
	LastListenerEstablished pgtype.Timestamp `json:"lastListenerEstablished"`
	ID                      pgtype.UUID      `json:"id"`
}

func (q *Queries) UpdateWorkerActiveStatus(ctx context.Context, db DBTX, arg UpdateWorkerActiveStatusParams) (*Worker, error) {
	row := db.QueryRow(ctx, updateWorkerActiveStatus, arg.Isactive, arg.LastListenerEstablished, arg.ID)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.LastHeartbeatAt,
		&i.Name,
		&i.DispatcherId,
		&i.MaxRuns,
		&i.IsActive,
		&i.LastListenerEstablished,
		&i.IsPaused,
		&i.Type,
		&i.WebhookId,
		&i.Language,
		&i.LanguageVersion,
		&i.Os,
		&i.RuntimeExtra,
		&i.SdkVersion,
	)
	return &i, err
}

const updateWorkerHeartbeat = `-- name: UpdateWorkerHeartbeat :one
UPDATE
    "Worker"
SET
    "updatedAt" = CURRENT_TIMESTAMP,
    "lastHeartbeatAt" = $1::timestamp
WHERE
    "id" = $2::uuid
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused", type, "webhookId", language, "languageVersion", os, "runtimeExtra", "sdkVersion"
`

type UpdateWorkerHeartbeatParams struct {
	LastHeartbeatAt pgtype.Timestamp `json:"lastHeartbeatAt"`
	ID              pgtype.UUID      `json:"id"`
}

func (q *Queries) UpdateWorkerHeartbeat(ctx context.Context, db DBTX, arg UpdateWorkerHeartbeatParams) (*Worker, error) {
	row := db.QueryRow(ctx, updateWorkerHeartbeat, arg.LastHeartbeatAt, arg.ID)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.LastHeartbeatAt,
		&i.Name,
		&i.DispatcherId,
		&i.MaxRuns,
		&i.IsActive,
		&i.LastListenerEstablished,
		&i.IsPaused,
		&i.Type,
		&i.WebhookId,
		&i.Language,
		&i.LanguageVersion,
		&i.Os,
		&i.RuntimeExtra,
		&i.SdkVersion,
	)
	return &i, err
}

const updateWorkersByWebhookId = `-- name: UpdateWorkersByWebhookId :many
UPDATE "Worker"
SET "isActive" = $1::boolean
WHERE
  "tenantId" = $2::uuid AND
  "webhookId" = $3::uuid
RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "lastHeartbeatAt", name, "dispatcherId", "maxRuns", "isActive", "lastListenerEstablished", "isPaused", type, "webhookId", language, "languageVersion", os, "runtimeExtra", "sdkVersion"
`

type UpdateWorkersByWebhookIdParams struct {
	Isactive  bool        `json:"isactive"`
	Tenantid  pgtype.UUID `json:"tenantid"`
	Webhookid pgtype.UUID `json:"webhookid"`
}

func (q *Queries) UpdateWorkersByWebhookId(ctx context.Context, db DBTX, arg UpdateWorkersByWebhookIdParams) ([]*Worker, error) {
	rows, err := db.Query(ctx, updateWorkersByWebhookId, arg.Isactive, arg.Tenantid, arg.Webhookid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Worker
	for rows.Next() {
		var i Worker
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TenantId,
			&i.LastHeartbeatAt,
			&i.Name,
			&i.DispatcherId,
			&i.MaxRuns,
			&i.IsActive,
			&i.LastListenerEstablished,
			&i.IsPaused,
			&i.Type,
			&i.WebhookId,
			&i.Language,
			&i.LanguageVersion,
			&i.Os,
			&i.RuntimeExtra,
			&i.SdkVersion,
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

const upsertService = `-- name: UpsertService :one
INSERT INTO "Service" (
    "id",
    "createdAt",
    "updatedAt",
    "name",
    "tenantId"
)
VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1::text,
    $2::uuid
)
ON CONFLICT ("tenantId", "name") DO UPDATE
SET
    "updatedAt" = CURRENT_TIMESTAMP
WHERE
    "Service"."tenantId" = $2 AND "Service"."name" = $1::text
RETURNING id, "createdAt", "updatedAt", "deletedAt", name, description, "tenantId"
`

type UpsertServiceParams struct {
	Name     string      `json:"name"`
	Tenantid pgtype.UUID `json:"tenantid"`
}

func (q *Queries) UpsertService(ctx context.Context, db DBTX, arg UpsertServiceParams) (*Service, error) {
	row := db.QueryRow(ctx, upsertService, arg.Name, arg.Tenantid)
	var i Service
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Name,
		&i.Description,
		&i.TenantId,
	)
	return &i, err
}

const upsertWorkerLabel = `-- name: UpsertWorkerLabel :one
INSERT INTO "WorkerLabel" (
    "createdAt",
    "updatedAt",
    "workerId",
    "key",
    "intValue",
    "strValue"
) VALUES (
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1::uuid,
    $2::text,
    $3::int,
    $4::text
) ON CONFLICT ("workerId", "key") DO UPDATE
SET
    "updatedAt" = CURRENT_TIMESTAMP,
    "intValue" = $3::int,
    "strValue" = $4::text
RETURNING id, "createdAt", "updatedAt", "workerId", key, "strValue", "intValue"
`

type UpsertWorkerLabelParams struct {
	Workerid pgtype.UUID `json:"workerid"`
	Key      string      `json:"key"`
	IntValue pgtype.Int4 `json:"intValue"`
	StrValue pgtype.Text `json:"strValue"`
}

func (q *Queries) UpsertWorkerLabel(ctx context.Context, db DBTX, arg UpsertWorkerLabelParams) (*WorkerLabel, error) {
	row := db.QueryRow(ctx, upsertWorkerLabel,
		arg.Workerid,
		arg.Key,
		arg.IntValue,
		arg.StrValue,
	)
	var i WorkerLabel
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.WorkerId,
		&i.Key,
		&i.StrValue,
		&i.IntValue,
	)
	return &i, err
}
