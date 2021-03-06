package database

import (
	"encoding/json"
	"pulltg/utils"
	"time"
)

// DataList struct
type DataList struct {
	ID        int64     `json:"id" gorm:"primary_key;column:id"`
	Types     int       `json:"types" gorm:"index:idx_name_types_id;column:types"`
	Keys      string    `json:"keys" gorm:"varchar(128);index:idx_name_keys_id;column:keys"`
	Title     string    `json:"title" gorm:"varchar(128);index:idx_name_title_id;column:title"`
	Total     int       `json:"total" gorm:"column:total"`
	Completed int       `json:"completed" gorm:"column:completed"`
	Path      string    `json:"path" gorm:"varchar(128);column:path"`
	Percent   int       `json:"percent" gorm:"column:percent"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

// TableName change table name
func (DataList) TableName() string {
	return "data"
}

// GetData List
func (datalist *DataList) GetData(types int) (dataList []*DataList, err error) {
	time.Sleep(time.Duration(100) * time.Millisecond)
	if err = Eloquent.Find(&dataList, "types = ?", types).Error; err != nil {
		return
	}
	return
}

// Insert Data
func (datalist *DataList) Insert() error {
	Eloquent.Create(&datalist)
	return nil
}

// UpdateCompleted data
func (datalist *DataList) UpdateCompleted(keys string) (update DataList, err error) {
	time.Sleep(time.Duration(100) * time.Millisecond)
	if err = Eloquent.First(&update, "keys = ?", keys).Error; err != nil {
		return
	}
	datalist.Percent = utils.Round(float64(update.Completed+1) / float64(update.Total) * float64(100))
	datalist.Completed = update.Completed + 1
	if err = Eloquent.Model(&update).Updates(&datalist).Error; err != nil {
		return
	}
	return
}

// UpdateStatus data
func (datalist *DataList) UpdateStatus(keys string) (update DataList, err error) {
	time.Sleep(time.Duration(100) * time.Millisecond)
	if err = Eloquent.First(&update, "keys = ?", keys).Error; err != nil {
		return
	}
	if err = Eloquent.Model(&update).Updates(&datalist).Error; err != nil {
		return
	}
	return
}

// Encode Encode
func Encode(datalist []*DataList) ([]byte, error) {
	var buf []byte
	var err error

	if buf, err = json.Marshal(datalist); err != nil {
		return buf, err
	}
	return buf, nil
}
