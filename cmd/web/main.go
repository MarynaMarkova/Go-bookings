package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/MarynaMarkova/Go-bookings/internal/config"
	"github.com/MarynaMarkova/Go-bookings/internal/driver"
	"github.com/MarynaMarkova/Go-bookings/internal/handlers"
	"github.com/MarynaMarkova/Go-bookings/internal/helpers"
	"github.com/MarynaMarkova/Go-bookings/internal/models"
	"github.com/MarynaMarkova/Go-bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	from := "me@here.com"
	auth := smtp.PlainAuth("", from, "", "localhost")
	err = smtp.SendMail("localhost:1025", auth, from, []string{"you@there.com"}, []byte("Hello, world"))
	if err != nil {
		log.Print(err)
	}

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	
	srv := &http.Server {
	Addr: 		portNumber,
	Handler: 	routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
	log.Fatal(err)
	}
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

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL(`host=localhost port=5432 dbname=bookings user=postgres password=free23lancer`)
	// database.yml password
	if err !=nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}

//To run the program: go run ./cmd/web/.