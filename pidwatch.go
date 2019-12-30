/*
author: foolbread
file: pidwatch/pidwatch.go
date: 2019/12/26 3:50 PM
*/
package main

import (
	"flag"
	"io/ioutil"
	"fmt"
	"strconv"
	"syscall"
	"os/exec"
	"time"
	"os/signal"
	"os"
)

func init(){
	flag.String("h","","pidwatch <pidfile name> <command> [<cmdarg1> ...]")
	flag.Parse()
	//to do parser args
	args := flag.Args()
	fmt.Println("args:",args)
	if len(args) < 2{
		fmt.Println("the args are not enough!")
		os.Exit(0)
	}
	pid_file = args[0]
	cmd_program = args[1]
	if len(args) > 2{
		cmd_args = args[2:]
	}
}

var pid_file string
var cmd_program string
var cmd_args []string

func main(){
	go run()
	signalHandle()
}

func run(){
	ti := time.Tick(5*time.Second)
	t := time.Now()
	for {
		pid :=getPidFromFile(pid_file)
		fmt.Println(t.Format("2006-01-02 15:04:05"),"start check program:",cmd_program,"pid:",pid)

		exist := checkProgram(pid)
		//double check
		if !exist{
			time.Sleep(2*time.Second)
			pid = getPidFromFile(pid_file)
			exist = checkProgram(pid)
			if !exist{
				forkProgam(cmd_program,cmd_args)
				fmt.Println("pid:",pid,"program is not exist!")
				fmt.Println("new start program:",cmd_program)
			}
		}

		fmt.Println(t.Format("2006-01-02 15:04:05"),"finish check program:",cmd_program,"pid:",pid)
		t = <-ti
	}
}

func signalHandle(){
	ch := make(chan os.Signal)
	signal.Notify(ch,syscall.SIGINT,syscall.SIGTERM,syscall.SIGHUP,syscall.SIGCHLD)
	sig := <-ch

	switch sig {
	case syscall.SIGINT,syscall.SIGTERM,syscall.SIGHUP:
		fmt.Printf("received sig=%v\n",sig)
		syscall.Kill(getPidFromFile(pid_file),syscall.SIGTERM)
		os.Exit(0)
	case syscall.SIGCHLD:
		fmt.Println("received child sig!")
	}
}

func getPidFromFile(pidfile string)int{
	ret := -1
	data,err := ioutil.ReadFile(pidfile)
	if err != nil{
		fmt.Println(err)
		return ret
	}

	ret,err = strconv.Atoi(string(data))
	if err != nil{
		fmt.Println(err)
		return ret
	}

	return ret
}

func checkProgram(pid int)bool{
	if err := syscall.Kill(pid,0); err != nil{
		fmt.Println(err)
		return false
	}
	return true
}

func forkProgam(program string, programargs []string)error{
	cmd :=exec.Command(program,programargs...)
	return cmd.Start()
}