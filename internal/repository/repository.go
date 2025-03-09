package repository

import (
	"github.com/adorufus/imgupper/pkg/database"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Repositories struct {
	User   UserRepository
	Health HealthRepository
	Cr2    Cr2Repository
}

func NewRepositories(db *database.Database, s3Client *s3.Client) *Repositories {
	return &Repositories{
		User:   NewUserRepository(db),
		Health: NewHealthRepository(db),
		Cr2:    NewCr2Repository(db, s3Client),
	}
}
