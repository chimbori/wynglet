-- name: CreateFormSubmission :one
INSERT INTO form_submissions (form_id, domain, ip_address, form_data, is_spam, email_sent_at)
  VALUES ($1, $2, $3, $4, $5, $6)
  RETURNING _id, form_id, submitted_at, domain, ip_address, form_data, is_spam, email_sent_at;

-- name: ListForms :many
SELECT form_id, COUNT(*) as count
  FROM form_submissions
  GROUP BY form_id
  ORDER BY form_id;

-- name: GetFormSubmission :one
SELECT _id, form_id, submitted_at, domain, ip_address, form_data, is_spam, email_sent_at
  FROM form_submissions
  WHERE _id = @_id;

-- name: ListFormSubmissions :many
SELECT _id, form_id, submitted_at, domain, ip_address, form_data, is_spam, email_sent_at
  FROM form_submissions
  WHERE form_id = @form_id
  ORDER BY submitted_at DESC
  LIMIT $1 OFFSET $2;

-- name: CountFormSubmissions :one
SELECT COUNT(*)
  FROM form_submissions
  WHERE form_id = @form_id;

-- name: DeleteFormSubmission :exec
DELETE FROM form_submissions
  WHERE _id = @_id;

-- name: ExportFormSubmissions :many
SELECT _id, form_id, submitted_at, domain, ip_address, form_data, is_spam, email_sent_at
  FROM form_submissions
  WHERE form_id = @form_id
  ORDER BY submitted_at DESC;
