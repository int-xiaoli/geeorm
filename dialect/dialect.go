// 设置不同类型数据库的方言，并注册
package dialect

import "reflect"

// dialectsMap 存储所有已注册的数据库方言，键为驱动名称（如 "sqlite3"），值为对应的方言实现。
var dialectsMap = map[string]Dialect{}

// Dialect 接口定义了数据库方言必须实现的方法。
// 不同的数据库（如 MySQL、PostgreSQL、SQLite）在数据类型和 SQL 语法上存在差异，
// 通过实现此接口可以适配不同数据库的特性。
type Dialect interface {
	// DataTypeOf 根据 Go 类型返回对应的数据库字段类型。
	// 例如：Go 的 int 类型在 SQLite 中对应 "integer"，在 MySQL 中可能对应 "int"。
	DataTypeOf(typ reflect.Value) string

	// TableExistSQL 生成检查表是否存在的 SQL 语句。
	// 不同数据库检查表存在的语法不同，例如：
	// - SQLite: SELECT name FROM sqlite_master WHERE type='table' AND name = ?
	// - MySQL: SELECT table_name FROM information_schema.tables WHERE table_name = ?
	TableExistSQL(tableName string) (string, []interface{})
}

// RegisterDialect 注册一个新的数据库方言。
// name: 数据库驱动名称（如 "sqlite3"、"mysql"）
// dialect: 方言的具体实现（实现了 Dialect 接口的结构体）
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect 根据驱动名称获取对应的方言实现。
// name: 数据库驱动名称
// 返回值: 方言实现和是否存在该方言的布尔值
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
