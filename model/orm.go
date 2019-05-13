// orm
// 反射,运行时动态获取对象信息的方法
package model

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// 表信息
type TableInfo struct {
	Name   string
	Fields []*FiledInfo
}

// 字段信息
type FiledInfo struct {
	Name     string
	Value    reflect.Value
	Relation map[string]string
}

// 实体信息
type ModelInfo struct {
	Name  string
	Model interface{}
}

// 实体映射，key为表名，value为实体信息
var ModelMapping map[string]ModelInfo

// 表信息解析
func getTableInfo(model interface{}) (*TableInfo, error) {
	tableinfo := &TableInfo{}
	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}
	// TableName
	tableNameSplit := strings.Split(modelType.String(), ".")
	tableName := tableNameSplit[len(tableNameSplit)-1]
	tableNameLower := strings.ToLower(tableName)
	tableinfo.Name = tableNameLower
	// FieldInfo
	tableinfo.Fields = getFieldsInfo(modelType, modelValue)
	return tableinfo, nil
}

// 字段信息解析
func getFieldsInfo(modelType reflect.Type, modelValue reflect.Value) []*FiledInfo {
	var fields []*FiledInfo
	for i := 0; i < modelType.NumField(); i++ {
		// 实例化FiledInfo
		filedinfo := &FiledInfo{}
		filedinfo.Relation = make(map[string]string)
		filed := modelType.Field(i)
		// Name
		filedinfo.Name = filed.Name
		// Value
		filedinfo.Value = modelValue.Field(i)
		// Relation
		if strings.Index(string(modelType.Field(i).Tag), ":") != -1 {
			filedinfo.Relation[modelType.Field(i).Name] = modelType.Field(i).Tag.Get("name")
		}
		fields = append(fields, filedinfo)
	}
	return fields
}

// 实体注册到ModelMapping
func Register(model interface{}) {
	if ModelMapping == nil {
		ModelMapping = make(map[string]ModelInfo)
	}
	tbInfo, _ := getTableInfo(model)
	ModelMapping[tbInfo.Name] = ModelInfo{
		Name:  tbInfo.Name,
		Model: model,
	}
}

// 创建插入语句
// insert into profile (Name,Gender,Age,Height,Weight) values (?,?,?,?,?)
func generateInsertStatement(model interface{}) (string, []interface{}, error) {
	tbInfo, err := getTableInfo(model)
	if err != nil {
		panic(err)
	}
	if len(tbInfo.Fields) == 0 {
		panic(fmt.Sprintf("no fields"))
	}
	sqlStr := "insert into " + strings.ToLower(tbInfo.Name)
	fieldStr := ""
	valueStr := ""
	var params []interface{}
	for _, fieldInfo := range tbInfo.Fields {
		fieldStr += fieldInfo.Relation[fieldInfo.Name] + ","
		valueStr += "?,"
		params = append(params, fieldInfo.Value.Interface())
	}
	fieldStr = strings.TrimRight(fieldStr, ",")
	valueStr = strings.TrimRight(valueStr, ",")
	sqlStr += " (" + fieldStr + ") values (" + valueStr + ")"
	return sqlStr, params, nil
}

// 创建数据库操作
type DB struct {
	*sql.DB
	ResultChan chan interface{}
}

func (d *DB) Insert(model interface{}) error {
	insertSql, params, _ := generateInsertStatement(model)
	log.Printf("[Hope engine]: <Orm> generating %s", insertSql)
	_, err := d.Exec(insertSql, params...)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) Run() {
	d.ResultChan = make(chan interface{})
	go func() {
		for {
			result := <-d.ResultChan
			err := d.Insert(result)
			if err != nil {
				log.Printf("[Hope engine]: <Orm> save success")
			}
		}
	}()
}

// 数据库操作对象
func NewDB(driver, address string) *DB {
	db, err := sql.Open(driver, address)
	if err != nil {
		panic(err)
	}
	log.Printf("<Hope engine>: connect database success")
	return &DB{DB: db}
}
