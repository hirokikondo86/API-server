package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

/********************************************
				DATABASE接続
*********************************************/
func DbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "go_sample"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

/********************************************
				 SELECT ALL
*********************************************/
func ShowAll(w http.ResponseWriter, r *http.Request) {
	db := DbConn()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users")
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}
	users := Users{}
	var id int
	var name, mail, pass string
	var created, modified string
	for rows.Next() {
		if err := rows.Scan(&id, &name, &mail, &pass, &created, &modified); err != nil {
			panic(err.Error())
		}
		users.Id = id
		users.Name = name
		users.Mail = mail
		users.Pass = pass
		users.Created = created
		users.Modified = modified
		jsonBytes, err := json.Marshal(users)
		if err != nil {
			fmt.Println("Marshl error:", err)
			return
		}
		fmt.Fprint(w, string(jsonBytes), "\n")
	}
}

/********************************************
					SELECT
*********************************************/
func Show(w http.ResponseWriter, r *http.Request) {
	db := DbConn()
	defer db.Close()
	vars := mux.Vars(r)
	nid := vars["id"]
	row, err := db.Query("SELECT * FROM users WHERE id = ?", nid)
	defer row.Close()
	if err != nil {
		panic(err.Error())
	}
	users := Users{}
	var id int
	var name, mail, pass string
	var created, modified string

	row.Scan(&id, &name, &mail, &pass, &created, &modified)

	users.Id = id
	users.Name = name
	users.Mail = mail
	users.Pass = pass
	users.Created = created
	users.Modified = modified

	jsonBytes, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Marshl error:", err)
		return
	}
	fmt.Fprint(w, string(jsonBytes), "\n")
}

/********************************************
					INSERT
*********************************************/
func Insert(w http.ResponseWriter, r *http.Request) {
	db := DbConn()
	defer db.Close()

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	user := Users{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ins, err := db.Prepare("INSERT INTO users (user_name, mail_address, password, created, modified) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	ins.Exec(user.Name, user.Mail, user.Pass, time.Now(), time.Now())
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
}

/********************************************
				  UPDATE
*********************************************/
func Update(w http.ResponseWriter, r *http.Request) {
	db := DbConn()
	defer db.Close()

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	user := Users{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	upd, err := db.Prepare("UPDATE users SET user_name=? , mail_address=?, password=?, modified=? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	vars := mux.Vars(r)
	upd.Exec(user.Name, user.Mail, user.Pass, time.Now(), vars["id"])
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
}

/********************************************
				  DELETE
*********************************************/
func Delete(w http.ResponseWriter, r *http.Request) {
	db := DbConn()
	defer db.Close()
	vars := mux.Vars(r)
	del, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	del.Exec(vars["id"])
}
