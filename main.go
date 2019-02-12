package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-ego/gse"
	"log"
	"os"
	"strings"
	"text/template"
)

var TEMPLATE = `<!DOCTYPE netscape-bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<html>
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <title>Bookmarks</title>
 </head>
 <body>
  <h1>Bookmarks</h1>
  <dl>
   <p> </p>
   {{range $v, $k := .}}
    <dt><a href="{{$k.Url}}" add_date="{{$k.Created}}" {{if $k.Updated}}last_modified="{{$k.Updated}}"{{end}}{{if $k.Tag}} tags="{{$k.Tag}}"{{end}}>{{$k.Title}}</a>{{end}}</dt>
  </dl>
  <p> </p>
 </body>
</html>
`

type BookMark struct {
	Title   string   // 标题
	Url     string   // 地址
	Created string   // 创建时间
	Updated string   // 修改时间
	Tags    []string // 标签
	Tag     string   // 标签
}

//type BookMarks []BookMark

var bookmarks []BookMark
var seg gse.Segmenter

func main() {
	// 加载分词字典
	seg.LoadDict()
	// 分析书签
	analyze("./bookmarks_2019_2_12.html")
	// 生成书签文件
	tmpl, err := template.New("bookmark").Parse(TEMPLATE)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, bookmarks)
	if err != nil {
		panic(err)
	}
}

func analyze(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	// use the goquery document...
	urls := doc.Find("dl > dt > a")
	for i := range urls.Nodes {
		aurl := urls.Eq(i)

		title := aurl.Text()
		url, _ := aurl.Attr("href")
		created, _ := aurl.Attr("add_date")
		updated, _ := aurl.Attr("last_modified")

		var book BookMark
		book.Title = title
		book.Url = url
		book.Created = created
		book.Updated = updated

		segments := seg.Segment([]byte(book.Title))
		for _, s := range segments {
			word := s.Token().Text()
			if s.Token().Pos() == "n" && len([]rune(word)) > 1 {
				book.Tags = append(book.Tags, word)
			}
		}
		book.Tag = strings.Join(book.Tags, ",")
		bookmarks = append(bookmarks, book)
	}
}
