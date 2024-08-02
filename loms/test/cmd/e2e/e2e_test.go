package e2e

import (
	"github.com/stretchr/testify/suite"
	"route256/loms/test/suite/e2e"
	"testing"
)

func TestGRPC(t *testing.T) {
	suite.Run(t, &e2e.GRPCSuite{})
}
