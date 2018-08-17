package mysqltype

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

var (
	user         = os.Getenv("GORM_EXT_DB_USER")
	password     = os.Getenv("GORM_EXT_DB_PASSWORD")
	host         = os.Getenv("GORM_EXT_DB_HOST_NAME")
	port         = os.Getenv("GORM_EXT_DB_PORT")
	databaseName = os.Getenv("GORM_EXT_DB_NAME")
)

func init() {
	time.Local = time.UTC
	if err := initDB(); err != nil {
		panic(err.Error())
	}
}
func initDB() error {
	db, err := openDB("")
	if err != nil {
		return err
	}
	if databaseName == "" {
		databaseName = "gormext_094598c0_a11b_11e8_b689_7c7a9175123c"
	}
	if err := db.Exec(fmt.Sprintf("drop database if exists %s", databaseName)).Error; err != nil {
		return err
	}
	if err := db.Exec(fmt.Sprintf("create database if not exists %s", databaseName)).Error; err != nil {
		return err
	}

	db, err = openDB(databaseName)
	if err != nil {
		return err
	}
	DB = db

	return nil
}
func openDB(databaseName string) (*gorm.DB, error) {

	if user == "" {
		user = "gormexttest"
	}

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "3306"
	}

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=UTC", user, password, host, port, databaseName))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

func runTest(m *testing.M) int {
	defer testDown()
	return m.Run()
}

func testDown() {
	if err := DB.Exec("drop database " + databaseName).Error; err != nil {
		panic(err.Error())
	}
	DB.Close()
}
