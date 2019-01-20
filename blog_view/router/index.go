package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type post struct {
	ID       string    `json:"_id"`
	Title    string    `json:"title"`
	Time     time.Time `json:"time"`
	Views    int       `json:"views"`
	Tags     []string  `json:"tags"`
	Contents string    `json:"contents"`
}

func Index(res http.ResponseWriter, req *http.Request) {
	if len(req.URL.Path) > 1 {
		if req.URL.Path == "/sitemap.xml" {
			FileReader, err := ioutil.ReadFile("/home/juunini/blog_view/views/sitemap.xml")
			if err != nil {
				stderr(err)
				return
			}
			res.Header().Set("Content-Type", "application/xml")
			if _, err := res.Write(FileReader); err != nil {
				stderr(err)
				return
			}
			return
		} else if req.URL.Path == "/robots.txt" {
			FileReader, err := ioutil.ReadFile("/home/juunini/blog_view/views/robots.txt")
			if err != nil {
				stderr(err)
				return
			}
			res.Header().Set("Content-Type", "text/plain")
			if _, err := res.Write(FileReader); err != nil {
				stderr(err)
				return
			}
			return
		}
	}

	var index []byte

	if doctype, err := ioutil.ReadFile("/home/juunini/blog_view/views/doctype.html"); err != nil {
		stderr(err)
		return
	} else {
		index = append(index, doctype...)
	}
	if head, err := ioutil.ReadFile("/home/juunini/blog_view/views/head.html"); err != nil {
		stderr(err)
		return
	} else {
		index = append(index, head...)
	}
	if main, err := ioutil.ReadFile("/home/juunini/blog_view/views/main.html"); err != nil {
		stderr(err)
		return
	} else {
		index = append(index, main...)
	}
	List, err := http.Get("http://127.0.0.1:9959/blog?query={post{_id,title,time,views,tags}}")
	if err != nil {
		stderr(err)
		return
	}
	body, err := ioutil.ReadAll(List.Body)
	if err != nil {
		stderr(err)
		return
	}
	var list struct {
		Post []post `json:"post"`
	}
	if err := json.Unmarshal(body, &list); err != nil {
		stderr(err)
		return
	}
	for _, row := range list.Post {
		var tag string
		for _, t := range row.Tags {
			tag += fmt.Sprintf("<span class=\"tag\">" + t + "</span>")
		}
		index = append(index, []byte(fmt.Sprintf(`
<a href="/log/%s" class="list">
	<span class="title">%s</span>
	<span class="time">%s</span>
	<span class="views">%d</span>
	<span class="tags">%s</span>
</a>
`, row.ID[strings.LastIndex(row.ID, "(")+2:len(row.ID)-2], row.Title, row.Time.String()[:19], row.Views, tag))...)
	}
	if err := List.Body.Close(); err != nil {
		stderr(err)
		return
	}
	if footer, err := ioutil.ReadFile("/home/juunini/blog_view/views/footer.html"); err != nil {
		stderr(err)
		return
	} else {
		index = append(index, footer...)
	}

	if _, err := res.Write(index); err != nil {
		stderr(err)
		return
	}
}

func stderr(err error) {
	if _, err := fmt.Fprintln(os.Stderr, time.Now().String()[:19], err); err != nil {
		return
	}
}