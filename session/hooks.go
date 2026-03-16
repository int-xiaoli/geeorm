package session

import (
	"geeorm/log"
	"reflect"
)

// Hooks constants
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// CallMethod calls the registered hooks
// CallMethod 调用指定对象上的方法，该方法接收 Session 作为参数。
// 如果 value 不为 nil，则在 value 上查找并调用方法；否则在 s.RefTable().Model 上查找并调用方法。
// 若方法返回 error 类型的值且不为 nil，则记录错误日志。
func (s *Session) CallMethod(method string, value interface{}) {
	// 首先尝试从模型实例上获取指定名称的方法
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	// 如果传入了 value 参数，则优先使用 value 上的方法
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	// 准备方法调用的参数，将当前 Session 作为参数传入
	param := []reflect.Value{reflect.ValueOf(s)}
	// 检查方法是否有效（存在）
	if fm.IsValid() {
		// 调用方法并获取返回值
		if v := fm.Call(param); len(v) > 0 {
			// 如果第一个返回值是 error 类型，检查是否为非 nil 错误
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}
