package preprocessing

import (
	"net/url"
)

type documentCollection struct {
	docList []document
}

type document struct {
	id int64
	docUrl url.URL
	docContent string
}