package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"joyful.go/go-concurrency/final/cmd/data"
)

const webPort = "8080"

func main() {

	//connect to the database
	db := initDB()
	db.Ping()
	//create sessions
	session := initSession()
	//create channels
	//create waitgroup
	wg := sync.WaitGroup{}
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	//setup the application config
	app := Config{
		Session:  session,
		DB:       db,
		Wait:     &wg,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Models:   data.New(db),
	}
	//set up mail
	app.Mailer = app.createMail()
	go app.listenForMail()

	//listen for web connection
	app.serve()

	// listen for signals
	app.listenForShutdown()

}

func (app *Config) createMail() Mail {

	errorChan := make(chan error)
	mailerChan := make(chan Message, 100)

	mailerDone := make(chan bool)

	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encrption:   "none",
		FromAddress: "info@mycompany.com",
		FromName:    "info",
		ErrChan:     errorChan,
		Wait:        app.Wait,
		MailerChan:  mailerChan,
		DoneChan:    mailerDone,
	}
	return m

}

func initSession() *scs.SessionManager {

	gob.Register(data.User{})
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

func initRedis() *redis.Pool {
	redisURI := "127.0.0.1:6379"
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			// return redis.Dial("tcp", os.Getenv("REDIS"))
			return redis.Dial("tcp", redisURI)
		},
	}
	return redisPool
}

func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	// perform any clean up
	app.InfoLog.Println("cleaning up...")
	// block until waitgroup is empty
	app.Wait.Wait()
	app.Mailer.DoneChan <- true
	app.InfoLog.Println("closing channels and shutting down application...")
	close(app.Mailer.MailerChan)
	close(app.Mailer.ErrChan)
	close(app.Mailer.DoneChan)
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting web server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to database")
	}
	return conn
}

func connectToDB() *sql.DB {
	counts := 0
	// dsn := os.Getenv("DSN")
	dsn := "host=localhost port=5432 user=postgres password=password dbname=concurrency sslmode=disable timezone=UTC connect_timeout=5"
	fmt.Printf("The environment DSN is set to %s\n", dsn)

	for {
		connection, err := opendDB(dsn)
		if err != nil {
			log.Println("postres not yet ready...\n")
			counts++
		} else {
			log.Print("Connected to database!\n")
			return connection
		}

		if counts > 10 {
			return nil
		}

		log.Print("Backing off for 1 second\n")
		time.Sleep(1 * time.Second)
	}

}

func opendDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}
