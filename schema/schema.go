package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

// Field represents a column of database
type Field struct {
	Name string // 字段名
	Type string // 数据库类型（如 "int", "varchar"）
	Tag  string // geeorm 标签（用于约束条件）
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}       // 原始结构体实例
	Name       string            // 表名（结构体名）
	Fields     []*Field          // 所有字段
	FieldNames []string          // 字段名列表
	fieldMap   map[string]*Field // 字段名→Field 的映射（快速查找）
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// dest是表的结构体
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// reflect.ValueOf(dest): 将 interface{} 转换为 reflect.Value
	// reflect.Indirect(): 自动解引用指针，如果 dest 是指针则获取指向的值，如果是值则直接返回
	// .Type(): 从 reflect.Value 获取 reflect.Type，得到结构体的类型信息
	// 整体作用：无论传入的是结构体指针还是结构体值，都能统一获取到结构体类型
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(), // modelType.Name(): 获取结构体的名称作为表名
		fieldMap: make(map[string]*Field),
	}

	// modelType.NumField(): 获取结构体的字段总数
	// 注意：只能对结构体类型调用，对指针类型调用会 panic
	for i := 0; i < modelType.NumField(); i++ {
		// modelType.Field(i): 获取第 i 个字段的 reflect.StructField
		// StructField 包含字段的元数据：Name(字段名)、Type(字段类型)、Tag(字段标签) 等
		p := modelType.Field(i)

		// p.Anonymous: 判断是否为嵌入字段（匿名嵌套的结构体）
		// ast.IsExported(p.Name): 判断字段名是否首字母大写（导出字段）
		// ORM 只处理导出字段，私有字段不映射到数据库
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name, // 字段名，如 "ID", "Name"

				// 这里是核心反射用法，分步解析：
				// 1. p.Type: 字段的 reflect.Type（如 int, string）
				// 2. reflect.New(p.Type): 创建该类型的新实例（零值），返回指针类型的 reflect.Value
				//    例如：p.Type=int → 创建 *int，值为 0
				// 3. reflect.Indirect(): 解引用指针，得到值的 reflect.Value
				//    例如：*int → int
				// 4. d.DataTypeOf(): 将 reflect.Value 传给方言，通过 typ.Kind() 判断类型
				//    返回对应的 SQL 数据类型字符串（如 "integer", "text"）
				// 为什么这么复杂？因为 DataTypeOf 需要 reflect.Value 参数，而不是 reflect.Type
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}

			// p.Tag.Lookup("geeorm"): 读取结构体字段 tag 中标记为 "geeorm" 的值
			// Lookup 比直接 Get 更安全，可以区分"不存在"和"空字符串"
			// 示例：Name string `geeorm:"NOT NULL"` → v = "NOT NULL"
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}

			// 将解析后的字段添加到 Schema 的三个数据结构中：
			// Fields: 有序的字段列表，用于遍历所有字段
			schema.Fields = append(schema.Fields, field)
			// FieldNames: 字段名列表，用于生成 SQL 列名
			schema.FieldNames = append(schema.FieldNames, p.Name)
			// fieldMap: 哈希表，用于 O(1) 时间复杂度的字段查找
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

// RecordValues 从目标结构体实例中提取所有字段的值，按照 Schema 定义的字段顺序返回。
// 该方法主要用于将结构体数据转换为数据库操作所需的值列表。
//
// 参数:
//   - dest: 目标结构体实例，可以是指针或值类型，函数内部会自动解引用
//
// 返回值:
//   - []interface{}: 包含所有字段值的切片，顺序与 Schema.Fields 中定义的字段顺序一致
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
