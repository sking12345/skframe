package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"reflect"
	"skframe/pkg/config"
	"skframe/pkg/logger"
	"strings"
	"sync"
	"time"
)

type DbConnectPool struct {
	DbPtr       *sqlx.DB //连接指针
	Statue      int      //状态 1，正常 2:异常
	LastUseTime int64
	SyncLock    sync.Mutex
	Tx          *sqlx.Tx
}

var masterPool []DbConnectPool
var slavePool []DbConnectPool
var masterUseIndex = 0
var slaveUserIndex = 0

var InUse = 2 //使用中
var Yes = 1   //正常
var No = 0    //不正常

func setupMasterDB() { //主数据
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		config.Get("database.mysql_master.username"),
		config.Get("database.mysql_master.password"),
		config.Get("database.mysql_master.host"),
		config.Get("database.mysql_master.port"),
		config.Get("database.mysql_master.database"),
		config.Get("database.mysql_master.charset"),
	)
	maxConnect := config.GetInt("database.mysql_master.max_connections")
	nowTime := time.Now()
	for i := 0; i < maxConnect; i++ {
		database, err := sqlx.Open("mysql", dsn)
		if err != nil {
			break
		}
		masterPool = append(masterPool, DbConnectPool{
			DbPtr:       database,
			Statue:      Yes,
			LastUseTime: nowTime.Unix(),
		})
	}
	if len(masterPool) <= 0 {
		panic(errors.New("database connection not supported"))
	}
	logger.InfoString("db", "master", fmt.Sprintf("connect:%d", len(masterPool)))

}

func setupSlaverDB() { //主数据
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		config.Get("database.mysql_save.username"),
		config.Get("database.mysql_save.password"),
		config.Get("database.mysql_save.host"),
		config.Get("database.mysql_save.port"),
		config.Get("database.mysql_save.database"),
		config.Get("database.mysql_save.charset"),
	)
	fmt.Println(dsn)
	maxConnect := config.GetInt("database.mysql_master.max_connections")
	nowTime := time.Now()
	for i := 0; i < maxConnect; i++ {
		database, err := sqlx.Open("mysql", dsn)
		if err != nil {
			break
		}
		slavePool = append(slavePool, DbConnectPool{
			DbPtr:       database,
			Statue:      Yes,
			LastUseTime: nowTime.Unix(),
		})
	}
	if len(slavePool) <= 0 {
		panic(errors.New("database connection not supported"))
	}
	logger.InfoString("db", "slave", fmt.Sprintf("connect:%d", len(slavePool)))
}

func ConnectDB() {
	setupMasterDB()
	if config.GetBool("database.master_slave") == true {
		setupSlaverDB()
	}
}

func Destruct() {
	for _, item := range masterPool {
		item.DbPtr.Close()
	}
	if config.GetBool("database.master_slave") == true {
		for _, item := range slavePool {
			item.DbPtr.Close()
		}
	}
}

func CheckOneConnect(db *sqlx.DB) error { //检测当前连接是否正常
	return db.Ping()
}

func CheckAllConnect() {
	timeNowUnix := time.Now().Unix()
	for index, item := range masterPool {
		if item.Statue == InUse {
			continue
		}
		if item.DbPtr.Ping() != nil {
			masterPool = append(masterPool[:index], masterPool[index+1:]...)
		}
		masterPool[index].LastUseTime = timeNowUnix
	}

	if config.GetBool("database.master_slave") == true {
		for index, item := range slavePool {
			if item.Statue == InUse {
				continue
			}
			if item.DbPtr.Ping() != nil {
				slavePool = append(slavePool[:index], slavePool[index+1:]...)
			}
			slavePool[index].LastUseTime = timeNowUnix
		}
	}
}

func GetMasterDB() *DbConnectPool {
	connect := &masterPool[masterUseIndex]
	masterUseIndex = (masterUseIndex + 1) % len(masterPool)
	return connect
}

func GetSlaverDB() *DbConnectPool {
	connect := slavePool[slaveUserIndex]
	slaveUserIndex = (slaveUserIndex + 1) % len(slavePool)
	return &connect
}

type CountStruct struct {
	Count int `db:"count"`
}

