package output

import "github.com/gookit/color"

type Outputter interface {
	Debug(message interface{})
	DebugBytes(message []byte)
	Error(message interface{})
	Info(message string)
	RenderError(a ...interface{}) string
	RenderInfo(a ...interface{}) string
	RenderSuccess(a ...interface{}) string
	RenderWarn(a ...interface{}) string
	Success(message string)
	Warn(message string)
}

var ErrorColor = color.FgRed
var DebugColor = color.FgGray
var InfoColor = color.FgCyan
var PassColor = color.FgCyan
var SuccessColor = color.FgCyan
var WarnColor = color.FgYellow
