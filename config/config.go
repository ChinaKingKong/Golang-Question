package config

import (
	"golang-question/errorx"
	"reflect"
	"sync"
)

type Manager[T any] interface {
	Get() T
	Update(T) errorx.Error
	OnChange(func(T)) (cancel func())
	Watch() Manager[T]
	InitData(T) Manager[T]
}

type localManager[T any] struct {
	data     T
	mu       sync.RWMutex
	callback func(T)
}

// Get 从本地管理器中获取当前存储的数据。
// 返回存储的数据。
// 需要注意，这是一个并发安全的方法，使用了读写锁来保证数据的一致性。
func (m *localManager[T]) Get() T {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data
}

// Update 用于更新 localManager 中的数据。
// 参数 newData 是新的数据。
// 如果提供了回调函数 callback，则在更新数据后调用该函数，并将 newData 作为参数传入。
// 如果操作成功，返回 nil；否则返回 errorx.Error 类型的错误。
func (m *localManager[T]) Update(newData T) errorx.Error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = newData
	if m.callback != nil {
		m.callback(newData)
	}
	return nil
}

// OnChange 为localManager类型的方法，它接收一个类型为T的回调函数作为参数。
// 当localManager管理的数据发生变化时，会自动调用该回调函数。
//
// 参数：
//     callback: 当数据发生变化时调用的回调函数，接收一个类型为T的参数。
//
// 返回值：
//     返回一个取消函数，调用该函数可以取消之前设置的回调函数。
//
// 注意：
//     当调用取消函数后，即使数据发生变化，之前设置的回调函数也不会再被调用。
func (m *localManager[T]) OnChange(callback func(T)) (cancel func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callback = callback
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.callback = nil
	}
}

// Watch 返回当前localManager对象本身，实现Manager接口
func (m *localManager[T]) Watch() Manager[T] { 
	return m
}

// isZeroValue 函数用于判断给定的值是否是其类型的零值。
//
// 参数：
//   - value: 要判断的值，类型为泛型T。
//
// 返回值：
//   - bool: 如果value是其类型的零值，则返回true；否则返回false。
func isZeroValue[T any](value T) bool {
	return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
}

// InitData 使用初始数据初始化localManager的数据
// 如果localManager的数据是零值，则将localManager的数据设置为initialData
// 返回初始化后的localManager实例
func (m *localManager[T]) InitData(initialData T) Manager[T] {
	m.mu.Lock()
	defer m.mu.Unlock()
	if isZeroValue(m.data) {
		m.data = initialData
	}
	return m
}

// Local 返回一个本地管理器的实例
// 参数：
// - T: 泛型类型，代表管理器将管理的元素类型
//
// 返回值：
// - Manager[T]: 返回本地管理器的实例
func Local[T any]() Manager[T] {
	return &localManager[T]{}
}

// Etcd implementation remains as a TODO
func Etcd[T any]() Manager[T] {
	//TODO: implement
	return nil
}
