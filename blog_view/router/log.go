package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Log(res http.ResponseWriter, req *http.Request) {
	var logMain []byte

	if doctype, err := ioutil.ReadFile("/home/juunini/blog_view/views/doctype.html"); err != nil {
		stderr(err)
		return
	} else {
		logMain = append(logMain, doctype...)
	}

	Req, err := http.Get("http://127.0.0.1:9959/blog?query={post(_id:\"" + req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:] + "\"){title,time,views,tags,contents}}")
	if err != nil {
		stderr(err)
		return
	}
	body, err := ioutil.ReadAll(Req.Body)
	if err != nil {
		stderr(err)
		return
	}
	var logContents struct {
		Post []post `json:"post"`
	}
	if err := json.Unmarshal(body, &logContents); err != nil {
		stderr(err)
		return
	}

	logMain = append(logMain, []byte(`<meta name="keyword" content="`+strings.Join(logContents.Post[0].Tags, ",")+`">`)...)

	if head, err := ioutil.ReadFile("/home/juunini/blog_view/views/head.html"); err != nil {
		stderr(err)
		return
	} else {
		logMain = append(logMain, head...)
	}

	var tag string
	for _, t := range logContents.Post[0].Tags {
		tag += fmt.Sprintf("<a class=\"tag\" href=\"/tags/" + t + "\">" + t + "</a>")
	}
	logMain = append(logMain, []byte(fmt.Sprintf(`
<div class="content-title">
	<h2 class="title">%s</h2>
	<span class="time">%s</span>
	<span class="views">%d</span>
	<p class="tags">%s</p>
</div>
`, logContents.Post[0].Title, logContents.Post[0].Time.String()[:19], logContents.Post[0].Views, tag))...)

	if main, err := ioutil.ReadFile("/home/juunini/blog_view/views/main.html"); err != nil {
		stderr(err)
		return
	} else {
		logMain = append(logMain, main...)
	}

	logMain = append(logMain, []byte(logContents.Post[0].Contents)...)

	if err := req.Body.Close(); err != nil {
		stderr(err)
		return
	}

	logMain = append(logMain, []byte(`
<p style="text-align:center;margin-top:60px;">
<img src="/img/CCL.png" alt="" style="width:200px;">
<br>
이 저작물은 크리에이티브 커먼즈 저작자표시 4.0 국제 라이선스에 따라 이용할 수 있습니다.
</p>
`)...)

	if footer, err := ioutil.ReadFile("/home/juunini/blog_view/views/footer.html"); err != nil {
		stderr(err)
		return
	} else {
		logMain = append(logMain, footer...)
	}

	fmt.Println(time.Now().String()[:19], "["+req.RemoteAddr[:strings.LastIndex(req.RemoteAddr, ":")]+"]", logContents.Post[0].Title, logContents.Post[0].Views)

	if _, err := res.Write(logMain); err != nil {
		stderr(err)
		return
	}
}
