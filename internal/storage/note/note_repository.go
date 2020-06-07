package note

import (
	"github.com/jinzhu/gorm"
	"github.com/pass-wall/passwall-server/model"
)

// Repository ...
type Repository struct {
	db *gorm.DB
}

// NewRepository ...
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// All ...
func (p *Repository) All(schema string) ([]model.Note, error) {
	notes := []model.Note{}
	err := p.db.Table(schema + ".notes").Find(&notes).Error
	return notes, err
}

// FindAll ...
func (p *Repository) FindAll(argsStr map[string]string, argsInt map[string]int, schema string) ([]model.Note, error) {
	notes := []model.Note{}

	query := p.db
	query = query.Table(schema + ".notes")
	query = query.Limit(argsInt["limit"])
	if argsInt["limit"] > 0 {
		// offset can't be declared without a valid limit
		query = query.Offset(argsInt["offset"])
	}

	query = query.Order(argsStr["order"])

	// TODO: This is not working because notes are encrypted
	if argsStr["search"] != "" {
		query = query.Where("note LIKE ?", "%"+argsStr["search"]+"%")
	}

	err := query.Find(&notes).Error
	return notes, err
}

// FindByID ...
func (p *Repository) FindByID(id uint, schema string) (*model.Note, error) {
	note := new(model.Note)
	err := p.db.Table(schema+".notes").Where(`id = ?`, id).First(&note).Error
	return note, err
}

// Save ...
func (p *Repository) Save(note *model.Note, schema string) (*model.Note, error) {
	err := p.db.Table(schema + ".notes").Save(&note).Error
	return note, err
}

// Delete ...
func (p *Repository) Delete(id uint, schema string) error {
	err := p.db.Table(schema + ".notes").Delete(&model.Note{ID: id}).Error
	return err
}

// Migrate ...
func (p *Repository) Migrate(schema string) error {
	return p.db.Table(schema + ".notes").AutoMigrate(&model.Note{}).Error
}
