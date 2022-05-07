package bootstrap

import "database/sql"


type Dependency struct {
	cfg             Config
	db              *sql.DB
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

func (dep *Dependency) GetDB() *sql.DB {
	if dep.db == nil {
		dep.db = NewMysqlDB(dep.GetConfig().Database.Write);
	}
	return dep.db
}