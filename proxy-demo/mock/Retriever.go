package mock

type Retriever struct {
	Contents string
}

func (r Retriever) Post(url string, form map[string]string) string {
	r.Contents = form["contents"]
	panic("implement me")
}

func (r Retriever) Get(url string) string {
	return r.Contents
}
