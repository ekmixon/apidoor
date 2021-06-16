package gateway

import (
	"errors"
	"log"
	"strings"
)

type block struct {
	value   string
	isParam bool
}

type URITemplate struct {
	path []block
}

func NewURITemplate(path string) *URITemplate {
	u := &URITemplate{}
	if len(path) == 0 {
		log.Fatal("invalid path")
	}

	slice := strings.Split(path[1:], "/")
	for _, v := range slice {
		isParam := strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}")
		var value string
		if isParam {
			value = v[1 : len(v)-1]
		} else {
			value = v
		}
		u.path = append(u.path, block{
			value,
			isParam,
		})
	}

	return u
}

func (u *URITemplate) TemplateMatch(t URITemplate) (map[string]string, bool) {
	params := make(map[string]string)
	if len(u.path) != len(t.path) {
		return nil, false
	}

	for i := 0; i < len(u.path); i++ {
		if t.path[i].isParam {
			params[t.path[i].value] = u.path[i].value
		} else if u.path[i].value != t.path[i].value {
			return nil, false
		}
	}

	return params, true
}

func (u *URITemplate) JoinPath() string {
	var s []string
	for _, v := range u.path {
		s = append(s, v.value)
	}

	return strings.Join(s, "/")
}

func (u *URITemplate) AllocateParameter(m map[string]string) error {
	for i, block := range u.path {
		if block.isParam {
			if v, ok := m[block.value]; !ok {
				return errors.New("no such parameter")
			} else {
				u.path[i].value = v
			}
		}
	}

	return nil
}
