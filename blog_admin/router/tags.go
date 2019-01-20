package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Tags(res http.ResponseWriter, req *http.Request) {
	var tagsMain []byte

	if doctype, err := ioutil.ReadFile("/home/juunini/blog_view/views/doctype.html"); err != nil {
		stderr(err)
		return
	} else {
		tagsMain = append(tagsMain, doctype...)
	}
	if head, err := ioutil.ReadFile("/home/juunini/blog_view/views/head.html"); err != nil {
		stderr(err)
		return
	} else {
		tagsMain = append(tagsMain, head...)
	}
	if main, err := ioutil.ReadFile("/home/juunini/blog_view/views/main.html"); err != nil {
		stderr(err)
		return
	} else {
		tagsMain = append(tagsMain, main...)
	}

	tagsMain = append(tagsMain, []byte(`<div class="tags">`)...)

	Req, err := http.Get("http://127.0.0.1:9959/blog?query={tags}")
	if err != nil {
		stderr(err)
		return
	}
	body, err := ioutil.ReadAll(Req.Body)
	if err != nil {
		stderr(err)
		return
	}
	var tagList map[string][]string
	if err := json.Unmarshal(body, &tagList); err != nil {
		stderr(err)
		return
	}
	for _, tag := range tagList["tags"] {
		tagsMain = append(tagsMain, []byte("<a class=\"tag\" href=\"/tags/"+tag+"\">"+tag+"</a>")...)
	}
	if err := Req.Body.Close(); err != nil {
		stderr(err)
		return
	}

	tagsMain = append(tagsMain, []byte(`</div>`)...)

	if footer, err := ioutil.ReadFile("/home/juunini/blog_view/views/footer.html"); err != nil {
		stderr(err)
		return
	} else {
		tagsMain = append(tagsMain, footer...)
	}

	if _, err := res.Write(tagsMain); err != nil {
		stderr(err)
		return
	}
}

func TagsList(res http.ResponseWriter, req *http.Request) {
	var tagsMain []byte

	if doctype, err := ioutil.ReadFile("/home/juunini/blog_view/views/doctype.html"); err != nil {
		stderr(err)
		return
	} else {
		tagsMain = append(tagsMain, doctype...)
	}
	if head, err := ioutil.ReadFile("/home/juunini/blog_view/views/head.html"); err != nil {
		stderr(err)
		return
	} else {
		tagsMain = append(tagsMain, head...)
	}
	if main, err := ioutil.ReadFile("/home/juunini/blog_view/views/main.html"); err != nil {
		stderr(err)
		return
	} else {
		tagsMain = append(tagsMain, main...)
	}

	tagsMain = append(tagsMain, []byte(`<div class="tags">`)...)

	Req, err := http.Get("http://127.0.0.1:9959/blog?query={tags}")
	if err != nil {
		stderr(err)
		return
	}
	body, err := ioutil.ReadAll(Req.Body)
	if err != nil {
		stderr(err)
		return
	}
	var tagList map[string][]string
	if err := json.Unmarshal(body, &tagList); err != nil {
		stderr(err)
		return
	}
	for _, tag := range tagList["tags"] {
		tagsMain = append(tagsMain, []byte("<a class=\"tag\" href=\"/tags/"+tag+"\">"+tag+"</a>")...)
	}
	if err := Req.Body.Close(); err != nil {
		stderr(err)
		return
	}

	tagsMain = append(tagsMain, []byte(`</div>`)...)

	Req, err = http.Get("http://127.0.0.1:9959/blog?query={post(tag:\"" + req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:] + "\"){_id,title,time,views,tags}}")
	if err != nil {
		stderr(err)
		return
	}
	body, err = ioutil.ReadAll(Req.Body)
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
		tagsMain = append(tagsMain, []byte(fmt.Sprintf(`
<a href="/log/%s" class="list">
	<span class="title">%s</span>
	<span class="time">%s</span>
	<span class="views">%d</span>
	<span class="tags">%s</span>
</a>
`, row.ID[strings.LastIndex(row.ID, "(")+2:len(row.ID)-2], row.Title, row.Time.String()[:19], row.Views, tag))...)
	}
	if err := Req.Body.Close(); err != nil {
		stderr(err)
		return
	}

	if footer, err := ioutil.ReadFile("/home/juunini/blog_view/views/footer.html"); err != nil {
		stderr(err)
		return
	} else {
		tagsMain = append(tagsMain, footer...)
	}

	if _, err := res.Write(tagsMain); err != nil {
		stderr(err)
		return
	}
}
