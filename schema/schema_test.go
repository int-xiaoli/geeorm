package schema

import (
	"geeorm/dialect"
	"testing"
)

// schema_test.go
type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&User{}, TestDial)

	// 打印 schema 的详细信息
	t.Logf("Schema Name: %s", schema.Name)
	t.Logf("Schema Model: %T", schema.Model)
	t.Logf("Field Count: %d", len(schema.Fields))
	t.Logf("Field Names: %v", schema.FieldNames)

	for i, field := range schema.Fields {
		t.Logf("Field %d: Name=%s, Type=%s, Tag=%s",
			i, field.Name, field.Type, field.Tag)
	}

	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}
