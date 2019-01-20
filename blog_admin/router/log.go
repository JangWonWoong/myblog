package router

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

	logMain = append(logMain, []byte(fmt.Sprintf(`
<div class="content-title">
	<h2 class="title"><input type="text" id="title" value="%s" style="width:100%%;padding:10px 20px;color:#d7dae0;background:#21252b;border:3px solid #282c34;border-radius:10px;outline:none;"></h2>
	<span class="time">%s</span>
	<span class="views">%d</span>
	<p class="tags"><input type="text" id="tags" class="tag" value="%s" style="width:100%%;padding:10px 20px;color:#d7dae0;background:#21252b;border:3px solid #282c34;border-radius:10px;outline:none;"></p>
	<span style="display:inline-block;margin-right:10px;padding:10px 14px;background:crimson;color:white;font-size:16px;border-radius:10px;cursor:pointer;" onclick="del('%s')">삭제</span>
	<span style="display:inline-block;padding:10px 14px;background:blue;color:white;font-size:16px;border-radius:10px;cursor:pointer;" onclick="mod('%s')">수정</span>
</div>
`, logContents.Post[0].Title, logContents.Post[0].Time.String()[:19], logContents.Post[0].Views, strings.Join(logContents.Post[0].Tags, ","), req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:], req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]))...)

	if main, err := ioutil.ReadFile("/home/juunini/blog_view/views/main.html"); err != nil {
		stderr(err)
		return
	} else {
		logMain = append(logMain, main...)
	}

	logMain = append(logMain, []byte(`<textarea id="contents" style="width:100%;height:600px;padding:20px;color:#d7dae0;background:#21252b;border:3px solid #282c34;border-radius:10px;outline:none;">`+logContents.Post[0].Contents+`</textarea>`)...)

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
<script>
	function del(id) {
		let xhr = new XMLHttpRequest();
		xhr.open("GET", "/delete/" + id);
		xhr.send();
		xhr.onload = function() {
			if (this.response === "done") {
				alert("삭제되었습니다.");
				location.href = "/";
			} else {
				alert("실패했습니다.");
				location.href = "/";
			}
		};
	}

	function mod(id) {
		let xhr = new XMLHttpRequest();
		xhr.open("PUT", "/update/" + id);

		let form = new FormData();
		form.append("title", document.getElementById("title").value);
		form.append("tags", document.getElementById("tags").value);
		form.append("contents", document.getElementById("contents").value.replace(/\n/g, ""));

		xhr.send(form);
		xhr.onload = function() {
			if (this.response === "done") {
				alert("수정되었습니다.");
				location.href = "/log/" + id;
			} else {
				alert("실패했습니다.");
				location.href = "/log/" + id;
			}
		};
	}
</script>
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

func Delete(res http.ResponseWriter, req *http.Request) {
	_id := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]

	if _, err := http.PostForm(`http://127.0.0.1:9959/blog`, url.Values{"query": {fmt.Sprintf(`mutation{delete(_id:"%s")}`, _id)}}); err != nil {
		stderr(err)
		return
	}

	deleteSitemap(_id)

	if _, err := res.Write([]byte("done")); err != nil {
		stderr(err)
		return
	}
}

func Update(res http.ResponseWriter, req *http.Request) {
	_id := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
	title := req.FormValue("title")
	tags := req.FormValue("tags")
	contents := req.FormValue("contents")

	tagsList := strings.Split(tags, ",")
	tags = ""
	for _, t := range tagsList {
		tags += fmt.Sprintf(`"%s",`, t)
	}

	if _, err := http.PostForm("http://127.0.0.1:9959/blog", url.Values{"query": {fmt.Sprintf(`mutation{update(_id:"%s",title:"%s",tags:[%s],contents:"%s"){_id}}`, _id, title, tags[:len(tags)-1], contents)}}); err != nil {
		stderr(err)
		return
	}

	updateSitemap(_id)

	if _, err := res.Write([]byte("done")); err != nil {
		stderr(err)
		return
	}
}

