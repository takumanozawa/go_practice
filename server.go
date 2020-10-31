/*
####################HTTPリクエスト#########################################
1,リクエスト行　GET(POST)（メソッドという）+URl+使用されているhttpのバージョン
2,リクエストヘッダ(1,2行)
[<名前>:<値>]という形で配信される
Accept:クライアントが受け取りたい形（例としてhtmlなど）
Accept-Cartest:クライアントが受け取りたい形の文字タイプ
Authorization:ベーシック認証の証明書をサーバーに送付する時に使用される
Cookie:クライアントのブラウザに設定されているクッキーを送り返さなければならない
Comntent-length:リクエスト本体の長さ
Host:サーバー名とポート番号
Referrer:リクエストされたページにリンクしていた以前のアドレス
User-Agent:呼び出しているクライアント
３、空行
4,ボディ
補足のメモ書きみたいなもの
first_name=takuma&charaset=utf-8など&で情報を繋げて入れる

####################HTTPレスポンス#######################################
１、ステータス行
１XX：サーバーがリクエストを受信していて、処理し始めている
２XX：成功。クライアントの要求が受け入れられた
３XX：受け入れられたが、クライアントのやることが残っている
4XX:クライアントエラー（クライアントのリクエストに謝り）
５XX:サーバーエラー

2,レスポンスヘッダ
Allow:サーバーがどのリクエストメソッドをサポートしているかを伝える
Location:次のURLをどこに要求すればいいか伝える
Server:レスポンスを返しているサーバーのドメイン名
Set-Cookie：クライアントのクッキーを設定する
WWW-Authenticate：リクエストヘッダのAuthorizationにどのような認証情報を含めるべきか決める

3,空行
４、ボディ
htmlなど応答の中身が入っている



##http.HandleとはServeHTTP関数を持つだけのインターフェイス
type Handler interface{
	//requestとwriterを受け取って返すだけの関数
	ServeHTTP(ResponseWriter,*Request)
}
##http.HandleとはURLに対応するhttp.HandlerをDefaultServeMuxに登録するだけの関数
(DefaultServeMuxとはデフォルトでhttp.ServeMux型の構造体でURLに対応した、http.Handlerを実行するいわゆるルーター)
DefaultServeMuxも広義の意味ではServeHTTP 関数を持つ http.Handler。http.ListenAndServe関数が	nilの場合、DefaultServeMux が Handler として指定されます。

func Handle(pattern string,handler Handler){
	DefaultServeMux.Handle(pattern,handler)
}
##http.HandlerFuncとはfunc(ResponseWriter,*Request)の別名の型で構造体を宣言することなくhttp.Handlerを用意することができます。

##ServeMux→リクエストを登録済みのurlパターンリストと照合して、マッチしたハンドラーを呼び出す
*/
package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	//pathパッケージはスラッシュ区切りのファイルパスを実装するためのものを実装している
	//Base関数はスラッシュ区切りの一番最後の値を返す
	"path"
	"strconv"
)

type Post struct {
	Db      *sql.DB
	Id      int    `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author" `
}

/*handleRequestでそれぞれのCRUD関数に作業をふり分ける*/

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = handleGet(w, r)
	case "POST":
		err = handlePost(w, r)
	case "PUT":
		err = handlePut(w, r)
	case "DELETE":
		err = handleDelete(w, r)

	}
	if err != nil {
		//func Error(w ResponseWriter, error string, code int)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	//selectでidの列を取得する
	post, err := retrieve(id)

	if err != nil {
		return
	}
	//postのメソッドretrieve()で取得した列json形式に変換してhttpに書き込んで返す
	output, err := json.MarshalIndent(&post, "", "\t\t")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

}

func handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	//Contentlengthには渡すデータのサイズが入っている
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var post Post
	json.Unmarshal(body, &post)
	err = post.update()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

func handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	post, err := retrieve(id)
	if err != nil {
		return
	}
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	//jsonを格納するための構造体を作成し、実際のjsonデータを読み込み格納する
	json.Unmarshal(body, &post)
	err = post.update()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return

}

func handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	post, err := retrieve(id)
	if err != nil {
		return
	}
	err = post.delete()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}
func main() {
	server := http.Server{
		Addr: ":" + os.Getenv("PORT"), //
	}
	http.HandleFunc("/post/", handleRequest)
	server.ListenAndServe()

}
