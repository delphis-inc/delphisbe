package datastore

import (
	"sync"
)

type testBackend struct {
	db              Datastore
	discussionMutex sync.Mutex
}

//func MakeDatastore(ctx context.Context, testData TestData) (Datastore, func() error, error) {
//	url := "postgres://chatham_local@localhost:5432/"
//
//	dbName, err := createTestDatabase(ctx, url+"?sslmode=disable")
//	if err != nil {
//		return nil, nil, errors.Wrap(err, "failed to create test database")
//	}
//
//	url = url + dbName + "?sslmode=disable"
//
//	db, err := connectTestDatabase(ctx, url)
//	if err != nil {
//		return nil, nil, errors.Wrap(err, "failed to create testing datastore")
//	}
//
//	// Create test tables with test data
//	closeFunc, err := db.CreateTestTables(ctx, testData)
//	if err != nil {
//		return nil, nil, errors.Wrap(err, "failed to create test tables")
//	}
//
//	return db, closeFunc, nil
//}
//
//func connectTestDatabase(ctx context.Context, url string) (Datastore, error) {
//	// Initialize gorm
//	testGorm, err := gorm.Open("postgres", url)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to open testing database - gorm")
//	}
//
//	// Initialize sql
//	testDB, err := sql.Open("postgres", url)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to open testing database - SQL")
//	}
//
//	// Create db object
//	return &delphisDB{
//		sql:       testGorm,
//		pg:        testDB,
//		prepStmts: &dbPrepStmts{},
//	}, nil
//}
//
//func createTestDatabase(ctx context.Context, url string) (string, error) {
//	db, err := gorm.Open("postgres", url)
//	if err != nil {
//		return "", err
//	}
//
//	rand.Seed(time.Now().UnixNano())
//	name := "tests_" + strconv.Itoa(rand.Int())
//
//	logrus.Infof("Database Name: %v\n", name)
//
//	db = db.Exec(fmt.Sprintf(`create database %s;`, name))
//	if db.Error != nil {
//		return "", err
//	}
//
//	return name, nil
//}
