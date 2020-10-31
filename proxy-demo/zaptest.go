package main

import (
	"fmt"
	"proxy-demo/mock"
	"proxy-demo/real"
)

type Retriever interface {
	Get(url string) string
}

func download(r Retriever) string {
	return r.Get("http://www.baidu.com")
}

type Poster interface {
	Post(url string, form map[string]string) string
}

func post(poster Poster) {
	poster.Post("http://www.baidu.com",
		map[string]string{
		"name":"leosanqing",
		"course":"golang",
		})
}


type RetrieverPoster interface {
	Retriever
	Poster
}



func main() {
	var r Retriever
	r = mock.Retriever{"sdfasd "}
	r = real.Retriever{}
	fmt.Println(download(r))

}
