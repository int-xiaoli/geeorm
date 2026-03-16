// Package session 放置操作数据库表相关的代码
package session

import (
	"fmt"
	"geeorm/schema"
	"reflect"
	"strings"
)

// Model 设置当前会话要操作的模型（结构体）
// 参数 value: 要操作的结构体实例
// 返回值: 当前 Session 实例，支持链式调用
// 功能: 如果当前没有设置 refTable 或者传入的模型类型与当前不同，
// 则解析新的模型结构并更新 refTable
func (s *Session) Model(value interface{}) *Session {
	// nil or different model,update refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// RefTable 获取当前会话关联的表结构信息
// 返回值: 指向 schema.Schema 的指针，包含表的元数据信息
// 注意: 如果未设置模型会直接 panic，确保先调用 Model() 方法
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		panic("model is not set")
	}
	return s.refTable
}

// CreateTable 根据当前模型创建数据库表
// 返回值: 创建表时可能发生的错误
// 功能: 遍历模型的所有字段，构建 CREATE TABLE SQL 语句并执行
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	// 遍历所有字段，构建每个字段的定义（字段名 + 数据类型 + 标签约束）
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	// 将所有字段定义用逗号连接
	desc := strings.Join(columns, ",")
	// 执行 CREATE TABLE 语句
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// DropTable 删除当前模型对应的数据库表
// 返回值: 删除表时可能发生的错误
// 功能: 执行 DROP TABLE IF EXISTS 语句，安全地删除表（如果表不存在也不会报错）
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

// HasTable 检查当前模型对应的表是否存在
// 返回值: true 表示表存在，false 表示表不存在
// 功能: 使用数据库方言提供的表存在性检查 SQL，查询系统表来验证表是否存在
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
