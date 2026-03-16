package clause

import (
	"fmt"
	"strings"
)

// generator 定义了一个函数类型，用于生成 SQL 子句。
// 接收任意数量的 interface{} 参数，返回生成的 SQL 片段和对应的参数列表。
type generator func(values ...interface{}) (string, []interface{})

// generators 是一个全局映射，将不同的 SQL 子句类型（Type）映射到对应的生成函数。
var generators map[Type]generator

// init 函数在包初始化时执行，用于注册各种 SQL 子句的生成函数。
func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert   // 注册 INSERT 子句生成器
	generators[VALUES] = _values   // 注册 VALUES 子句生成器
	generators[SELECT] = _select   // 注册 SELECT 子句生成器
	generators[LIMIT] = _limit     // 注册 LIMIT 子句生成器
	generators[WHERE] = _where     // 注册 WHERE 子句生成器
	generators[ORDERBY] = _orderBy // 注册 ORDER BY 子句生成器
	generators[UPDATE] = _update   // 注册 UPDATE 子句生成器
	generators[DELETE] = _delete   // 注册 DELETE 子句生成器
	generators[COUNT] = _count     // 注册 COUNT 子句生成器
}

// genBindVars 生成指定数量的占位符字符串，例如 "?, ?, ?"。
// num 表示需要生成的占位符数量。
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

// _insert 生成 INSERT 子句。
// values[0] 是表名，values[1] 是字段名列表（[]string）。
// 返回格式: "INSERT INTO tableName (field1,field2,...)"。
func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

// _values 生成 VALUES 子句。
// 每个 value 是一个 []interface{}，代表一行数据。
// 返回格式: "VALUES (?,?), (?,?)" 并附带所有参数。
func _values(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars

}

// _select 生成 SELECT 子句。
// values[0] 是表名，values[1] 是字段名列表（[]string）。
// 返回格式: "SELECT field1,field2,... FROM tableName"。
func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

// _limit 生成 LIMIT 子句。
// values 包含限制的数量。
// 返回格式: "LIMIT ?" 并附带参数。
func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

// _where 生成 WHERE 子句。
// values[0] 是条件描述（如 "id = ?"），其余是参数。
// 返回格式: "WHERE condition" 并附带参数。
func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

// _orderBy 生成 ORDER BY 子句。
// values[0] 是排序字段（如 "id DESC"）。
// 返回格式: "ORDER BY field"。
func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

// _update 生成 UPDATE 子句。
// values[0] 是表名，values[1] 是字段名和值的映射（map[string]interface{}）。
// 返回格式: "UPDATE tableName SET field1 = ?, field2 = ?" 并附带所有参数。
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

// _delete 生成 DELETE 子句。
// values[0] 是表名。
// 返回格式: "DELETE FROM tableName"。
func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

// _count 生成 COUNT 子句。
// values[0] 是表名。
// 返回格式: "SELECT count(*) FROM tableName"。
func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}
