package main

import (
	"flag"
	"fmt"
	//"os"
	"time"

	//"github.com/rususlasan/url-checker/pkg/checker"
)

const (
	concurrency = 5
	timeout = time.Second * 5
)

func main() {
	flag.Parse()
	fmt.Println(flag.Args())
	//urlChecker, err := checker.NewChecker(concurrency, timeout)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//urlChecker.Check(flag.Args())
}
