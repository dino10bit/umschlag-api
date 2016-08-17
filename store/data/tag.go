package data

import (
	"github.com/jinzhu/gorm"
	"github.com/umschlag/umschlag-api/model"
)

// GetTags retrieves all available tags from the database.
func (db *data) GetTags() (*model.Tags, error) {
	records := &model.Tags{}

	err := db.Order(
		"name ASC",
	).Find(
		records,
	).Error

	return records, err
}

// CreateTag creates a new tag.
func (db *data) CreateTag(record *model.Tag) error {
	return db.Create(
		record,
	).Error
}

// UpdateTag updates a tag.
func (db *data) UpdateTag(record *model.Tag) error {
	return db.Save(
		record,
	).Error
}

// DeleteTag deletes a tag.
func (db *data) DeleteTag(record *model.Tag) error {
	return db.Delete(
		record,
	).Error
}

// GetTag retrieves a specific tag from the database.
func (db *data) GetTag(id string) (*model.Tag, *gorm.DB) {
	record := &model.Tag{}

	res := db.Where(
		"id = ?",
		id,
	).Or(
		"slug = ?",
		id,
	).Model(
		record,
	).First(
		record,
	)

	return record, res
}
