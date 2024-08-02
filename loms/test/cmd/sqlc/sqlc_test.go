package sqlc

import (
	"github.com/stretchr/testify/suite"
	"route256/loms/test/suite/sqlc"
	"testing"
)

func TestSQLC(t *testing.T) {
	suite.Run(t, &sqlc.RepositorySuite{})
}
