package gormx

import (
	"ginAdmin/internal/app/config"
	"ginAdmin/internal/app/model/gormx/entity"
	"ginAdmin/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

// Config 配置参数
type Config struct {
	Debug        bool
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TablePrefix  string
}

// NewDB 创建DB实例
func NewDB(c *Config) (*gorm.DB, func(), error) {
	var gormDB *gorm.DB
	var err error
	switch c.DBType {
	case "mysql":
		gormDB, err = gorm.Open(mysql.Open(c.DSN), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   c.TablePrefix, // 表名前缀，`User`表为`t_users`
				SingularTable: true,          // 使用单数表名，启用该选项后，`User` 表将是`user`
			},
		})
	case "postgres":
		gormDB, err = gorm.Open(postgres.Open(c.DSN), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   c.TablePrefix, // 表名前缀，`User`表为`t_users`
				SingularTable: true,          // 使用单数表名，启用该选项后，`User` 表将是`user`
			},
		})
	case "sqlite3":
		gormDB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   c.TablePrefix, // 表名前缀，`User`表为`t_users`
				SingularTable: true,          // 使用单数表名，启用该选项后，`User` 表将是`user`
			},
		})
	}

	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, err
	}

	if c.Debug {
		gormDB = gormDB.Debug()
	}

	cleanFunc := func() {
		err := sqlDB.Close()
		if err != nil {
			logger.Errorf("Gorm db close error: %s", err.Error())
		}
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, cleanFunc, err
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	return gormDB, cleanFunc, nil
}

// AutoMigrate 自动映射数据表
func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate(
		new(entity.Demo),
		new(entity.Menu),
		new(entity.Role),
		new(entity.UserRole),
		new(entity.User),
		new(entity.RouterResource),
	)
}
