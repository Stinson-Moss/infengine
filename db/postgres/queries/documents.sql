-- Document CRUD

-- name: CreateDocument :one
INSERT INTO Documents (Guid, Title, Description, Content)
VALUES ($1, $2, $3, $4)
RETURNING Id, Guid, Title, Description, Content, Created;

-- name: GetDocumentByGuid :one
SELECT Id, Guid, Title, Description, Content, Created 
FROM Documents 
WHERE Guid = $1 LIMIT 1;

-- name: UpdateDocument :one
UPDATE Documents
SET Title = $2, Description = $3, Content = $4
WHERE Guid = $1
RETURNING Id, Guid, Title, Description, Content, Created;

-- name: DeleteDocument :exec
DELETE FROM Documents 
WHERE Guid = $1;

-- Author, Tag, Link creating when Documents are creating
-- name: GetOrCreateAuthor :one
INSERT INTO Authors (Name) 
VALUES ($1)
ON CONFLICT (Name) DO UPDATE SET Name = EXCLUDED.Name
RETURNING Id, Name;

-- name: GetOrCreateTag :one
INSERT INTO Tags (Name) 
VALUES ($1)
ON CONFLICT (Name) DO UPDATE SET Name = EXCLUDED.Name
RETURNING Id, Name;

-- name: GetOrCreateLink :one
INSERT INTO Links (Link) 
VALUES ($1)
ON CONFLICT (Link) DO UPDATE SET Link = EXCLUDED.Link
RETURNING Id, Link;


-- These will be in a transaction
-- name: LinkDocumentToAuthor :exec
INSERT INTO Document_Authors (DocumentId, AuthorId)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: LinkDocumentToTag :exec
INSERT INTO Document_Tags (DocumentId, TagId)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: LinkDocumentToLink :exec
INSERT INTO Document_Links (DocumentId, LinkId)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- AI queriy method tools as well as ui interface method tools
-- name: GetDocumentsByAuthors :many
SELECT d.Id, d.Guid, d.Title, d.Description, d.Content, d.Created
FROM Documents d
JOIN Document_Authors da ON d.Id = da.DocumentId
JOIN Authors a ON da.AuthorId = a.Id
WHERE a.Name = ANY($1::text[]);

-- name: GetDocumentsByTimeRange :many
SELECT Id, Guid, Title, Description, Content, Created 
FROM Documents
WHERE Created >= $1 AND Created <= $2
ORDER BY Created DESC;

-- name: GetDocumentsByTags :many
SELECT d.Id, d.Guid, d.Title, d.Description, d.Content, d.Created
FROM Documents d
JOIN Document_Tags dt ON d.Id = dt.DocumentId
JOIN Tags t ON dt.TagId = t.Id
WHERE t.Name = ANY($1::text[]);

-- name: GetDocumentsByRelatedLinks :many
SELECT DISTINCT d.Id, d.Guid, d.Title, d.Description, d.Content, d.Created
FROM Documents d
JOIN Document_Links dl ON d.Id = dl.DocumentId
JOIN Links l ON dl.LinkId = l.Id
WHERE l.Link LIKE $1;