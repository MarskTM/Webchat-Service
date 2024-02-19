package infrastructure

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	APPPORT    = "APP_PORT"
	DBHOST     = "DB_HOST"
	DBPORT     = "DB_PORT"
	DBUSER     = "DB_USER"
	DBPASSWORD = "DB_PASSWORD"
	DBNAME     = "DB_NAME"

	ROOTPATH = "ROOT_PATH"
	HTTPURL  = "HTTP_URL"
	ENV      = "ENV"
)

var (
	env string

	appPort    string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string

	httpURL    string
	rootPath   string
	staticPath string

	InfoLog *log.Logger
	ErrLog  *log.Logger

	ZapSugar  *zap.SugaredLogger
	ZapLogger *zap.Logger

	db          *gorm.DB
	storagePath string
)

// Hàm sẽ trả về giá trị của biến môi trường
func getStringEnvParameter(envParam string) string {
	if value, ok := os.LookupEnv(envParam); ok {
		return value
	} else {
		return ""
	}
}

func goDotEnvVariable(key string, version int) string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading.env file")
	}
	return os.Getenv(key)
}

func loadEnvParameters() {
	root, _ := os.Getwd()
	env = getStringEnvParameter(ENV)
	appPort = getStringEnvParameter(APPPORT)
	dbPort = getStringEnvParameter(DBPORT)

	InfoLog.Printf("Environment: %s\n", env)
	dbHost = getStringEnvParameter(DBHOST)
	dbUser = getStringEnvParameter(DBUSER)
	dbPassword = getStringEnvParameter(DBPASSWORD)
	dbName = getStringEnvParameter(DBNAME)

	httpURL = getStringEnvParameter(HTTPURL)

	staticPath = root + "/static"
	storagePath = "webchat_storage"
}

func init() {
	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
	ErrLog = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	zapSgLogger, _ := zap.NewDevelopment()
	ZapLogger, _ = zap.NewDevelopment()
	ZapSugar = zapSgLogger.Sugar()

	var initDB bool
	flag.BoolVar(&initDB, "db", false, "allow recreate model database in postgres")
	flag.Parse()

	loadEnvParameters()
	if err := InitDatabase(initDB); err != nil {
		ErrLog.Println("error initialize database: ", err)
	}
}

// ----------------------------------------------------------------
// Declare Geters

func GetDBName() string {
	return dbName
}

// GetDB export db
func GetDB() *gorm.DB {
	return db
}

// GetHTTPURL export http url
func GetHTTPURL() string {
	return httpURL
}