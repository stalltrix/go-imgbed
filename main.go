package main

import (
    "encoding/json"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "sync"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"strings"
	"embed"
    "html/template"
)

var (
    user string
    pass string
    dataDir string
    imgMap map[string]struct{}
    imgLock sync.RWMutex
	fs http.Handler 
	allowedExt = map[string]bool{
        ".jpg":  true,
        ".jpeg": true,
        ".png":  true,
        ".webp": true,
    }
)

//go:embed tpl/*
var templateFS embed.FS

var tmpl = template.Must(template.ParseFS(templateFS, "tpl/*.html"))
var data =struct{
	Name string
}{}
type Img_file struct {
    State string `json:"state"`
	OK bool `json:"ok"`
}

func checkAuth(r *http.Request) bool {
    u, p, ok := r.BasicAuth()
	if !ok {
		return false
	}
	if len(u)<2 || len(u)>255 {
		return false
	}
	if len(p)<8 || len(p)>255 {
		return false
	}
    return u == user && p == pass
}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if !checkAuth(r) {
        w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    if r.Method != "POST" {
		tmpl.ExecuteTemplate(w, "upload.html", data)
        return
    }
	w.Header().Set("Content-Type", "application/json")
    file, header, err := r.FormFile("file")
    if err != nil {
		files := Img_file{
			State: err.Error(),
			OK: false,
		}
		json.NewEncoder(w).Encode(files)
		return
    }
    defer file.Close()
	ext:=strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExt[ext] {
        files := Img_file{
			State: "only jpg, png, webp allowed",
			OK: false,
		}
		json.NewEncoder(w).Encode(files)
		return
    }
	f, err := os.CreateTemp(dataDir, "upload-*")
	if err != nil {
		files := Img_file{
			State: err.Error(),
			OK: false,
		}
		json.NewEncoder(w).Encode(files)
		return
	}
	tmp_path:=f.Name()
	h := sha256.New()
	mw := io.MultiWriter(h, f)
	n, err := io.Copy(mw, io.LimitReader(file, 15<<20+1))
	if err != nil {
		f.Close()
		os.Remove(tmp_path)
		files := Img_file{
			State: err.Error(),
			OK: false,
		}
		json.NewEncoder(w).Encode(files)
		return
	}
	if n > 15<<20 {
		f.Close()
		os.Remove(tmp_path)
		files := Img_file{
			State: "too big",
			OK: false,
		}
		json.NewEncoder(w).Encode(files)
		return
	}
    sum := h.Sum(nil)
	img_file := base64.RawURLEncoding.EncodeToString(sum) + ext
	f.Close()
	imgLock.RLock()
	_,ok:=imgMap[img_file]
    imgLock.RUnlock()
	if ok {
		os.Remove(tmp_path)
		files := Img_file{
        State: "/file/"+img_file,
		OK: true,
		}
		json.NewEncoder(w).Encode(files)
		return
	}
	path := filepath.Join(dataDir, img_file)
	err = os.Rename(tmp_path, path)
	if err!=nil {
		os.Remove(tmp_path)
		files := Img_file{
			State: err.Error(),
			OK: false,
		}
		json.NewEncoder(w).Encode(files)
		return
	}
    imgLock.Lock()
	imgMap[img_file]=struct{}{}
    imgLock.Unlock()
	
	files := Img_file{
        State: "/file/"+img_file,
		OK: true,
    }
	json.NewEncoder(w).Encode(files)
	return
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(strings.TrimPrefix(r.URL.Path, "/file/"))
	imgLock.RLock()
	_,ok:=imgMap[filename]
    imgLock.RUnlock()
	if !ok {
		http.NotFound(w, r)
        return
	}
	fs.ServeHTTP(w, r)
}
func loadImages(dir string) error {
	imgMap = make(map[string]struct{})
    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }
    for _, e := range entries {
        if e.IsDir() {
            continue
        }
        name := e.Name()
        ext := strings.ToLower(filepath.Ext(name))
        if !allowedExt[ext] {
            continue
        }
        imgMap[name]=struct{}{}
    }
    return nil
}

func main() {
	argc:=len(os.Args)
	if argc <=1{
		log.Println("go-imgbed [config.json]")
		return
	}
	cfg,err:=resolv(os.Args[1])
	if err!=nil{
		log.Println("read config err:",err)
		return
	}
	user=cfg.User
	pass=cfg.Pass
	dataDir=cfg.SaveFile
	if dataDir == "" {
		dataDir= "img_data"
	}
	if cfg.Upload == "" {
		cfg.Upload = "/upload"
	}
	os.MkdirAll(dataDir, 0755)
	err = loadImages(dataDir)
	if err != nil {
		log.Println("load image dir err:", err)
	}
	data.Name=cfg.Upload
    http.HandleFunc(cfg.Upload, uploadHandler)
	http.HandleFunc("/file/", viewHandler)
	ssl_on:=false
	if cfg.Crt!="" && cfg.Key!="" {
		ssl_on=true
	}
	fs = http.StripPrefix("/file/", http.FileServer(http.Dir(dataDir)))
	if ssl_on {
		log.Println("server Listen https on",cfg.Listen)
		err := http.ListenAndServeTLS(cfg.Listen, cfg.Crt, cfg.Key, nil)
		if err != nil {
			log.Println(err)
		}
	}else{
		log.Println("server Listen http on",cfg.Listen)
		err = http.ListenAndServe(cfg.Listen, nil)
		if err != nil {
			log.Println(err)
		}
	}
}