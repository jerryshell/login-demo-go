package main

import (
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"jerryshell.cn/login_demo/dao"
	"jerryshell.cn/login_demo/session"
)

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/userinfo", userinfo)
	http.HandleFunc("/logout", logout)
}

func main() {
	log.Println("Server is running at http://localhost:8080/. Press Ctrl+C to stop.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("html/login.html")
	checkError(err)

	err = t.Execute(w, nil)
	checkError(err)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/", 302)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	log.Println("login", username, password)

	user := dao.FindUserByUsernameAndPassword(username, password)
	if user == nil {
		message(w, r, "登录失败！")
		return
	}
	// 登陆成功
	sess := session.GetSession(w, r)
	sess.SetAttr("user", user)
	http.Redirect(w, r, "/userinfo", 302)
}

func message(w http.ResponseWriter, r *http.Request, message string) {
	t, err := template.ParseFiles("html/message.html")
	checkError(err)

	data := make(map[string]string)
	data["Message"] = message
	err = t.Execute(w, data)
	checkError(err)
}

func userinfo(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(w, r)
	user, exist := sess.GetAttr("user")
	if !exist {
		http.Redirect(w, r, "/", 302)
		return
	}
	t, err := template.ParseFiles("html/userinfo.html")
	checkError(err)
	t.Execute(w, user)
}

func logout(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(w, r)
	sess.DelAttr("user")
	http.Redirect(w, r, "/", 302)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
