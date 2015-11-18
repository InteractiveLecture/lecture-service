package paginator

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginator(t *testing.T) {
	u, _ := url.Parse("http://example.com?page=1&sort=hugo&hugo.dir=asc")

	p, err := ParsePages(u)
	assert.Nil(t, err)
	assert.Equal(t, 1, p.Number)
	assert.Equal(t, 10, p.Size)
	assert.Equal(t, 1, len(p.Sorts))
	sort := p.Sorts[0]
	assert.Equal(t, ASC, sort.Direction)
	assert.Equal(t, "hugo", sort.Name)

	u, _ = url.Parse("http://example.com?sort=hugo&hugo.dir=asc&size=30&sort=bla")

	p, err = ParsePages(u)

	assert.Nil(t, err)
	sorts := p.Sorts
	assert.Equal(t, 2, len(sorts))
	assert.Equal(t, DESC, sorts[1].Direction)
	assert.Equal(t, 0, p.Number)

}
