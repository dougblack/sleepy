package sleepy

import (
	"net/url"
)

type (
	GetNotSupported    struct{}
	PostNotSupported   struct{}
	PutNotSupported    struct{}
	DeleteNotSupported struct{}
)

func (GetNotSupported) Get(values ...url.Values) (int, interface{}) {
	return 405, ""
}

func (PostNotSupported) Post(values ...url.Values) (int, interface{}) {
	return 405, ""
}

func (PutNotSupported) Put(values ...url.Values) (int, interface{}) {
	return 405, ""
}

func (DeleteNotSupported) Delete(values ...url.Values) (int, interface{}) {
	return 405, ""
}
