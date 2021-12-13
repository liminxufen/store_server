package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

//time normal define
type TimeNormal struct { //自定义日期格式, 非gorm指定的RFC3339Nano
	time.Time
}

func (t *TimeNormal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	str := string(data)
	timeStr := strings.Trim(str, "\"")
	if len(timeStr) == 0 {
		//timeStr = "1970-01-01 00:00:00"
	}
	if len(timeStr) == 10 {
		timeStr = fmt.Sprintf("%v 00:00:00", timeStr)
	}
	if strings.Contains(timeStr, "0000-00-00") { //兼容历史数据
		timeStr = strings.Replace(timeStr, "0000-00-00", "1970-01-01", -1)
	}
	if strings.Contains(timeStr, "0000:00:00") { //兼容历史数据
		timeStr = strings.Replace(timeStr, "0000:00:00", "1970-01-01 00:00:00", -1)
	}
	t1, err := time.Parse("2006-01-02 15:04:05", timeStr)
	*t = TimeNormal{t1}
	return err
}

func (t TimeNormal) MarshalJSON() ([]byte, error) {
	tune := t.Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
}

//value insert timestamp into mysql need this function
func (t TimeNormal) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

//scan valueof time.Time
func (t *TimeNormal) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = TimeNormal{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

/******************** 使gorm支持[]int64结构 ******************/
type Int64s []int64

func (c Int64s) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *Int64s) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), c)
}

/******************** 使gorm支持[]string结构 ******************/
type Strings []string

func (c Strings) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *Strings) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), c)
}