func Count(table string, where map[string]interface{}, txCon *sqlx.Tx) int {

	var whereField []string
	var values []interface{}
	for key, item := range where {
		if reflect.ValueOf(item).Kind() == reflect.Map {
			for key1, item1 := range item.(map[string]interface{}) {
				whereField = append(whereField, fmt.Sprintf("%s %s ?", key, key1))
				values = append(values, item1)
			}
		} else {
			values = append(values, item)
			whereField = append(whereField, key+"=?")
		}
	}
	whereStr := strings.Join(whereField, ",")
	sqlStr := fmt.Sprintf("select count(*) as count from %s where %s", table, whereStr)
	var countInfo []CountStruct
	var err error
	if txCon == nil {
		var connectStruct *DbConnectPool
		if config.GetBool("database.master_slave") == true {
			connectStruct = GetSlaverDB()
		} else {
			connectStruct = GetMasterDB()
		}
		connectStruct.SyncLock.Lock()
		connectStruct.Statue = InUse
		err = connectStruct.DbPtr.Select(&countInfo, sqlStr, values...)
		connectStruct.Statue = Yes
		connectStruct.LastUseTime = time.Now().Unix()
		connectStruct.SyncLock.Unlock()

	} else {
		err = txCon.Select(&countInfo, sqlStr, values...)
	}
	if err != nil {
		logger.LogWarnIf(err)
		return 0
	}
	return countInfo[0].Count

}

func Find(table string, filedInfo map[string]string, where map[string]interface{}, txCon *sqlx.Tx) (error, []map[string]interface{}) {

	var whereField []string
	var values []interface{}
	for key, item := range where {
		if reflect.ValueOf(item).Kind() == reflect.Map {
			for key1, item1 := range item.(map[string]interface{}) {
				whereField = append(whereField, fmt.Sprintf("%s %s ?", key, key1))
				values = append(values, item1)
			}
		} else {
			values = append(values, item)
			whereField = append(whereField, key+"=?")
		}
	}
	var findField []string
	for key, _ := range filedInfo {
		findField = append(findField, key)
	}
	whereStr := strings.Join(whereField, ",")
	sqlStr := fmt.Sprintf("select %s from %s where %s", strings.Join(findField, ","), table, whereStr)
	var rows *sql.Rows
	var err error
	if txCon == nil {
		var connectStruct *DbConnectPool
		if config.GetBool("database.master_slave") == true {
			connectStruct = GetSlaverDB()
		} else {
			connectStruct = GetMasterDB()
		}
		connectStruct.SyncLock.Lock()
		connectStruct.Statue = InUse
		rows, err = connectStruct.DbPtr.Query(sqlStr, values...)
		connectStruct.Statue = Yes
		connectStruct.LastUseTime = time.Now().Unix()
		connectStruct.SyncLock.Unlock()
	} else {
		rows, err = txCon.Query(sqlStr, values...)
	}
	if err != nil {
		return err, nil
	}
	columns, _ := rows.Columns()
	length := len(columns)
	result := []map[string]interface{}{}
	for rows.Next() {
		value := make([]interface{}, length)
		columnPointers := make([]interface{}, length)
		for i := 0; i < length; i++ {
			columnPointers[i] = &value[i]
		}
		rows.Scan(columnPointers...)
		data := make(map[string]interface{})
		for i := 0; i < length; i++ {
			columnName := columns[i]
			columnValue := columnPointers[i].(*interface{})
			data[columnName] = *columnValue
			switch filedInfo[columnName] {
			case "string":
				data[columnName] = cast.ToString(data[columnName])
			case "uint":
				data[columnName] = cast.ToUint(data[columnName])
			case "uint64":
				data[columnName] = cast.ToUint64(data[columnName])
			case "uint8":
				data[columnName] = cast.ToUint8(data[columnName])
			case "uint32":
				data[columnName] = cast.ToUint32(data[columnName])
			case "uint16":
				data[columnName] = cast.ToUint16(data[columnName])
			case "int":
				data[columnName] = cast.ToInt(data[columnName])
			case "int64":
				data[columnName] = cast.ToInt64(data[columnName])
			case "int8":
				data[columnName] = cast.ToInt8(data[columnName])
			case "int32":
				data[columnName] = cast.ToInt32(data[columnName])
			case "int16":
				data[columnName] = cast.ToInt16(data[columnName])
			case "float":
				data[columnName] = cast.ToFloat32(data[columnName])
			case "decimal":
				varStr := cast.ToString(data[columnName])
				val, _ := decimal.NewFromString(varStr)
				data[columnName] = val
			default:
				data[columnName] = cast.ToString(data[columnName])
			}
		}
		result = append(result, data)
	}
	return nil, result
}

