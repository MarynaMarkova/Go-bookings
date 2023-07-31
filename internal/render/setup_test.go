package render

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/MarynaMarkova/Go-bookings/internal/config"
	"github.com/MarynaMarkova/Go-bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

var functions = template.FuncMap{
	"humanDate": render.HumanDate,
	"formatDate": render.FormatDate,
	"iterate": render.Iterate,
	"add": render.Add,
}

func TestMain(m *testing.M){

	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})


	// change this to true when in production
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog 

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	app.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int){

}

func (tw *myWriter) Write(b []byte)(int, error){
	length := len(b)
	return length, nil
}