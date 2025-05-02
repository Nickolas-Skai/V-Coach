package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/cohune-cabbage/di/internal/data"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	db                     *sql.DB
	addr                   *string
	logger                 *slog.Logger
	templateCache          map[string]*template.Template
	questionModel          *data.QuestionModel
	InterviewResponseModel *data.InterviewResponseModel
	signUpModel            *data.SignUpModel
	loginModel             *data.LoginModel
	sessionManager         *data.SessionManager
}

type HomePageData struct {
	Title       string
	Header      string
	Description string
}

func main() {
	err := godotenv.Load(".envrc")
	if err != nil {
		log.Fatalf("Error loading .envrc file")
	}

	addr := flag.String("addr", os.Getenv("ADDRESS"), "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("VCOACH_DB_DSN"), "PostgreSQL DSN")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("database connection pool established")
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		addr:                   addr,
		logger:                 logger,
		templateCache:          templateCache,
		db:                     db,
		signUpModel:            data.NewSignUpModel(db),
		loginModel:             data.NewLoginModel(db),
		questionModel:          &data.QuestionModel{DB: db},
		InterviewResponseModel: data.NewInterviewResponseModel(db),
		sessionManager:         data.NewSessionManager(db, os.Getenv("SESSION_KEY")),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

//set models

type QuestionModel struct {
	DB *sql.DB
}
type SignUpModel struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password_hash"`
	Role      string    `json:"role"`
	Age       *int      `json:"age,omitempty"`
	SchoolID  *int      `json:"school_id,omitempty"`
	CoachID   *int      `json:"coach_id,omitempty"`
}
