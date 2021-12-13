package dblogic

import (
	osql "database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/store_server/dbtools/driver"
	"github.com/store_server/logger"
	"reflect"
	"strings"
)

//base driver
type BaseDriver struct {
	*driver.CMSDriver
	opDB *gorm.DB
}

func (bd *BaseDriver) clone() *BaseDriver {
	cv := &BaseDriver{
		CMSDriver: bd.CMSDriver,
		opDB:      bd.opDB.Begin(),
	}
	return cv
}

func (bd *BaseDriver) bBegin() (cv *BaseDriver, err error) { //开始事务
	if bd.opDB == nil {
		return nil, fmt.Errorf("invalid db pointer")
	}
	cv = bd.clone()
	err = cv.opDB.Error
	logger.Entry().Debugf("BaseDriver[db: %p] begin transaction...", cv.opDB)
	return
}

func (bd *BaseDriver) bRollback() { //支持回滚
	logger.Entry().Debugf("BaseDriver[db: %p] rollback...", bd.opDB)
	if err := bd.opDB.Rollback().Error; err != nil {
		logger.Entry().Errorf("BaseDriver rollback err: %v", err)
	}
}

func (bd *BaseDriver) bCommit() { //提交事务
	logger.Entry().Debugf("BaseDriver[db: %p] commit...", bd.opDB)
	if err := bd.opDB.Commit().Error; err != nil {
		logger.Entry().Errorf("BaseDriver commit err: %v", err)
		bd.Rollback()
	}
}

/*-------------------------- 通用属性方法封装 -------------------------*/
//原生query语句, 返回所有字段
func (bd *BaseDriver) ExecRawQuerySql(sql string, page, pagesize int64,
	model interface{}) ([]interface{}, int64, error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Entry().Errorf("exec raw sql query error: %v", err)
		}
	}()
	if len(sql) == 0 {
		return nil, 0, gorm.ErrInvalidSQL
	}
	if strings.HasSuffix(sql, ";") {
		sql = strings.TrimSuffix(sql, ";")
	}
	if page == 0 {
		page = 1
	}
	if pagesize == 0 {
		pagesize = 10000
	}
	modelType := reflect.TypeOf(model)
	//logger.Entry().Errorf("get model type: %v", modelType)
	var total int64
	if page == 1 {
		if strings.Contains(sql, "where") {
			var tmp []interface{}
			total = bd.opDB.Model(model).Raw(sql).Scan(&tmp).RowsAffected
			tmp = nil
		} else {
			bd.opDB.Model(model).Count(&total)
		}
	}
	offset := (page - 1) * pagesize
	var err error
	var rows *osql.Rows
	if strings.Contains(strings.ToLower(sql), "limit") {
		rows, err = bd.opDB.Model(model).Raw(sql).Rows()
	} else {
		rows, err = bd.opDB.Model(model).Offset(offset).Limit(pagesize).Raw(sql).Rows()
	}
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()
	res := make([]interface{}, 0, pagesize)
	for rows.Next() {
		r := reflect.New(modelType.Elem())
		t := r.Elem().Addr().Interface()
		err = bd.opDB.ScanRows(rows, t)
		if err != nil {
			return nil, total, err
		}
		res = append(res, t)
	}
	return res, total, nil
}

//原生query语句, 返回自定义字段
func (bd *BaseDriver) ExecRawQuerySqlCustomized(sql string, page,
	pagesize int64) ([][]interface{}, error) {
	results := make([][]interface{}, 0)
	if len(sql) <= 0 {
		return nil, gorm.ErrInvalidSQL
	}
	if strings.HasSuffix(sql, ";") {
		sql = strings.TrimSuffix(sql, ";")
	}
	if page == 0 {
		page = 1
	}
	if pagesize == 0 {
		pagesize = 10000
	}
	offset := (page - 1) * pagesize
	var err error
	var rows *osql.Rows
	if strings.Contains(strings.ToLower(sql), "limit") {
		rows, err = bd.opDB.Raw(sql).Rows()
	} else {
		rows, err = bd.opDB.Raw(sql).Offset(offset).Limit(pagesize).Rows()
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]osql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		tmp := make([]interface{}, 0)
		for _, v := range values {
			tmp = append(tmp, string(v))
		}
		results = append(results, tmp)
	}
	return results, err
}

//原生update, insert or delete语句
func (bd *BaseDriver) ExecRawUpdateOrInsertOrDeleteSql(sql string,
	model interface{}) error {
	if len(sql) == 0 {
		return gorm.ErrInvalidSQL
	}
	err := bd.opDB.Model(model).Exec(sql).Error
	if err != nil {
		return gorm.ErrInvalidSQL
	}
	tx, err := bd.bBegin()
	if err != nil {
		tx.bRollback()
	} else {
		tx.bCommit()
	}
	return err
}

