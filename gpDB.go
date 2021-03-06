/**
gp数据库连接
create by gloomy 2017-3-30 15:24:11
*/
package gutil

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// GP数据库连接对象
// create by gloomy 2017-3-30 15:27:26
type GpDBStruct struct {
	DbUser       string //数据库用户名
	DbHost       string //数据库地址
	DbPort       int    //数据库端口
	DbPass       string //数据库密码
	DbName       string //数据库库名
	MaxOpenConns int    // 用于设置最大打开的连接数，默认值为0表示不限制
	MaxIdleConns int    // 用于设置闲置的连接数
}

// GP数据库连接
// create by gloomy 2017-3-30 15:29:12
func GpSqlConntion(model GpDBStruct) (*sql.DB, error) {
	dbClause := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", model.DbUser, model.DbPass, model.DbHost, model.DbPort, model.DbName)
	db, err := sql.Open("postgres", dbClause)
	if err != nil {
		return nil, fmt.Errorf("gp can't connection gpClause: %s err: %s ", dbClause, err.Error())
	}
	db.SetMaxOpenConns(model.MaxOpenConns)
	db.SetMaxIdleConns(model.MaxIdleConns)
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("gp can't ping dbClause: %s err: %s ", dbClause, err.Error())
	}
	return db, err
}

// GP数据库关闭
// create by gloomy 2017-3-30 15:39:10
func GpSqlClose(db *sql.DB) {
	if db == nil {
		return
	}
	db.Close()
}

// 查询方法
// create by gloomy 2017-3-30 16:04:53
// dbs数据库连接对象 model数据库对象 sqlStr 要执行的sql语句 param执行SQL的语句参数化传递
// 查询返回条数  错误对象输出
func GpSqlSelect(dbs *sql.DB, model GpDBStruct, sqlStr string, param ...interface{}) (*sql.Rows, error) {
	var (
		row *sql.Rows
		err error
	)
	if dbs == nil || dbs.Ping() != nil {
		GpSqlClose(dbs)
		dbs, err = GpSqlConntion(model)
		if err != nil {
			return nil, fmt.Errorf("GpSqlSelect gpConntion err! err: %s ", err.Error())
		}
	}
	if param == nil {
		row, err = dbs.Query(sqlStr)
	} else {
		row, err = dbs.Query(sqlStr, param...)
	}
	if err != nil {
		return nil, fmt.Errorf("GpSqlSelect query can't select sql: %s err: %s \n", sqlStr, err.Error())
	}
	return row, nil
}

/**
数据库运行方法
创建人:邵炜
创建时间:2015年12月29日17:33:06
修正时间：2017年03月11日16:21:36
输入参数: dbs数据库连接对象 model数据库对象 sqlStr 要执行的sql语句  param执行SQL的语句参数化传递
输出参数: 执行结果对象  错误对象输出
*/
func GpSqlExec(dbs *sql.DB, model GpDBStruct, sqlStr string, param ...interface{}) (sql.Result, error) {
	var (
		exec sql.Result
		err  error
	)
	if dbs == nil || dbs.Ping() != nil {
		GpSqlClose(dbs)
		dbs, err = GpSqlConntion(model)
		if err != nil {
			return nil, fmt.Errorf("GpSqlSelect gpConntion err! err: %s ", err.Error())
		}
	}

	if param == nil {
		exec, err = dbs.Exec(sqlStr)
	} else {
		exec, err = dbs.Exec(sqlStr, param...)
	}

	if err != nil {
		return nil, fmt.Errorf("GpSqlExec query can't select sql: %s err: %s \n", sqlStr, err.Error())
	}

	return exec, err
}

// 查询所有字段值
// create by gloomy 2017-5-12 16:38:58
func GPSelectUnknowColumn(dbs *sql.DB, model GpDBStruct, sqlStr string, param ...interface{}) (*[]string, *[][]string, error) {

	var (
		row             *sql.Rows
		err             error
		columnNameArray []string
		dataList        [][]string
	)
	if dbs == nil || dbs.Ping() != nil {
		MySqlClose(dbs)
		dbs, err = GpSqlConntion(model)
		if err != nil {
			return nil, nil, err
		}
	}
	if param == nil {
		row, err = dbs.Query(sqlStr)
	} else {
		row, err = dbs.Query(sqlStr, param...)
	}
	if err != nil {
		return nil, nil, err
	}
	defer row.Close()
	columnNameArray, err = row.Columns()
	if err != nil {
		return nil, nil, err
	}
	values := make([]sql.RawBytes, len(columnNameArray))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for row.Next() {
		var rowData []string
		err = row.Scan(scanArgs...)
		if err != nil {
			continue
		}
		for _, col := range values {
			if col == nil {
				rowData = append(rowData, "")
			} else {
				rowData = append(rowData, string(col))
			}
		}
		dataList = append(dataList, rowData)
	}
	return &columnNameArray, &dataList, nil
}
