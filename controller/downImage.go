package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// DownloadImages Download Images
func DownloadImages(l []string, p string) bool {
	for i, item := range l {
		wg.Add(1)
		go SavePic(item, p, i)
	}
	wg.Wait()
	fmt.Println("over")
	return true
}

var wg sync.WaitGroup

//SavePic Save Pic
func SavePic(url, path string, i int) {
	defer wg.Add(-1)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(body)
	n := strconv.Itoa(i)
	si := ZeroFill(n)
	p := strings.Join([]string{path, "/"}, "")
	fileName := strings.Join([]string{p, si, ".jpg"}, "")
	_ = ioutil.WriteFile(fileName, body, 0755)
	// fmt.Println(fileName)
	/*
		保存json completed + 1
	*/
}

// ZeroFill leading zero fill
func ZeroFill(i string) (x string) {
	if len(i) == 1 {
		x = strings.Join([]string{"000", i}, "")
	} else if len(i) == 2 {
		x = strings.Join([]string{"00", i}, "")
	} else if len(i) == 3 {
		x = strings.Join([]string{"0", i}, "")
	} else {
		x = i
	}
	return
}