func Create(table string, data map[string]interface{}, txCon *sqlx.Tx) (int64, error) {
	var field []string
	var value []string
	var sqlData []interface{}
	for key, val := range data {
		field = append(field, key)
		value = append(value, "?")
		sqlData = append(sqlData, val)
	}
	fieldStr := strings.Join(field, ",")
	valueStr := strings.Join(value, ",")
	sqlStr := fmt.Sprintf("insert into %s(%s) values(%s)", table, fieldStr, valueStr)
	var result sql.Result
	var err error
	if txCon == nil {
		connectStruct := GetMasterDB()
		connectStruct.SyncLock.Lock()
		connectStruct.Statue = InUse
		connectStruct.LastUseTime = time.Now().Unix()
		result, err = connectStruct.DbPtr.Exec(sqlStr, sqlData...)
		connectStruct.Statue = Yes
		connectStruct.SyncLock.Unlock()
	} else {
		result, err = txCon.Exec(sqlStr, sqlData...)
	}
	if err != nil {
		logger.LogWarnIf(err)
		return 0, err
	}
	lastId, err1 := result.LastInsertId()

	return lastId, err1

}

func Update(table string, data map[string]interface{}, where map[string]interface{}, txCon *sqlx.Tx) (int64, error) {
	var setField []string
	var values []interface{}
	for key, val := range data {
		setField = append(setField, key+"=?")
		values = append(values, val)
	}
	SetStr := strings.Join(setField, ",")
	var whereField []string
	for key, item := range where {
		if reflect.ValueOf(item).Kind() == reflect.Map {
			for key1, item1 := range item.(map[string]interface{}) {
				whereField = append(whereField, fmt.Sprintf("%s %s ?", key, key1))
				values = append(values, item1)
			}
		} else {
			values = append(values, item)
			whereField = append(whereField, key+"=?")
		}
	}
	whereStr := strings.Join(whereField, ",")
	sqlStr := fmt.Sprintf("update %s set %s where %s", table, SetStr, whereStr)
	var result sql.Result
	var err error
	if txCon == nil {
		connectStruct := GetMasterDB()
		connectStruct.SyncLock.Lock()
		connectStruct.Statue = InUse
		connectStruct.LastUseTime = time.Now().Unix()
		result, err = connectStruct.DbPtr.Exec(sqlStr, values...)
		connectStruct.SyncLock.Unlock()
	} else {
		result, err = txCon.Exec(sqlStr, values...)
	}
	if err != nil {
		logger.LogWarnIf(err)
		return 0, err
	}
	rowNum, err1 := result.RowsAffected()
	return rowNum, err1
}

func Del(table string, where map[string]interface{}, txCon *sqlx.Tx) (int64, error) {
	var whereField []string
	var values []interface{}
	for key, item := range where {
		if reflect.ValueOf(item).Kind() == reflect.Map {
			for key1, item1 := range item.(map[string]interface{}) {
				whereField = append(whereField, fmt.Sprintf("%s %s ?", key, key1))
				values = append(values, item1)
			}
		} else {
			values = append(values, item)
			whereField = append(whereField, key+"=?")
		}
	}
	whereStr := strings.Join(whereField, ",")
	sqlStr := fmt.Sprintf("delete  from %s where %s", table, whereStr)
	var result sql.Result
	var err error
	if txCon == nil {
		connectStruct := GetMasterDB()
		connectStruct.SyncLock.Lock()
		connectStruct.Statue = InUse
		connectStruct.LastUseTime = time.Now().Unix()
		result, err = connectStruct.DbPtr.Exec(sqlStr, values...)
		connectStruct.Statue = Yes
		connectStruct.SyncLock.Unlock()
	} else {
		result, err = txCon.Exec(sqlStr, values...)
	}
	if err != nil {
		logger.LogWarnIf(err)
		return 0, err
	}
	rowNum, err1 := result.RowsAffected()
	return rowNum, err1
}

func Begin() (*DbConnectPool, error) {
	connectStruct := GetMasterDB()
	connectStruct.SyncLock.Lock()
	connectStruct.Statue = InUse
	connectStruct.LastUseTime = time.Now().Unix()
	tx, err := connectStruct.DbPtr.Beginx()
	if err != nil {
		logger.LogIf(err)
		connectStruct.Statue = Yes
		connectStruct.SyncLock.Unlock()
	}
	connectStruct.Tx = tx
	return connectStruct, err
}

func Rollback(conTx *DbConnectPool) {
	if conTx.Tx != nil {
		conTx.Tx.Rollback()
	}
	conTx.Statue = Yes
	conTx.SyncLock.Unlock()

}

func Commit(conTx *DbConnectPool) {
	if conTx.Tx != nil {
		conTx.Tx.Commit()
	}
	conTx.Statue = Yes
	conTx.SyncLock.Unlock()
}
