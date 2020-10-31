/*
SQLの処理の流れ
①sqlの解析ー＞解析済みのsql情報が存在していないかチェックする。存在していない場合、sqlの検証を行う
->実行計画の作成を行う->作成されたsql情報を共有プールに格納
➁sqlの実行
③行の取得
*/

package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
	if err != nil {
		panic(err)
	}
}

/*　sql　行の取得*/
func retrieve(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRow("select id,content,author from posts where id=$1", id).Scan(&post.Id, &post.Content, &post.Author)
	return

}

/*sqlの解析*/
func (post *Post) create(err error) {

	statement := "insert into posts (content,author) values returning id"
	/*共有プールに存在するか確認→解析結果をオブジェクトにして返す*/
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	return
}

func (post *Post) update() (err error) {
	//update データベース名、set 列名
	_, err = Db.Exec("update posts set content=$2 ,author=$3 where id=$1 ", post.Id, post.Content, post.Author)
	return
}

func (post *Post) delete() (err error) {
	_, err = Db.Exec("delete from posts where id=$1", post.Id)
	return
}
