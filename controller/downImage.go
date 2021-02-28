package controller

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"pulltg/database"
	"pulltg/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// DownloadImages Download Images
func DownloadImages(l []string, p, port string) bool {
	addr := strings.Join([]string{"localhost", port}, ":")
	u := url.URL{Scheme: "ws", Host: addr, Path: "/api/downlist"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return false
	}
	for i, item := range l {
		wg.Add(1)
		time.Sleep(time.Duration(1) * time.Second)
		go SavePic(item, p, i, conn)
	}
	wg.Wait()
	time.Sleep(time.Duration(1) * time.Second)
	ChangeDataStatus(p, conn)
	return true
}

var wg sync.WaitGroup

//SavePic Save Pic
func SavePic(url, path string, i int, conn *websocket.Conn) {
	defer wg.Add(-1)
	// resp, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	n := strconv.Itoa(i + 1)
	si := ZeroFill(n)
	p := strings.Join([]string{path, "/"}, "")
	tp := strings.Split(url, ".")
	tn := len(tp) - 1
	typ := strings.Join([]string{".", tp[tn]}, "")
	fileName := strings.Join([]string{p, si, typ}, "")
	key := utils.MakeMD5(path)

	var datalist database.DataList
	datalist.UpdateCompleted(key)
	dataList, _ := datalist.GetData(true)
	saveData, _ := database.Encode(dataList)

	time.Sleep(time.Duration(3) * time.Second)
	// go SaveDataToFile(fileName, body)
	go WsWriter(conn, saveData)
	fmt.Println(fileName)
	return
}

// SaveDataToFile save data to file
func SaveDataToFile(dataFile string, saveData []byte) {
	_ = ioutil.WriteFile(dataFile, saveData, 0644)
	return
}

// WsWriter ws writer
func WsWriter(conn *websocket.Conn, saveData []byte) {
	conn.WriteMessage(websocket.TextMessage, saveData)
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
func ChangeDataStatus(path string, conn *websocket.Conn) {
	key := utils.MakeMD5(path)
	var datalist database.DataList
	datalist.UpdateStatus(key)
	dataList, _ := datalist.GetData(false)
	sendData, _ := database.Encode(dataList)
	go WsWriter(conn, sendData)
	return
}
