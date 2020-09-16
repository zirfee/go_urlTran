package urlServer

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
)

const AddForm = `<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>`

var store Store

func Start(masterAddr, port *string, enableRpc *bool) {
	if *masterAddr == "" {
		store = NewUrlStore("tt.record")
	} else {
		store = NewProxyStore(*masterAddr)
	}
	if *enableRpc {
		rpc.RegisterName("store", store)
		rpc.HandleHTTP()
	}
	http.HandleFunc("/add", add)
	http.HandleFunc("/", redirect)
	http.ListenAndServe(*port, nil)
}

func add(writer http.ResponseWriter, reader *http.Request) {
	/*url := reader.FormValue("url")
	if url==""{
		fmt.Fprintf(writer,AddForm)
		return
	}*/
	key := ""
	url := "noway"
	if err := store.Put(&url, &key); err != nil {
		log.Println(err)
		http.Error(writer, "添加失败", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(writer, "http://localhost:8080/%s", key)
}

func redirect(writer http.ResponseWriter, reader *http.Request) {
	shortUrl := reader.URL.Path[1:]
	trueUrl := ""
	store.Get(&shortUrl, &trueUrl)
	if trueUrl == "" {
		http.NotFound(writer, reader)
		return
	}
	http.Redirect(writer, reader, trueUrl, http.StatusFound)
}