//原生join query语句
func (bd *BaseDriver) JoinQueryWithSql(sql string, page,
	pagesize int64) ([][]interface{}, error) {
	results := make([][]interface{}, 0)
	if len(sql) <= 0 {
		return nil, gorm.ErrInvalidSQL
	}
	if strings.HasSuffix(sql, ";") {
		sql = strings.TrimSuffix(sql, ";")
	}
	if page == 0 {
		page = 1
	}
	if pagesize == 0 {
		pagesize = 10000
	}
	offset := (page - 1) * pagesize
	var err error
	var rows *osql.Rows
	if strings.Contains(strings.ToLower(sql), "limit") {
		rows, err = bd.opDB.Raw(sql).Rows()
	} else {
		rows, err = bd.opDB.Raw(sql).Offset(offset).Limit(pagesize).Rows()
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]osql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		tmp := make([]interface{}, 0)
		for _, v := range values {
			tmp = append(tmp, string(v))
		}
		results = append(results, tmp)
	}
	return results, err
}

//导出model所有记录
func (bd *BaseDriver) ExportAllRecords(sql string) ([][]interface{}, error) {
	rows, err := bd.opDB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]osql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	res := make([][]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		tmp := make([]interface{}, 0)
		for _, v := range values {
			tmp = append(tmp, string(v))
		}
		res = append(res, tmp)
	}
	if err = rows.Err(); err != nil {
		return res, err
	}
	return res, nil
}

//通用model query
func (bd *BaseDriver) QueryWithModel(model interface{}, conds map[string]interface{},
	page, pagesize int64) ([]interface{}, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pagesize <= 0 {
		pagesize = 100
	}
	modelType := reflect.TypeOf(model)
	var total int64
	var err error
	if page == 1 {
		if len(conds) != 0 {
			var tmp []interface{}
			total = bd.opDB.Model(model).Where(conds).Scan(&tmp).RowsAffected
			tmp = nil
		} else {
			bd.opDB.Model(model).Count(&total)
		}
	}
	offset := (page - 1) * pagesize
	rows, err := bd.opDB.Model(model).Where(conds).Offset(offset).Limit(pagesize).Rows()
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()
	res := make([]interface{}, 0, pagesize)
	for rows.Next() {
		r := reflect.New(modelType.Elem())
		t := r.Elem().Addr().Interface()
		err = bd.opDB.ScanRows(rows, t)
		if err != nil {
			return nil, total, err
		}
		res = append(res, t)
	}
	return res, total, err
}

//通用model scope
func (bd *BaseDriver) WithScope(scope [3]interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(scope) < 3 {
			return db
		}
		sql := "%s %s %v"
		field := scope[0].(string)
		op := scope[1].(string)
		value := scope[2]
		switch op {
		case "gt":
			op = ">"
		case "ge":
			op = ">="
		case "eq":
			op = "=="
		case "ne":
			op = "!="
		case "lt":
			op = "<"
		case "le":
			op = "<="
		}
		sql = fmt.Sprintf(sql, field, op, value)
		return db.Where(sql)
	}
}

//通用model update
func (bd *BaseDriver) UpdateWithModel(model interface{}, conds interface{},
	updateAttrs interface{}, args ...interface{}) (int64, error) {
	tx, err := bd.bBegin()
	if err != nil {
		return 0, err
	}
	ret := tx.opDB.Model(model).Where(conds, args...).Update(updateAttrs)
	if ret.Error != nil {
		tx.bRollback()
		return 0, ret.Error
	}
	tx.bCommit()
	return ret.RowsAffected, nil
}

//通用model insert
func (bd *BaseDriver) InsertWithModel(modelValue interface{}) (int64, error) {
	tx, err := bd.bBegin()
	if err != nil {
		return 0, err
	}
	db := tx.opDB.Create(modelValue)
	if db.Error != nil {
		tx.bRollback()
	} else {
		tx.bCommit()
	}
	return db.RowsAffected, err
}

//通用model delete with id
func (bd *BaseDriver) DeleteWithModelID(model interface{},
	id interface{}) (int64, error) {
	tx, err := bd.bBegin()
	if err != nil {
		return 0, err
	}
	var ret *gorm.DB
	ret = tx.opDB.Delete(model, id)
	if ret.Error != nil {
		tx.bRollback()
		return 0, ret.Error
	}
	if ret.RowsAffected != 1 {
		tx.bRollback()
		return 0, fmt.Errorf("%d rows will be delete, please confirm it", ret.RowsAffected)
	}
	tx.bCommit()
	return ret.RowsAffected, nil
}

//通用model delete
func (bd *BaseDriver) DeleteWithModel(model interface{},
	conds interface{}, args ...interface{}) (int64, error) {
	tx, err := bd.bBegin()
	if err != nil {
		return 0, err
	}
	ret := tx.opDB.Where(conds, args...).Delete(model)
	if ret.Error != nil {
		tx.bRollback()
		return 0, ret.Error
	}
	if ret.RowsAffected > 100 {
		tx.bRollback()
		return 0, fmt.Errorf("%d rows wiil be delete, please confirm it", ret.RowsAffected)
	}
	tx.bCommit()
	return ret.RowsAffected, nil
}
