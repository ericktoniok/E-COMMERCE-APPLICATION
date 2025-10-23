package appctx

import (
	"mini-ecommerce/backend/internal/auth"
	"mini-ecommerce/backend/internal/config"
	"gorm.io/gorm"
)

type Context struct {
	DB  *gorm.DB
	JWT *auth.JWTManager
	Cfg *config.Config
}
