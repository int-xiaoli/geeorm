package clause

import "strings"

// Clause 结构体用于存储和管理 SQL 语句的不同子句。
// sql: 存储每个子句类型对应的 SQL 片段字符串。
// sqlVars: 存储每个子句类型对应的参数列表（用于防止 SQL 注入）。
type Clause struct {
	sql     map[Type]string        // 按 Type 存储 SQL 片段
	sqlVars map[Type][]interface{} // 按 Type 存储对应的参数
}

// Type 是一个整数类型，用于表示不同的 SQL 子句类型。
// 使用 iota 自动生成常量值（0,1,2,...）。
type Type int

const (
	INSERT  Type = iota // INSERT 子句 (值为 0)
	VALUES              // VALUES 子句 (值为 1)
	SELECT              // SELECT 子句 (值为 2)
	LIMIT               // LIMIT 子句 (值为 3)
	WHERE               // WHERE 子句 (值为 4)
	ORDERBY             // ORDER BY 子句 (值为 5)
	UPDATE              // UPDATE子句 6
	DELETE              // DELETE子句 7
	COUNT               // COUNT子句 8
)

// Set 方法用于设置指定类型的 SQL 子句。
// name: 子句类型（如 INSERT, WHERE 等）
// vars: 传递给生成器函数的参数（如表名、字段名、条件值等）
// 该方法会调用 generators[name] 对应的生成函数，获取 SQL 片段和参数，并存储到 Clause 中。
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...) // 调用对应的生成器函数
	c.sql[name] = sql                      // 存储 SQL 片段
	c.sqlVars[name] = vars                 // 存储参数
}

// Build 方法根据指定的顺序拼接 SQL 子句，生成完整的 SQL 语句。
// orders: 指定要拼接的子句类型及其顺序（如 []Type{SELECT, WHERE, LIMIT}）
// 返回值: 完整的 SQL 字符串和所有参数的扁平化列表。
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string              // 存储按顺序排列的 SQL 片段
	var vars []interface{}         // 存储按顺序排列的所有参数
	for _, order := range orders { // 遍历指定的子句顺序
		if sql, ok := c.sql[order]; ok { // 如果该子句已设置
			sqls = append(sqls, sql)                 // 添加 SQL 片段
			vars = append(vars, c.sqlVars[order]...) // 添加对应参数（展开）
		}
	}
	return strings.Join(sqls, " "), vars // 用空格连接所有 SQL 片段
}
