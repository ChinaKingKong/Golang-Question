package errorx

import (
	"fmt"
	"runtime"
)

type Error interface {
	error
	fmt.Formatter
	Unwrap() error
	Cause() error
	Code() int
	Type() ErrType
	Stack() Stack
}

type ErrType string

const (
	ErrTypeNotFound   ErrType = "not_found"
	ErrTypeTimeout    ErrType = "timeout"
	ErrTypeInvalid    ErrType = "invalid"
	ErrTypeConflict   ErrType = "conflict"
	ErrTypePermission ErrType = "permission_denied"
	ErrTypeInternal   ErrType = "internal_error"
	ErrTypeUnavailable ErrType = "unavailable"
	// 可以在此处添加更多错误类型
)

type Frame struct {
	Name string
	File string
	Line int
}

type Stack []Frame

type customError struct {
	msg   string
	code  int
	errType ErrType
	stack Stack
	cause error
}

func (e *customError) Error() string {
	return e.msg
}

func (e *customError) Unwrap() error {
	return e.cause
}

func (e *customError) Cause() error {
	return e.cause
}

func (e *customError) Code() int {
	return e.code
}

func (e *customError) Type() ErrType {
	return e.errType
}

func (e *customError) Stack() Stack {
	return e.stack
}

func (e *customError) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "%s", e.Error())
}

// captureStack 函数用于捕获当前调用栈信息并返回为一个 Stack 类型的切片
// 返回值是一个包含调用栈信息的 Stack 切片
//
// 该函数通过调用 runtime.Caller 遍历调用栈，跳过前两个栈帧（通常是 captureStack 函数本身和它的调用者）
// 对于每个栈帧，函数获取其对应的程序计数器（PC）、文件名、行号和函数信息
// 然后将这些信息封装成一个 Frame 结构体，并添加到返回的 Stack 切片中
// 当遇到无法获取栈帧信息的情况时，函数结束遍历并返回当前已捕获的栈信息
func captureStack() Stack {
	stack := Stack{}
	for i := 2; ; i++ { // Skip the first two frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		stack = append(stack, Frame{
			Name: fn.Name(),
			File: file,
			Line: line,
		})
	}
	return stack
}

// Wrap 将给定的错误包装为 Error 类型，并返回一个新的 Error 实例。
// 如果输入的 error 为 nil，则返回 nil。
// 如果输入的 error 不为 nil，则将其错误信息和堆栈信息包装在 customError 结构体中，并返回该结构体指针。
func Wrap(err error) Error {
	if err == nil {
		return nil
	}
	return &customError{
		msg:   err.Error(),
		cause: err,
		stack: captureStack(),
	}
}

// New 创建一个新的自定义错误对象
//
// 参数:
//     msg: 错误信息字符串
//
// 返回值:
//     返回一个实现了 error 接口的自定义错误对象
func New(msg string) Error {
	return &customError{
		msg:   msg,
		stack: captureStack(),
	}
}

// C 创建一个自定义错误对象
// code：错误码
// msg：错误信息
// 返回值：返回创建的自定义错误对象
func C(code int, msg string) Error {
	return &customError{
		msg:   msg,
		code:  code,
		stack: captureStack(),
	}
}

// Cf 函数用于创建一个自定义错误对象。
//
// 参数:
// - code: int 类型，错误码。
// - format: string 类型，格式化字符串，用于构建错误信息。
// - args: ...interface{} 类型，可变参数列表，用于格式化错误信息。
//
// 返回值:
// - Error: 返回一个 Error 类型的自定义错误对象。
func Cf(code int, format string, args ...interface{}) Error {
	return &customError{
		msg:   fmt.Sprintf(format, args...),
		code:  code,
		stack: captureStack(),
	}
}
