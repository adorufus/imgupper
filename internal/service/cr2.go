package service

import (
	"context"
	"mime/multipart"

	"github.com/adorufus/imgupper/internal/model"
)

type Cr2Service interface {
	ObjectUpload(ctx context.Context, req model.CR2UploadRequest, object multipart.File, handler *multipart.FileHeader) (model.CR2UploadResponse, error)
	ObjectFetchById(ctx context.Context, id int64) (model.CR2UploadResponse, error)
	ObjectFetchByUserId(ctx context.Context) ([]model.CR2UploadResponse, error)
}

type cr2Service struct {
	deps Deps
}

// ObjectFetchByUserId implements Cr2Service.
func (s *cr2Service) ObjectFetchByUserId(ctx context.Context) ([]model.CR2UploadResponse, error) {
	return s.deps.Repos.Cr2.GetByUserID(ctx)
}

// ObjectFetchById implements Cr2Service.
func (s *cr2Service) ObjectFetchById(ctx context.Context, id int64) (model.CR2UploadResponse, error) {
	return s.deps.Repos.Cr2.GetByID(ctx, id)
}

// ObjectUpload implements Cr2Service.
func (s *cr2Service) ObjectUpload(ctx context.Context, req model.CR2UploadRequest, object multipart.File, handler *multipart.FileHeader) (model.CR2UploadResponse, error) {
	return s.deps.Repos.Cr2.Create(ctx, req, object, handler)
}

func NewCr2Srvice(deps Deps) Cr2Service {
	return &cr2Service{
		deps: deps,
	}
}
