package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"pulltg/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DownloadImages Download Images
func DownloadImages(l []string, p, dataFile string) bool {
	for i, item := range l {
		wg.Add(1)
		time.Sleep(time.Duration(3) * time.Second)
		go SavePic(item, p, i, dataFile)
	}
	wg.Wait()
	ChangeDataStatus(dataFile, p)
	fmt.Println("over")
	return true
}

var wg sync.WaitGroup

//SavePic Save Pic
func SavePic(url, path string, i int, dataFile string) {
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
	data := GetDataFile(dataFile)
	for i, item := range data.Running {
		if item.Key == key {
			data.Running[i] = SaveData{Total: item.Total, Completed: item.Completed + 1, Key: item.Key, Path: item.Path}
			break
		}
	}
	saveData, _ := json.Marshal(data)
	_ = ioutil.WriteFile(dataFile, saveData, 0644)

	time.Sleep(time.Duration(3) * time.Second)
	// _ = ioutil.WriteFile(fileName, body, 0644)
	fmt.Println(fileName)
	/*
		保存json completed + 1
	*/
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

// GetDataFile get data file
func GetDataFile(d string) (j TempData) {
	data, _ := ioutil.ReadFile(d)
	var (
		index int = len(data)
	)
	index = bytes.IndexByte(data, 0)
	if index != -1 {
		data = data[:index]
	}
	if err := json.Unmarshal(data, &j); err != nil {
		return
	}
	return
}

// ChangeDataStatus change data status
func ChangeDataStatus(dataFile, path string) {
	key := utils.MakeMD5(path)
	data := GetDataFile(dataFile)
	Arrlen := len(data.Running)
	for i, item := range data.Running {
		if item.Key == key {
			if i == Arrlen-1 {
				data.Running = data.Running[0:i]
				data.Done = append(data.Done, SaveData{Total: item.Total, Completed: item.Completed, Key: item.Key, Path: item.Path})
				break
			} else {
				data.Running = append(data.Running[0:i], data.Running[i+1:]...)
				data.Done = append(data.Done, SaveData{Total: item.Total, Completed: item.Completed, Key: item.Key, Path: item.Path})
				break
			}
		}
	}
	saveData, _ := json.Marshal(data)
	fmt.Println(string(saveData))
	_ = ioutil.WriteFile(dataFile, saveData, 0644)
	return
}
