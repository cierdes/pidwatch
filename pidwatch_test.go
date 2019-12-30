/*
author: foolbread
file: main/
date: 2019/12/26 4:01 PM
*/
package main

import "testing"

func Test_getPidFromFile(t *testing.T){
	pid := getPidFromFile("my.pid")
	if pid < 0{
		t.Fatal(pid)
	}

	t.Log(pid)
}

func Test_checkProgram(t *testing.T){
	pid := -1
	exist := checkProgram(pid)
	if !exist{
		t.Log(pid,"is not exist!")
	}else{
		t.Log(pid,"is exist!")
	}
}

func Test_forkProgam(t *testing.T){
	err := forkProgam("/Users/foolbread/code/gocode/test/src/test_sso/test_sso",nil)
	if err != nil{
		t.Fatal(err)
	}
	t.Log("start ok!")
}