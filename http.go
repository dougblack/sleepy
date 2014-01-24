package sleepy

type GetNotSupported struct{}

func (GetNotSupported) Get(map[string][]string) string {
	return "Nope."
}

type PostNotSupported struct{}

func (PostNotSupported) Post(map[string][]string) string {
	return "Nope."
}

type PutNotSupported struct{}

func (PutNotSupported) Put(map[string][]string) string {
	return "Nope."
}

type DeleteNotSupported struct{}

func (DeleteNotSupported) Delete(map[string][]string) string {
	return "Nope."
}

