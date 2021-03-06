package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "net/http"
  "code.google.com/p/go-html-transform/html/transform"
  "strings"
)

func main() {
  c := make(chan string)
  go consume(os.Args[1], c)
  for i:=0; i<150; i++ {
    url := <- c
    fmt.Printf("%v\n", url)
  }
}

func consume(url string, c chan string) {
  urls, err := crawl(url)
  if err==nil {
    for i := 0; i<len(urls); i++ {
      c <- urls[i]
      go consume(urls[i], c)
    }
  }
}

func crawl(url string) (urls []string, err error) {
  content, err := getContent(url)
  if err==nil {
    urls, err := parseLinks(content)
    return urls, err
  }
  return urls, err
}

func parseLinks(content string) (urls []string, err error) {
  doc, err := transform.NewDoc(content)
  if err == nil {
    selector := transform.NewSelectorQuery("a")
    nodes := selector.Apply(doc)
    for i := 0; i<len(nodes); i++ {
      for j := 0; j<len(nodes[i].Attr); j++ {
        if nodes[i].Attr[j].Name == "href" && strings.HasPrefix(nodes[i].Attr[j].Value, "http") {
          urls = append(urls, nodes[i].Attr[j].Value)
        }
      }
    }
  }
  return urls, err
}

func getContent(url string) (content string, err error) {
  resp, err := http.Get(url)
  if err == nil {
    defer resp.Body.Close()
    var body []byte
    body, err = ioutil.ReadAll(resp.Body)
    if err == nil {
      content = string(body)
    }
  }
  return content, err
}