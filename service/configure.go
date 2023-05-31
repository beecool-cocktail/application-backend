package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Configure struct {
	Version string
	Logger  *loggerConfig
	DB      *mysqlConfig
	HTTP    *httpConfig
	Redis   *redisConfig
	Elastic *elasticConfig
	Others  *OtherConfig
}

type loggerConfig struct {
	Path     string
	FileName string
	Level    string
}

type mysqlConfig struct {
	MainDB mysqlDBParam
}

type mysqlDBParam struct {
	DriverName         string
	User               string
	Password           string
	Address            string
	DBName             string
	SetConnMaxIdleTime int
	SetMaxIdleConns    int //sec
	SetMaxOpenConns    int //sec
}

type httpConfig struct {
	Address         string
	Port            string
	IsTLS           bool
	CertificateFile string
	KeyFile         string
}

type redisConfig struct {
	Network            string
	Address            string
	Password           string
	DB                 int
	DialTimeoutSecond  int
	ReadTimeoutSecond  int
	WriteTimeoutSecond int
	PoolSize           int
}

type elasticConfig struct {
	Enable bool
	Urls   []string
}

type OtherConfig struct {
	GoogleOAuth2 *GoogleOAuth2
	File         *File
}

type GoogleOAuth2 struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

type File struct {
	Image *Image
}

type Image struct {
	PathInServer string
	PathInURL    string
}

func newConfigure(fileName string) (*Configure, error) {

	var configure Configure
	configFilePath := getFilePosition(fileName)

	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	//unmarshal to struct
	if err := json.Unmarshal(b, &configure); err != nil {
		return nil, err
	}

	return &configure, nil
}

func getFilePosition(fileName string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
	}

	var buf bytes.Buffer
	if fileName[0] == '/' {
		buf.WriteString(fileName)
	} else {
		buf.WriteString(dir)
		buf.WriteString("/")
		buf.WriteString(fileName)
	}

	return buf.String()
}
