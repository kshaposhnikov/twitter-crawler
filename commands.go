package main

import "flag"

type command struct {
	flag flag.Flag
	executor func()
}

type commands []command

var commandCache commands

func (cmd *commands) Init() {
	//commands
}
