package main

import (
	"github.com/robertkrimen/otto"
	"fmt"
	"net/http"
	"log"
)
func main() {


	//测试虚拟机能否运行js
	http.HandleFunc("/", sayHttp)
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//parse js function
func parseJSFunc(w http.ResponseWriter,r *http.Request) {
	r.ParseForm()

	for k, v := range r.Form {
		if k =="jstest" {

			//创建虚拟机
			vm := otto.New()
			var src=string(v[0])
			//设置Web中传递的js代码，并在虚拟机运行js
			vm.Run(src)

			//调用js对象
			var obj,_=vm.Object("object1")

			//设置js中的成员变量
			obj.Set("age",1000)

			//调用js成员成员变量
			var age,_=obj.Get("age")
			var name,_=obj.Get("name")
			fmt.Println("name:",name,"age:",age)

			//调用对象中的无参数方法
			var meth,_=obj.Call("sayHi",nil,nil)
			fmt.Println(meth)


			//调用对象中有参数的方法
			var mpar,_=obj.Call("sayHello",11,22)
			fmt.Println(mpar)
		}
	}
	fmt.Fprintf(w, "Hello world!")
}

func sayHttp(w http.ResponseWriter,r *http.Request){
	parseJSFunc(w,r)
}

