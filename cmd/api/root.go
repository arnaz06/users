package main

import (
	"database/sql"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/arnaz06/users"
	"github.com/arnaz06/users/cmd/logger"
	mysqlRepo "github.com/arnaz06/users/internal/mysql"
	service "github.com/arnaz06/users/user"
)

var (
	contextTimeout time.Duration
	userRepository users.UserRepository
	userService    users.UserService
	secretKey      string
	expiresTime    time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "users",
	Short: "users is a service for managing users",
}

func init() {
	logger.SetupLogs()
	cobra.OnInitialize(initApp)
}

// Execute the main function
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func initApp() {
	/*==== Key ======*/
	secretKey = os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY not set")
	}

	expiry, err := strconv.ParseInt(os.Getenv("TOKEN_EXPIRY_DATE"), 10, 16)
	if err != nil {
		log.Fatalf("TOKEN_EXPIRY_DATE not set %+v", err)
	}
	expiresTime = time.Duration(expiry) * time.Second

	/*==== CONTEXT-TIMEOUT ======*/
	t, err := strconv.ParseInt(os.Getenv("CONTEXT_TIMEOUT_MS"), 10, 16)
	if err != nil {
		log.Fatalf("CONTEXT_TIMEOUT_MS not set %+v", err)
	}
	contextTimeout = time.Duration(t) * time.Millisecond

	/*====MYSQL CONF======*/

	dsnMysql := os.Getenv("MYSQL_URI")
	if dsnMysql == "" {
		log.Fatal("MYSQL_URI not set")
	}

	db, err := sql.Open("mysql", dsnMysql)
	if err != nil {
		log.Fatalf("Can't open MYSQL connection to: %s, got err: %v", dsnMysql, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Can't connect to MYSQL DB got Error: %+v", err)
	}

	maxIdleConn := os.Getenv("MYSQL_MAX_IDLE_CONNECTION")
	if maxIdleConn == "" {
		log.Fatal("MYSQL_MAX_IDLE_CONNECTION not set")
	}

	mysqlMaxIdleCon, err := strconv.Atoi(maxIdleConn)
	if err != nil {
		log.Fatal("invalid MYSQL_MAX_IDLE_CONNECTION")
	}
	db.SetMaxIdleConns(mysqlMaxIdleCon)

	maxOpenConn := os.Getenv("MYSQL_MAX_OPEN_CONNECTION")
	if maxIdleConn == "" {
		log.Fatal("MYSQL_MAX_OPEN_CONNECTION not set")
	}

	mysqlMaxOpenCon, err := strconv.Atoi(maxOpenConn)
	if err != nil {
		log.Fatal("invalid MYSQL_MAX_OPEN_CONNECTION")
	}
	db.SetMaxOpenConns(mysqlMaxOpenCon)

	connLifeTime := os.Getenv("MYSQL_CONNECTION_LIFETIME_M")
	if maxIdleConn == "" {
		log.Fatal("MYSQL_CONNECTION_LIFETIME_M not set")
	}

	mysqlMaxConnLifetime, err := strconv.Atoi(connLifeTime)
	if err != nil {
		log.Fatal("invalid MYSQL_CONNECTION_LIFETIME_M")
	}
	db.SetConnMaxLifetime(time.Minute * time.Duration(mysqlMaxConnLifetime))

	userRepository = mysqlRepo.NewUserRepository(db)
	userService = service.NewUserService(userRepository)
}
