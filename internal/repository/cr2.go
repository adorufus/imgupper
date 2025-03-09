package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/adorufus/imgupper/internal/model"
	"github.com/adorufus/imgupper/pkg/database"
	"github.com/adorufus/imgupper/pkg/middleware"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// Cr2Repository defines the file repository interface
type Cr2Repository interface {
	Create(ctx context.Context, file model.CR2UploadRequest, object multipart.File, handler *multipart.FileHeader) (model.CR2UploadResponse, error)
	GetByID(ctx context.Context, id int64) (model.CR2UploadResponse, error)
	GetByUserID(ctx context.Context) ([]model.CR2UploadResponse, error)
	// GetAll(ctx context.Context) ([]model.File, error)
	// Update(ctx context.Context, file model.File) (model.File, error)
	// Delete(ctx context.Context, id int64) error
}

// cr2Repository implements FileRepository
type cr2Repository struct {
	db       *database.Database
	s3Client *s3.Client
}

// NewFileRepository creates a new FileRepository
func NewCr2Repository(db *database.Database, s3Client *s3.Client) Cr2Repository {
	return &cr2Repository{
		db:       db,
		s3Client: s3Client,
	}
}

// Create creates a new file record
func (r *cr2Repository) Create(ctx context.Context, file model.CR2UploadRequest, object multipart.File, handler *multipart.FileHeader) (model.CR2UploadResponse, error) {

	print("user id", file.UserID)
	// First, check if user exists
	var userExists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", file.UserID).Scan(&userExists)
	if err != nil {
		return model.CR2UploadResponse{}, err
	}

	if !userExists {
		return model.CR2UploadResponse{}, errors.New("User Not Found")
	}

	uid := file.UserID

	// Generate unique filename
	filename := fmt.Sprintf("u/%v/uploads/%s-%d%s",
		uid,
		uuid.New().String(),
		time.Now().Unix(),
		getFileExtension(handler.Filename),
	)

	client := r.s3Client
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("ember-imgupper"),
		Key:         aws.String(filename),
		Body:        object,
		ContentType: aws.String(handler.Header.Get("Content-Type")),
		ACL:         "public-read",
	})

	fmt.Print("hit this")

	if err != nil {
		return model.CR2UploadResponse{}, err
	}

	bucketUrl := "https://cdn.imgupper.web.id/" + filename

	query := `
		INSERT INTO files (user_id, filename, filesize, mime_type, bucket_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, user_id, filename, filesize, mime_type, bucket_url, created_at, updated_at
	`

	var createdFile model.CR2UploadResponse

	err = r.db.QueryRowContext(
		ctx,
		query,
		file.UserID,
		handler.Filename,
		handler.Size,
		handler.Header.Get("Content-Type"),
		bucketUrl,
	).Scan(
		&createdFile.ID,
		&createdFile.UserID,
		&createdFile.Filename,
		&createdFile.Filesize,
		&createdFile.MimeType,
		&createdFile.BucketURL,
		&createdFile.CreatedAt,
		&createdFile.UpdatedAt,
	)

	if err != nil {
		return model.CR2UploadResponse{}, fmt.Errorf("failed to create file record: %w", err)
	}

	return createdFile, nil
}

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

// GetByID gets a file by ID
func (r *cr2Repository) GetByID(ctx context.Context, id int64) (model.CR2UploadResponse, error) {
	query := `
		SELECT id, user_id, filename, filesize, mime_type, bucket_url, created_at, updated_at
		FROM files
		WHERE id = $1
	`

	var file model.CR2UploadResponse
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&file.ID,
		&file.UserID,
		&file.Filename,
		&file.Filesize,
		&file.MimeType,
		&file.BucketURL,
		&file.CreatedAt,
		&file.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.CR2UploadResponse{}, fmt.Errorf("file not found: %w", err)
		}
		return model.CR2UploadResponse{}, fmt.Errorf("failed to get file: %w", err)
	}

	return file, nil
}

// // GetByUserID gets all files for a specific user
func (r *cr2Repository) GetByUserID(ctx context.Context) ([]model.CR2UploadResponse, error) {
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return []model.CR2UploadResponse{}, err
	}

	uid := user.UserID

	query := `
		SELECT id, user_id, filename, filesize, mime_type, bucket_url, created_at, updated_at
		FROM files
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to query files for user: %w", err)
	}
	defer rows.Close()

	var files []model.CR2UploadResponse
	for rows.Next() {
		var file model.CR2UploadResponse
		if err := rows.Scan(
			&file.ID,
			&file.UserID,
			&file.Filename,
			&file.Filesize,
			&file.MimeType,
			&file.BucketURL,
			&file.CreatedAt,
			&file.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating file rows: %w", err)
	}

	return files, nil
}

// // GetAll gets all files
// func (r *fileRepository) GetAll(ctx context.Context) ([]model.File, error) {
// 	query := `
// 		SELECT id, user_id, filename, filesize, mime_type, bucket_url, created_at, updated_at
// 		FROM files
// 		ORDER BY created_at DESC
// 	`

// 	rows, err := r.db.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query files: %w", err)
// 	}
// 	defer rows.Close()

// 	var files []model.File
// 	for rows.Next() {
// 		var file model.File
// 		if err := rows.Scan(
// 			&file.ID,
// 			&file.UserID,
// 			&file.Filename,
// 			&file.Filesize,
// 			&file.MimeType,
// 			&file.BucketURL,
// 			&file.CreatedAt,
// 			&file.UpdatedAt,
// 		); err != nil {
// 			return nil, fmt.Errorf("failed to scan file: %w", err)
// 		}
// 		files = append(files, file)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating file rows: %w", err)
// 	}

// 	return files, nil
// }

// Update updates a file record
// func (r *fileRepository) Update(ctx context.Context, file model.File) (model.File, error) {
// 	query := `
// 		UPDATE files
// 		SET filename = $1, filesize = $2, mime_type = $3, bucket_url = $4, updated_at = NOW()
// 		WHERE id = $5
// 		RETURNING id, user_id, filename, filesize, mime_type, bucket_url, created_at, updated_at
// 	`

// 	var updatedFile model.File
// 	err := r.db.QueryRowContext(
// 		ctx,
// 		query,
// 		file.Filename,
// 		file.Filesize,
// 		file.MimeType,
// 		file.BucketURL,
// 		file.ID,
// 	).Scan(
// 		&updatedFile.ID,
// 		&updatedFile.UserID,
// 		&updatedFile.Filename,
// 		&updatedFile.Filesize,
// 		&updatedFile.MimeType,
// 		&updatedFile.BucketURL,
// 		&updatedFile.CreatedAt,
// 		&updatedFile.UpdatedAt,
// 	)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return model.File{}, fmt.Errorf("file not found: %w", err)
// 		}
// 		return model.File{}, fmt.Errorf("failed to update file: %w", err)
// 	}

// 	return updatedFile, nil
// }

// // Delete deletes a file record
// func (r *fileRepository) Delete(ctx context.Context, id int64) error {
// 	query := `
// 		DELETE FROM files
// 		WHERE id = $1
// 	`

// 	result, err := r.db.ExecContext(ctx, query, id)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete file: %w", err)
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get rows affected: %w", err)
// 	}

// 	if rowsAffected == 0 {
// 		return fmt.Errorf("file not found")
// 	}

// 	return nil
// }
