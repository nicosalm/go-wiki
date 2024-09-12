package main

import (
    "html/template"
    "io/fs"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strings"

    "github.com/yuin/goldmark"
)

type Page struct {
    Title   string
    Body    template.HTML // rendered HTML
    RawBody string        // raw markdown content
}

func (p *Page) save() error {
    filename := "data/" + p.Title + ".md"
    return os.WriteFile(filename, []byte(p.RawBody), 0600)
}

func loadPage(title string) (*Page, error) {
    // load .md files from the data directory
    filename := "data/" + title + ".md"
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    // markdown -> HTML using goldmark
    var htmlBody strings.Builder
    if err := goldmark.Convert(body, &htmlBody); err != nil {
        return nil, err
    }

    return &Page{
        Title:   title,
        Body:    template.HTML(htmlBody.String()), // rendered HTML
        RawBody: string(body),                     // raw markdown
    }, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    if title == "Directory" {
        serveDirectory(w)
        return
    }

    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body") // this will contain the raw Markdown
    p := &Page{Title: title, RawBody: body}

    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func serveDirectory(w http.ResponseWriter) {
    var files []string

    // list all .md files in the data directory
    err := filepath.Walk("data/", func(path string, info fs.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
            pageName := strings.TrimSuffix(info.Name(), ".md")
            files = append(files, pageName)
        }
        return nil
    })

    if err != nil {
        http.Error(w, "Error listing directory", http.StatusInternalServerError)
        return
    }

    var body strings.Builder
    body.WriteString(`<h1>Welcome to the wiki pages directory!</h1>`)
    body.WriteString(`<ul>`)
    for _, file := range files {
        body.WriteString(`<li><a href="/view/` + file + `">` + file + `</a></li>`)
    }
    body.WriteString(`</ul>`)
    body.WriteString(`Return to <a href="/view/FrontPage">FrontPage</a>`)

    w.Write([]byte(body.String()))
}

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

// redirect root to "/view/FrontPage"
func rootHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

func main() {
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))

    log.Fatal(http.ListenAndServe(":8080", nil))
}

