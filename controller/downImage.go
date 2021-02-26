package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"pulltg/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// DownloadImages Download Images
func DownloadImages(l []string, p, dataFile, port, doneFileName string) bool {
	addr := strings.Join([]string{"localhost", port}, ":")
	u := url.URL{Scheme: "ws", Host: addr, Path: "/api/downlist"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return false
	}
	for i, item := range l {
		wg.Add(1)
		// time.Sleep(time.Duration(3) * time.Second)
		go SavePic(item, p, i, dataFile, conn)
	}
	wg.Wait()
	ChangeDataStatus(dataFile, p, doneFileName, conn)
	return true
}

var wg sync.WaitGroup

//SavePic Save Pic
func SavePic(url, path string, i int, dataFile string, conn *websocket.Conn) {
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
	data := utils.GetDataFile(dataFile)
	for i, item := range data {
		if item.Key == key {
			precent := utils.Round(float64(item.Completed+1) / float64(item.Total) * float64(100))
			data[i] = &utils.SaveData{Total: item.Total, Completed: item.Completed + 1, Key: item.Key, Path: item.Path, Percent: precent}
			break
		}
	}
	saveData, _ := json.Marshal(data)
	_ = ioutil.WriteFile(dataFile, saveData, 0644)

	time.Sleep(time.Duration(3) * time.Second)
	_ = ioutil.WriteFile(fileName, body, 0644)
	conn.WriteMessage(websocket.TextMessage, saveData)
	// fmt.Println(fileName)
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
func ChangeDataStatus(dataFile, path, doneFileName string, conn *websocket.Conn) {
	key := utils.MakeMD5(path)
	data := utils.GetDataFile(dataFile)
	done := utils.GetDataFile(doneFileName)
	Arrlen := len(data)
	for i, item := range data {
		if item.Key == key {
			if i == Arrlen-1 {
				data = data[0:i]
				done = append(done, &utils.SaveData{Total: item.Total, Completed: item.Completed, Key: item.Key, Path: item.Path, Percent: 100})
				break
			} else {
				data = append(data[0:i], data[i+1:]...)
				done = append(done, &utils.SaveData{Total: item.Total, Completed: item.Completed, Key: item.Key, Path: item.Path, Percent: 100})
				break
			}
		}
	}
	saveData, _ := json.Marshal(data)
	_ = ioutil.WriteFile(dataFile, saveData, 0644)
	doneData, _ := json.Marshal(done)
	_ = ioutil.WriteFile(doneFileName, doneData, 0644)
	conn.WriteMessage(websocket.TextMessage, saveData)
	return
}
