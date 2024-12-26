package api

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type HandlersTestSuite struct {
	suite.Suite
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
