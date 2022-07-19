package main

import (
	"encoding/gob"
	"fmt"
	"github.com/marcionps/bookings/internal/driver"
	"github.com/marcionps/bookings/internal/helpers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/marcionps/bookings/internal/config"
	"github.com/marcionps/bookings/internal/handlers"
	"github.com/marcionps/bookings/internal/models"
	"github.com/marcionps/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	db.SQL.Close()

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("connecting to the database")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user='db_dev' password='guinness'")
	if err != nil {
		log.Fatal("cannot connect to database!")
	}
	log.Println("connected to the database successfully")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
