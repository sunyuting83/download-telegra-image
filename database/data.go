package database

import (
	"bytes"
	"encoding/binary"
	"pulltg/utils"
	"time"
)

// DataList struct
type DataList struct {
	ID        int64     `json:"id" gorm:"primary_key, column:id"`
	Title     string    `json:"title" gorm:"index:idx_name_title, column:title"`
	Total     int       `json:"total" gorm:"column:total"`
	Completed int       `json:"completed" gorm:"column:completed"`
	Keys      string    `json:"keys" gorm:"index:idx_name_key, column:keys"`
	Path      string    `json:"path" gorm:"column:path"`
	Percent   int       `json:"percent" gorm:"column:percent"`
	Type      bool      `json:"type" gorm:"index:idx_name_type, column:type"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

// TableName change table name
func (DataList) TableName() string {
	return "data"
}

// GetData List
func (datalist *DataList) GetData(types bool) (dataList []*DataList, err error) {
	if err = Eloquent.Find(&dataList, "type = ?", types).Error; err != nil {
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
func (datalist *DataList) UpdateCompleted(keys string) (update *DataList, err error) {

	if err = Eloquent.First(&update, "keys = ?", keys).Error; err != nil {
		return
	}
	datalist.Completed = datalist.Completed + 1
	datalist.Percent = utils.Round(float64(datalist.Completed+1) / float64(datalist.Total) * float64(100))
	if err = Eloquent.Model(&update).Updates(&datalist).Error; err != nil {
		return
	}
	return
}

// UpdateStatus data
func (datalist *DataList) UpdateStatus(keys string) (update *DataList, err error) {

	if err = Eloquent.First(&update, keys).Error; err != nil {
		return
	}
	datalist.Type = false
	if err = Eloquent.Model(&update).Updates(&datalist).Error; err != nil {
		return
	}
	return
}

// Encode en code
func Encode(datalist []*DataList) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, datalist); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
