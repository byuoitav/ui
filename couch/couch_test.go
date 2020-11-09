package couch

import (
	"context"
	"testing"

	"github.com/matryer/is"
)

func TestNewFail(t *testing.T) {
	is := is.New(t)

	_, err := New(context.TODO(), "")
	is.Equal("unable to build client: no URL specified", err.Error())
}
