package sleepy

type (
    GetNotSupported struct {}
    PostNotSupported struct {}
    PutNotSupported struct {}
    DeleteNotSupported struct {}
)

func (GetNotSupported) Get(map[string][]string) string {
	return "Nope."
}

func (PostNotSupported) Post(map[string][]string) string {
	return "Nope."
}

func (PutNotSupported) Put(map[string][]string) string {
	return "Nope."
}

func (DeleteNotSupported) Delete(map[string][]string) string {
	return "Nope."
}
