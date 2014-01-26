package sleepy

import (
	"net/url"
)

type BaseResource struct{}

func (BaseResource) Get(values url.Values) (int, interface{}) {
	return 405, ""
}

func (BaseResource) Post(values url.Values) (int, interface{}) {
	return 405, ""
}

func (BaseResource) Put(values url.Values) (int, interface{}) {
	return 405, ""
}

func (BaseResource) Delete(values url.Values) (int, interface{}) {
	return 405, ""
}
