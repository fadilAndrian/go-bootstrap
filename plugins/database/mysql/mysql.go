package mysql

import (
	"fmt"

	"github.com/konnco/go-bootstrap/internal/core"
)

type Plugin struct{}

func init() {
	core.Register(&Plugin{})
}

func (p *Plugin) Metadata() core.PluginMetadata {
	return core.PluginMetadata{
		Name:        "mysql",
		Category:    core.CategoryDatabase,
		Version:     "1.0.0",
		Description: "MySQL database adapter using github.com/go-sql-driver/mysql",
	}
}

func (p *Plugin) Category() core.Category {
	return core.CategoryDatabase
}

func (p *Plugin) Dependencies() []string {
	return []string{
		"github.com/go-sql-driver/mysql@latest",
	}
}

func (p *Plugin) Validate(ctx *core.GenerateContext) error {
	return nil
}

func (p *Plugin) Generate(ctx *core.GenerateContext) ([]core.GeneratedFile, error) {
	content := []byte(fmt.Sprintf(`package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"%s/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	sqlDB  *sql.DB
	logger *slog.Logger
}

func Connect(ctx context.Context, cfg config.DatabaseConfig, logger *slog.Logger) (*DB, error) {
	db, err := sql.Open("mysql", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql database: %%w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if !cfg.SkipPing {
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := db.PingContext(pingCtx); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to ping mysql database: %%w", err)
		}
	}

	logger.Info("connected to mysql database")
	return &DB{sqlDB: db, logger: logger}, nil
}

func (db *DB) SQL() *sql.DB {
	return db.sqlDB
}

func (db *DB) Ping(ctx context.Context) error {
	if db.sqlDB == nil {
		return fmt.Errorf("database is nil")
	}
	return db.sqlDB.PingContext(ctx)
}

func (db *DB) Close() error {
	if db.sqlDB != nil {
		db.logger.Info("closing mysql database")
		return db.sqlDB.Close()
	}
	return nil
}
`, ctx.ModulePath))

	return []core.GeneratedFile{
		{
			Path:    "internal/database/mysql.go",
			Content: content,
		},
	}, nil
}