func Create(res http.ResponseWriter, req *http.Request) {
	if strings.EqualFold(req.Method, "GET") {
		var logMain []byte

		if doctype, err := ioutil.ReadFile("/home/juunini/blog_view/views/doctype.html"); err != nil {
			stderr(err)
			return
		} else {
			logMain = append(logMain, doctype...)
		}

		if head, err := ioutil.ReadFile("/home/juunini/blog_view/views/head.html"); err != nil {
			stderr(err)
			return
		} else {
			logMain = append(logMain, head...)
		}

		logMain = append(logMain, []byte(`
<div class="content-title">
	<h2 class="title"><input type="text" id="title" placeholder="제목" style="width:100%;padding:10px 20px;color:#d7dae0;background:#21252b;border:3px solid #282c34;border-radius:10px;outline:none;"></h2>
	<p class="tags"><input type="text" id="tags" class="tag" placeholder="태그" style="width:100%;padding:10px 20px;color:#d7dae0;background:#21252b;border:3px solid #282c34;border-radius:10px;outline:none;"></p>
	<span style="display:inline-block;margin-right:10px;padding:10px 14px;background:green;color:white;font-size:16px;border-radius:10px;cursor:pointer;" onclick="create()">생성</span>
</div>`)...)

		if main, err := ioutil.ReadFile("/home/juunini/blog_view/views/main.html"); err != nil {
			stderr(err)
			return
		} else {
			logMain = append(logMain, main...)
		}

		logMain = append(logMain, []byte(`
<textarea id="contents" style="width:100%;height:600px;padding:20px;color:#d7dae0;background:#21252b;border:3px solid #282c34;border-radius:10px;outline:none;"></textarea>
<script>
	function create() {
		let xhr = new XMLHttpRequest();
		xhr.open("POST", "/create");

		let form = new FormData();
		form.append("title", document.getElementById("title").value);
		form.append("tags", document.getElementById("tags").value);
		form.append("contents", document.getElementById("contents").value.replace(/\n/g, ""));

		xhr.send(form)
		xhr.onload = function() {
			if (this.response === "done") {
				alert("생성에 성공하였습니다.");
				location.href = "/";
			} else {
				alert("생성에 실패하였습니다.");
				location.href = "/";
			}
		}
	}
</script>
`)...)

		if footer, err := ioutil.ReadFile("/home/juunini/blog_view/views/footer.html"); err != nil {
			stderr(err)
			return
		} else {
			logMain = append(logMain, footer...)
		}

		if _, err := res.Write(logMain); err != nil {
			stderr(err)
			return
		}
		return
	} else if strings.EqualFold(req.Method, "POST") {
		title := req.FormValue("title")
		tags := req.FormValue("tags")
		contents := req.FormValue("contents")

		var tag string
		for _, t := range strings.Split(tags, ",") {
			tag += fmt.Sprintf(`"%s",`, t)
		}

		if resp, err := http.PostForm("http://127.0.0.1:9959/blog", url.Values{"query": {fmt.Sprintf(`mutation{create(title:"%s",tags:[%s],contents:"%s"){_id}}`, title, tag[:len(tag)-1], contents)}}); err != nil {
			stderr(err)
			return
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				stderr(err)
				return
			}

			var Body map[string]map[string]string
			if err := json.Unmarshal(body, &Body); err != nil {
				stderr(err)
				return
			}

			result := strings.Replace(Body["create"]["_id"], "ObjectIdHex(\"", "", -1)
			createSitemap(result[:len(result)-2])

			if err := resp.Body.Close(); err != nil {
				stderr(err)
				return
			}
		}

		if _, err := res.Write([]byte("done")); err != nil {
			stderr(err)
			return
		}
	}
}

type urlset struct {
	Xmlns  string       `xml:"xmlns,attr"`
	Url    []sitemapUrl `xml:"url"`
}

type sitemapUrl struct {
	Loc        string `xml:"loc"`
	Lastmod    string `xml:"lastmod"`
	Changefreq string `xml:"changefreq"`
}

func createSitemap(_id string) {
	FileReader, err := ioutil.ReadFile("/home/juunini/blog_view/views/sitemap.xml")
	if err != nil {
		stderr(err)
		return
	}

	var Sitemap urlset
	if err := xml.Unmarshal(FileReader, &Sitemap); err != nil {
		stderr(err)
		return
	}

	Sitemap.Url = append(Sitemap.Url, sitemapUrl{"https://juunini.xyz/log/" + _id, time.Now().String()[:10], "monthly"})

	result, err := xml.Marshal(Sitemap)
	if err != nil {
		stderr(err)
		return
	}

	result = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), result...)

	if err := ioutil.WriteFile("/home/juunini/blog_view/views/sitemap.xml", result, 0644); err != nil {
		stderr(err)
		return
	}

	return
}

func updateSitemap(_id string) {
	FileReader, err := ioutil.ReadFile("/home/juunini/blog_view/views/sitemap.xml")
	if err != nil {
		stderr(err)
		return
	}

	var Sitemap urlset
	if err := xml.Unmarshal(FileReader, &Sitemap); err != nil {
		stderr(err)
		return
	}

	for i := range Sitemap.Url {
		if strings.Contains(Sitemap.Url[i].Loc, _id) {
			Sitemap.Url[i].Lastmod = time.Now().String()[:10]
			break
		}
	}

	result, err := xml.Marshal(Sitemap)
	if err != nil {
		stderr(err)
		return
	}

	result = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), result...)

	if err := ioutil.WriteFile("/home/juunini/blog_view/views/sitemap.xml", result, 0644); err != nil {
		stderr(err)
		return
	}

	return
}

func deleteSitemap(_id string) {
	FileReader, err := ioutil.ReadFile("/home/juunini/blog_view/views/sitemap.xml")
	if err != nil {
		stderr(err)
		return
	}

	var Sitemap urlset
	if err := xml.Unmarshal(FileReader, &Sitemap); err != nil {
		stderr(err)
		return
	}

	for i := range Sitemap.Url {
		if strings.Contains(Sitemap.Url[i].Loc, _id) {
			Sitemap.Url = append(Sitemap.Url[:i], Sitemap.Url[i+1:]...)
			break
		}
	}

	result, err := xml.Marshal(Sitemap)
	if err != nil {
		stderr(err)
		return
	}

	result = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), result...)

	if err := ioutil.WriteFile("/home/juunini/blog_view/views/sitemap.xml", result, 0644); err != nil {
		stderr(err)
		return
	}

	return
}
