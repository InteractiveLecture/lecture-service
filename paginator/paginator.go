package paginator

import (
	"net/url"
	"strconv"
	"strings"
)

type PageRequest struct {
	Number int
	Size   int
	Sorts  []Sort
}

type SortDirection int

const (
	ASC SortDirection = iota
	DESC
)

type Sort struct {
	Name      string
	Direction SortDirection
}

type DepthRequest struct {
	Layer       int
	Ancestors   int
	Descendants int
}

func ParseDepth(u *url.URL) (dr DepthRequest, err error) {
	dr = DepthRequest{}
	dr.Layer, err = strconv.Atoi(parsePart(u, "layer", "0"))
	if err != nil {
		return
	}
	dr.Ancestors, err = strconv.Atoi(parsePart(u, "ancestors", "-1"))
	if err != nil {
		return
	}
	dr.Descendants, err = strconv.Atoi(parsePart(u, "descendants", "-1"))
	return
}

func ParsePages(u *url.URL) (page PageRequest, err error) {
	page = PageRequest{}
	page.Number, err = strconv.Atoi(parsePart(u, "page", "0"))
	if err != nil {
		return
	}
	page.Size, err = strconv.Atoi(parsePart(u, "size", "10"))
	page.Sorts = parseSort(u)
	return
}

func parseSort(u *url.URL) []Sort {
	parts := u.Query()["sort"]
	var result = make([]Sort, 0)
	for _, v := range parts {
		sort := Sort{}
		sort.Name = v
		direction := parsePart(u, v+".dir", "desc")
		if strings.ToLower(direction) == "asc" {
			sort.Direction = ASC
		} else {
			sort.Direction = DESC
		}
		result = append(result, sort)
	}
	return result
}

func parsePart(u *url.URL, part string, defaultValue interface{}) string {
	parts := u.Query()[part]

	if len(parts) == 0 {
		return defaultValue.(string)
	}
	return parts[0]
}
