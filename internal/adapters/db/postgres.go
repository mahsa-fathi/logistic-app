package db

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"logistic-app/internal/app/domain"
	"logistic-app/internal/common/configs"
	"logistic-app/internal/common/errors"
	"strings"
	"time"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgresDB() (*Postgres, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
		configs.DBAddress, configs.DBUser, configs.DBPassword, configs.DBName, configs.DBPort,
	)
	db, e := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if e != nil {
		return nil, e
	}
	pdb := &Postgres{db: db}
	e = pdb.initializeDB()
	return pdb, e
}

func (p *Postgres) initializeDB() error {
	e := p.db.AutoMigrate(&domain.Customer{})
	if e != nil {
		return e
	}
	e = p.db.AutoMigrate(&domain.Provider{})
	if e != nil {
		return e
	}
	e = p.db.AutoMigrate(&domain.Order{})
	if e != nil {
		return e
	}

	if !p.db.Migrator().HasIndex(&domain.Order{}, "idx_ongoing_status") {
		statusList := strings.Join(domain.GetOngoingOrderStatus(), "', '")
		p.db.Exec(fmt.Sprintf(`
      CREATE INDEX CONCURRENTLY idx_ongoing_status
      ON orders(status)
      WHERE status IN ('%s');
    `, statusList))
	}

	e = p.db.AutoMigrate(&domain.PeriodicTask{})
	return e
}

func (p *Postgres) Close() {
	if sql, e := p.db.DB(); e == nil {
		_ = sql.Close()
	}
	return
}

func (p *Postgres) Ready() bool {
	if sql, e := p.db.DB(); e == nil {
		e = sql.Ping()
		if e != nil {
			return false
		}
		return true
	}
	return false
}

func (p *Postgres) GetAllProviders(ctx context.Context) ([]*domain.Provider, *errors.AppError) {
	var providers []*domain.Provider
	result := p.db.WithContext(ctx).Find(&providers)
	return providers, errors.ConvertGormErrors(result.Error)
}

func (p *Postgres) GetProvider(ctx context.Context, providerID uint) (*domain.Provider, *errors.AppError) {
	var provider *domain.Provider
	result := p.db.WithContext(ctx).First(&provider, providerID)
	return provider, errors.ConvertGormErrors(result.Error)
}

func (p *Postgres) CreateProvider(ctx context.Context, name, url *string) (*domain.Provider, *errors.AppError) {
	provider := &domain.Provider{
		Name: *name,
		Url:  *url,
	}
	result := p.db.WithContext(ctx).Create(&provider)
	return provider, errors.ConvertGormErrors(result.Error)
}

func (p *Postgres) GetCustomer(ctx context.Context, userID uint) (*domain.Customer, *errors.AppError) {
	var customer *domain.Customer
	result := p.db.WithContext(ctx).First(&customer, userID)
	return customer, errors.ConvertGormErrors(result.Error)
}

func (p *Postgres) CreateCustomer(ctx context.Context, name, phone, addr, postalCode *string) (*domain.Customer, *errors.AppError) {
	customer := &domain.Customer{
		PhoneNumber: *phone,
		Name:        name,
		Address:     *addr,
		PostalCode:  *postalCode,
	}
	result := p.db.WithContext(ctx).Create(&customer)
	return customer, errors.ConvertGormErrors(result.Error)
}

func (p *Postgres) CreateOrUpdatePeriodicTask(ctx context.Context, name string, interval int, failed bool, e *string) (*domain.PeriodicTask, *errors.AppError) {
	var task *domain.PeriodicTask
	t := time.Now()
	result := p.db.WithContext(ctx).
		Where(domain.PeriodicTask{JobName: name}).
		Assign(domain.PeriodicTask{
			IntervalInMinute: uint(interval),
			LastRunTime:      &t,
			Failed:           failed,
			Error:            e,
		}).
		FirstOrCreate(&task)
	return task, errors.ConvertGormErrors(result.Error)
}

func (p *Postgres) GetOrCreatePeriodicTask(ctx context.Context, name string, interval int) (*domain.PeriodicTask, *errors.AppError) {
	var task *domain.PeriodicTask
	result := p.db.WithContext(ctx).
		Where(domain.PeriodicTask{JobName: name}).
		Assign(domain.PeriodicTask{IntervalInMinute: uint(interval)}).
		FirstOrCreate(&task)
	return task, errors.ConvertGormErrors(result.Error)
}
