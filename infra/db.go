package infra

import (
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
)

// DB ReadWrite/ReadOnly オブジェクト
type DB struct {
	// Master ReadWrite Node
	Master *gorm.DB
	// Slave ReadOnly Node
	Slave *gorm.DB
}

var (
	db *DB
)

func getEnv(key, defVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defVal
}

func dataSourceKey(key, val string) string {
	return fmt.Sprintf("%s=%s ", key, val)
}

func dataSource(host, user, dbName, sslMode, passWord, port, schemaName string) string {
	dns := strings.Builder{}
	dns.WriteString(host)
	dns.WriteString(user)
	dns.WriteString(dbName)
	dns.WriteString(sslMode)
	dns.WriteString(passWord)
	dns.WriteString(port)
	if schemaName != "" {
		schemaName = dataSourceKey("search_path", schemaName)
		dns.WriteString(schemaName)
	}
	return dns.String()
}

func dbOpen(dialect, dns string) *gorm.DB {
	db, err := gorm.Open("postgres", dns)
	if err != nil {
		panic(err)
	}
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}
	return db
}

// NewDBInit DB接続
func NewDBInit() *DB {
	if db != nil {
		return db
	}

	var user, passWord, masterHost, slaveHost, port, dbName, schemaName, replica string
	var master, slave *gorm.DB

	user = dataSourceKey("user", getEnv("DB_USER", "postgres"))
	passWord = dataSourceKey("password", getEnv("DB_PASSWORD", "postgres"))
	masterHost = dataSourceKey("host", getEnv("DB_MASTER_NAME", "host.docker.internal"))
	port = dataSourceKey("port", getEnv("DB_PORT", "30543"))
	dbName = dataSourceKey("dbname", getEnv("DB_DBNAME", "postgres"))
	schemaName = getEnv("DB_SCHEMA", "stockpile")
	sslMode := dataSourceKey("sslmode", "disable")
	master = dbOpen("postgres", dataSource(masterHost, user, dbName, sslMode, passWord, port, schemaName))

	// レプリケーションされたDBの場合
	replica = getEnv("REPLICA", "N")
	if replica == "Y" {
		slaveHost = dataSourceKey("host", getEnv("DB_SLAVE_NAME", "host.docker.internal"))
		slave = dbOpen("postgres", dataSource(slaveHost, user, dbName, sslMode, passWord, port, schemaName))
	}
	sqlLog := getEnv("DB_LOG", "Y")
	if sqlLog == "Y" {
		if master != nil {
			master.LogMode(true)
		}
		if slave != nil {
			slave.LogMode(true)
		}
	}

	db = &DB{
		Master: master,
		Slave:  slave,
	}
	return db
}
