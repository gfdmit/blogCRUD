package repository

import (
	"blog/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(conf config.Config) (*Repository, error) {
	url := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=disable",
		conf.Host,
		conf.Username,
		conf.DB,
		conf.Password,
	)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %v", err)
	}

	return &Repository{db: db}, nil
}

func (repo Repository) Migrate(models ...interface{}) error {
	return repo.db.AutoMigrate(models...)
}

func (repo Repository) Create() {

}

func (repo Repository) Read() {

}

func (repo Repository) Update() {

}

func (repo Repository) Delete() {

}
