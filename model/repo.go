package model

import (
	"fmt"
	"time"

	"github.com/Machiel/slugify"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

// Repos is simply a collection of repo structs.
type Repos []*Repo

// Repo represents a repo model definition.
type Repo struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Org       *Org      `json:"org,omitempty"`
	OrgID     int       `json:"org_id" sql:"index"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	Public    bool      `json:"private" sql:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      Tags      `json:"tags,omitempty"`
}

// BeforeSave invokes required actions before persisting.
func (u *Repo) BeforeSave(db *gorm.DB) (err error) {
	if u.Slug == "" {
		for i := 0; true; i++ {
			if i == 0 {
				u.Slug = slugify.Slugify(u.Name)
			} else {
				u.Slug = slugify.Slugify(
					fmt.Sprintf("%s-%d", u.Name, i),
				)
			}

			notFound := db.Where(
				"slug = ?",
				u.Slug,
			).Not(
				"id",
				u.ID,
			).First(
				&Repo{},
			).RecordNotFound()

			if notFound {
				break
			}
		}
	}

	return nil
}

// AfterDelete invokes required actions after deletion.
func (u *Repo) AfterDelete(tx *gorm.DB) error {
	for _, tag := range u.Tags {
		if err := tx.Delete(tag).Error; err != nil {
			return err
		}
	}

	return nil
}

// Validate does some validation to be able to store the record.
func (u *Repo) Validate(db *gorm.DB) {
	if !govalidator.StringLength(u.Name, "1", "255") {
		db.AddError(fmt.Errorf("Name should be longer than 1 and shorter than 255"))
	}

	if u.Name != "" {
		notFound := db.Where(
			"name = ?",
			u.Name,
		).Not(
			"id",
			u.ID,
		).First(
			&Repo{},
		).RecordNotFound()

		if !notFound {
			db.AddError(fmt.Errorf("Name is already present"))
		}
	}
}