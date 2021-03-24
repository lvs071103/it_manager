package zfrtest1

type Retriever struct {
	Contents string
}

func (r Retriever) Get(url string) string {
	return r.Contents
}
