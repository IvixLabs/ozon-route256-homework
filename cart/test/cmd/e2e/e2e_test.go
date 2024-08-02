package e2e

import (
	"github.com/stretchr/testify/suite"
	"route256/cart/test/suite/e2e"
	"testing"
)

func TestRest(t *testing.T) {
	suite.Run(t, &e2e.RestSuite{})
}
