package main

import (
	"log/slog"
	"net/http"
	"strings"
	"text/template"
)

type User struct {
	Username string
	Password string
}

var registerdUsers = []User{
	{Username: "admin", Password: "admin123"},
}

var halamanUtama = []string{
	"Selamat datang di website",
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("views/base.html", "templates/header.html", "templates/footer.html", "templates/home.html")
	if err != nil {
		http.Error(w, "Gagal memuat template html", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, halamanUtama); err != nil {
		http.Error(w, "Gagal merender html", http.StatusInternalServerError)
	}
}

func TentangHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/tentang.html")
	if err != nil {
		http.Error(w, "Gagal memuat halaman html", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func DestinasiHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/destinasi.html")
	if err != nil {
		http.Error(w, "Gagal memuat halaman destinasi", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

var isLogin = false

func DasboardHandler(w http.ResponseWriter, r *http.Request) {

	if !isLogin {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("views/dasboard.html")
	if err != nil {
		http.Error(w, "Gagal memuat template html", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, halamanUtama); err != nil {
		http.Error(w, "Gagal merender html", http.StatusInternalServerError)
	}

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	isLogin = false
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, "Gagal memuat halaman Register", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func RegisterProsesHandler(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	if username == "" && password == "" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	registerdUsers = append(registerdUsers, User{
		Username: username,
		Password: password,
	})

	println("BERHASIL REGISTER -> User baru ditambahkan:", username)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method tidak di izinkan", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Gagal memuat template html"+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)

}

func LoginProsesHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	loginValid := false
	for _, u := range registerdUsers {
		if u.Username == username && u.Password == password {
			loginValid = true
			break
		}

	}

	if !loginValid {
		isLogin = true
		http.Redirect(w, r, "/dasboard", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// func LoginProses(w http.ResponseWriter, r *http.Request) {
// 	username := r.FormValue("username")
// 	password := r.FormValue("password")

// 	// Cetak ke terminal untuk debugging
// 	println("DEBUG -> Username:", username, "| Password:", password)

// 	if username == "admin" && password == "admin123" {
// 		http.Redirect(w, r, "/dasboard", http.StatusSeeOther)
// 		return
// 	}

// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
// }

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	// USER
	mux.HandleFunc("GET /{$}", HomeHandler)
	mux.HandleFunc("GET /tentang", TentangHandler)
	mux.HandleFunc("GET /destinasi", DestinasiHandler)
	mux.HandleFunc("GET /login", LoginPageHandler)
	mux.HandleFunc("GET /register", RegisterPageHandler)
	mux.HandleFunc("GET /logout", LogoutHandler)
	mux.HandleFunc("POST /register", RegisterProsesHandler)
	mux.HandleFunc("POST /login", LoginProsesHandler)
	mux.HandleFunc("GET /dasboard", DasboardHandler)

	slog.Info("Server berjalan di http://localhost:8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
