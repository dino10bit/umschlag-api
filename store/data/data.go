package data

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"github.com/o1egl/gormrus"
	"github.com/qor/validations"
	"github.com/umschlag/umschlag-api/config"
	"github.com/umschlag/umschlag-api/model"
	"github.com/umschlag/umschlag-api/store"

	// Register MySQL driver for GORM
	_ "github.com/jinzhu/gorm/dialects/mysql"

	// Register Postgres driver for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	dialects []string
)

// Store is a basic struct to represent the database handle.
type data struct {
	*gorm.DB
}

// Test creates an in-memory database connection.
func Test() store.Store {
	var (
		driver = "sqlite3"
		config = ":memory:"
	)

	if os.Getenv("DATABASE_DRIVER") != "" && os.Getenv("DATABASE_DRIVER") != "" {
		driver = os.Getenv("DATABASE_DRIVER")
		config = os.Getenv("DATABASE_CONFIG")
	}

	db, err := gorm.Open(driver, config)

	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("database connection failed")
	}

	return &data{
		setupDatabase(driver, db),
	}
}

// New initializes a new database connection.
func New(driver string, config string) store.Store {
	db, err := gorm.Open(driver, config)

	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Database connection failed")
	}

	return &data{
		setupDatabase(driver, db),
	}
}

// From takes an existing database connection.
func From(driver string, handle *sql.DB) store.Store {
	db, err := gorm.Open(driver, handle)

	if err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Database connection failed")
	}

	return &data{
		setupDatabase(driver, db),
	}
}

// Load initializes the database connection.
func Load() store.Store {
	driver := config.Database.Driver
	connect := ""

	if invalidDriver(driver) {
		logrus.Fatal("Unknown database driver selected")
	}

	switch driver {
	case "mysql":
		connect = fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.Database.Username,
			config.Database.Password,
			config.Database.Host,
			config.Database.Name,
		)
	case "postgres":
		connect = fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			config.Database.Username,
			config.Database.Password,
			config.Database.Host,
			config.Database.Name,
		)
	case "sqlite3":
		connect = config.Database.Name
	}

	logrus.Infof("Using database driver %s", driver)
	logrus.Infof("Using database config %s", connect)

	return New(
		driver,
		connect,
	)
}

func invalidDriver(driver string) bool {
	for _, dialect := range dialects {
		logrus.Debugf("Checking if %s matches %s database driver", dialect, driver)

		if dialect == driver {
			return false
		}
	}

	return true
}

func setupDatabase(driver string, db *gorm.DB) *gorm.DB {
	if config.Debug {
		db.LogMode(true)
		db.SetLogger(gormrus.New())
	}

	if err := prepareDatabase(driver, db); err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Database preparation failed")
	}

	if err := pingDatabase(driver, db); err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Database ping attempts failed")
	}

	if err := migrateDatabase(driver, db); err != nil {
		logrus.Errorln(err)
		logrus.Fatalln("Database migration failed")
	}

	return db
}

func prepareDatabase(driver string, db *gorm.DB) error {
	if driver == "mysql" {
		db.DB().SetMaxIdleConns(0)
	}

	validations.RegisterCallbacks(
		db,
	)

	return nil
}

func pingDatabase(driver string, db *gorm.DB) error {
	for i := 0; i < 30; i++ {
		err := db.DB().Ping()

		if err == nil {
			return nil
		}

		logrus.Infof("Database ping failed, retry in 1s")
		time.Sleep(time.Second)
	}

	return nil
}

func migrateDatabase(driver string, db *gorm.DB) error {
	db.AutoMigrate(
		&model.Org{},
		&model.Registry{},
		&model.Repo{},
		&model.Tag{},
		&model.Team{},
		&model.User{},
	)

	db.Model(
		&model.Org{},
	).AddUniqueIndex(
		"uix_orgs_registry_id_slug",
		"registry_id",
		"slug",
	)

	db.Model(
		&model.Org{},
	).AddUniqueIndex(
		"uix_orgs_registry_id_name",
		"registry_id",
		"name",
	)

	db.Model(
		&model.Repo{},
	).AddUniqueIndex(
		"uix_repos_org_id_slug",
		"org_id",
		"slug",
	)

	db.Model(
		&model.Repo{},
	).AddUniqueIndex(
		"uix_repos_org_id_name",
		"org_id",
		"name",
	)

	db.Model(
		&model.Tag{},
	).AddUniqueIndex(
		"uix_tags_repo_id_slug",
		"repo_id",
		"slug",
	)

	db.Model(
		&model.Tag{},
	).AddUniqueIndex(
		"uix_tags_repo_id_name",
		"repo_id",
		"name",
	)

	if db.First(&model.User{}).RecordNotFound() {
		record := &model.User{
			Username: "admin",
			Password: "admin",
			Email:    "admin@example.com",
			Active:   true,
			Admin:    true,
		}

		err := db.Create(
			record,
		).Error

		if err != nil {
			return fmt.Errorf(
				"Failed to create initial user. %s",
				err.Error(),
			)
		}
	}

	return nil
}

func init() {
	dialects = append(dialects, "mysql")
	dialects = append(dialects, "postgres")
}
