package datastore

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type testBackend struct {
	db              Datastore
	discussionMutex sync.Mutex
}

type TestData struct {
	Discussions []model.Discussion
}

func MakeDatastore(ctx context.Context, testData TestData) (Datastore, func() error, error) {
	url := "postgres://chatham_local@localhost:5432/"

	dbName, err := createTestDatabase(ctx, url+"?sslmode=disable")
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create test database")
	}

	testDB, err := gorm.Open("postgres", url+dbName+"?sslmode=disable")
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open testing database")
	}

	// Create db object
	db := &db{sql: testDB}

	// Create test tables with test data
	if err := createTestTables(ctx, testData, db); err != nil {
		return nil, nil, errors.Wrap(err, "failed to create test tables")
	}

	return db, func() error { return testDB.Close() }, nil
}

func createTestDatabase(ctx context.Context, url string) (string, error) {
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	name := "tests_" + strconv.Itoa(rand.Int())

	logrus.Infof("Database Name: %v\n", name)

	db = db.Exec(fmt.Sprintf(`create database %s;`, name))
	if db.Error != nil {
		return "", err
	}

	return name, nil
}

func createTestTables(ctx context.Context, data TestData, d *db) error {
	if err := writeDiscussions(ctx, data.Discussions, d); err != nil {
		return errors.Wrap(err, "failed to write discussions")
	}

	return nil
}

func writeDiscussions(ctx context.Context, testDiscussions []model.Discussion, d *db) error {
	// Create table with schema based on the Discussion model
	discussions := []model.Discussion{}
	if db := d.sql.AutoMigrate(&discussions); db.Error != nil {
		return db.Error
	}

	logrus.Infof("Do we have this table? %v\n", d.sql.HasTable("discussions"))

	// Iterate over test data to create test records
	for _, discussion := range testDiscussions {
		logrus.Infof("In here for Disc: %+v\n", discussion)
		if _, err := d.UpsertDiscussion(ctx, discussion); err != nil {
			return err
		}

		time.Sleep(3 * time.Second)
		if _, err := d.UpsertDiscussion(ctx, discussion); err != nil {
			return err
		}
	}

	return nil
}
