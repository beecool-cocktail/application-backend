package service

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func newHTTP(configure *Configure) (*gin.Engine, error) {
	if configure.HTTP == nil {
		return nil, errors.New("http configure is not initialed")
	}

	r := gin.Default()
	r.MaxMultipartMemory = 32 << 20

	return r, nil
}
