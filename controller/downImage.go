package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"pulltg/database"
	"pulltg/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//  https://blog.csdn.net/taoerchun/article/details/108402296
type WsConn struct {
	*websocket.Dialer
	Mux sync.RWMutex
}

// DownloadImages Download Images
func DownloadImages(l []string, p, port string, length int) bool {
	addr := strings.Join([]string{"localhost", port}, ":")
	u := url.URL{Scheme: "ws", Host: addr, Path: "/api/downlist"}
	var dialer *WsConn

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return false
	}
	for i, item := range l {
		wg.Add(1)
		time.Sleep(time.Duration(100) * time.Millisecond)
		go SavePic(item, p, i, conn, dialer)
	}
	wg.Wait()
	time.Sleep(time.Duration(200) * time.Millisecond)
	ChangeDataStatus(p, conn, length, dialer)
	return true
}

var wg sync.WaitGroup

//SavePic Save Pic
func SavePic(url, path string, i int, conn *websocket.Conn, dialer *WsConn) {
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
	n := strconv.Itoa(i + 1)
	si := ZeroFill(n)
	p := strings.Join([]string{path, "/"}, "")
	tp := strings.Split(url, ".")
	tn := len(tp) - 1
	typ := strings.Join([]string{".", tp[tn]}, "")
	fileName := strings.Join([]string{p, si, typ}, "")
	key := utils.MakeMD5(path)
	time.Sleep(time.Duration(100) * time.Millisecond)
	// time.Sleep(time.Duration(3) * time.Second)
	go SaveDataToFile(fileName, body)
	var datalist database.DataList
	datalist.UpdateCompleted(key)
	dataList, derr := datalist.GetData(0)
	if derr == nil {
		saveData, err := database.Encode(dataList)
		if err == nil {
			WsWriter(conn, saveData, dialer)
		}
	}
	// fmt.Println(fileName)
	return
}

// SaveDataToFile save data to file
func SaveDataToFile(dataFile string, saveData []byte) {
	_ = ioutil.WriteFile(dataFile, saveData, 0644)
	return
}

// WsWriter ws writer
func WsWriter(conn *websocket.Conn, saveData []byte, dialer *WsConn) {
	dialer.Mux.Lock()
	conn.WriteMessage(websocket.TextMessage, saveData)
	dialer.Mux.Unlock()
	return
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

// ChangeDataStatus change data status
func ChangeDataStatus(path string, conn *websocket.Conn, length int, dialer *WsConn) {
	key := utils.MakeMD5(path)
	var datalist database.DataList
	if length == GetFileCount(path) {
		datalist.Types = 1
		datalist.Completed = length
		datalist.Percent = 100
		datalist.UpdateStatus(key)
	}
	time.Sleep(time.Duration(100) * time.Millisecond)
	dataList, err := datalist.GetData(0)
	if err == nil {
		sendData, _ := database.Encode(dataList)
		WsWriter(conn, sendData, dialer)
	}
	return
}

// GetFileCount Get File Count
func GetFileCount(p string) int {
	i := 0
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return 0
	}
	for _, file := range files {
		if !file.IsDir() {
			i++
		}
	}
	return i
}
