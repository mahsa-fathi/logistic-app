package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"logistic-app/internal/common/configs"
)

type MockPostgres struct {
	Postgres
}

func NewMockPostgresDB() (*MockPostgres, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
		configs.DBTestAddress, configs.DBTestUser, configs.DBTestPassword, configs.DBTestName, configs.DBTestPort,
	)
	db, e := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if e != nil {
		return nil, e
	}
	pdb := &MockPostgres{Postgres{db: db}}
	e = pdb.initializeDB()
	return pdb, e
}

func (p *MockPostgres) Close() {
	p.db.Exec(`DROP TABLE orders`)
	p.db.Exec(`DROP TABLE customers`)
	p.db.Exec(`DROP TABLE providers`)
	p.db.Exec(`DROP TABLE periodic_tasks`)
	if sql, e := p.db.DB(); e == nil {
		_ = sql.Close()
	}
	return
}
