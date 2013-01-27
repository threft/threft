#!/usr/bin/env gorun
package main

/**
 * This file is used to quickly compile, clean and test the multiple packages/binaries of 
 * the threft project.
 * Typical setup would be to clone all the repositories at github.com/threft into yourGoDevpath/src/github.com/threft/
 * This would give you the following directories
 * yourGoDevpath/src/github.com/threft/tidm
 * yourGoDevpath/src/github.com/threft/threft-gen-go-golden
 * yourGoDevpath/src/github.com/threft/threft-gen-go
 * yourGoDevpath/src/github.com/threft/threft.github.com
 * yourGoDevpath/src/github.com/threft/threft
 * The shebang instruction at the top of this file requires launchpad.net/gorun to be installed.
 * Use: `go get go get launchpad.net/gorun`
 */

import (
	"fmt"
	"os/exec"
	"runtime"
)

var buildPkgs = []string{
	"tidm",
	"threft-gen-go/gen-go",
	// "threft-gen-go",
	"threft",
}

func main() {
	runtime.GOMAXPROCS(3)
	dones := make([]chan bool, 0)

	for _, pkg := range buildPkgs {
		ch := make(chan bool, 1)
		dones = append(dones, ch)
		go build(pkg, ch)
	}

	success := true
	for _, ch := range dones {
		if val := <-ch; val == false {
			success = false
		}
	}

	if success {
		fmt.Println("== Succesfull built project.")
	} else {
		fmt.Println("== There where built errors.")
	}
}

func build(pkg string, done chan bool) {

	cmd := exec.Command("go", "build")
	cmd.Dir = "../" + pkg
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", output)
		done <- false
		return
	}

	done <- true
}
