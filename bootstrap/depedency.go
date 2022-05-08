package bootstrap

import (
	"github.com/jmoiron/sqlx"
)

type Dependency struct {
	cfg Config
	db  *sqlx.DB
	// redis           *redis.Client
	// dispatcher      *channel.DispatcherHandler
	// smtp            *smtp.Client
	// encryptor       crypto.Repository
	// minioClient     *minio.Client
	// azureSharedCred *azblob.SharedKeyCredential
}

func NewDependency() *Dependency {
	dep := new(Dependency)
	return dep
}

func (dep *Dependency) GetConfig() Config {
	return dep.cfg
}

func (dep *Dependency) SetConfig(cfg Config) {
	dep.cfg = cfg
}

func (dep *Dependency) GetDB() *sqlx.DB {
	if dep.db == nil {
		dep.db = NewMysqlDB(dep.GetConfig().Database.Write)
	}
	return dep.db
}
