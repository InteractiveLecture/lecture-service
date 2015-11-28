package datamapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToParameters(t *testing.T) {
	stmt, parameters := prepare("insert into test values(%v)", 1, 2, []string{"hallo", "welt"})
	assert.Equal(t, 4, len(parameters))
	assert.Equal(t, "insert into test values($1,$2,$3,$4)", stmt)
}
