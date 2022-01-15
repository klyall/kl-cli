package output

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type SStdOut struct {
	Out io.Writer
}

func (s SStdOut) Debug(message interface{}) {
	fmt.Fprintf(s.Out, "%s\n", DebugColor.Render(message))
}

func (s SStdOut) DebugBytes(content []byte) {
	r := bytes.NewReader(content)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if line != "" {
			s.Debug(scanner.Text())
		}
	}
}

func (s SStdOut) Error(message interface{}) {
	s.printMessage(ErrorColor.Render("ERROR"), message)
}

func (s SStdOut) Info(message string) {
	s.printMessage(InfoColor.Render("INFO"), message)
}

func (s SStdOut) Success(message string) {
	s.printMessage(SuccessColor.Render("SUCCESS"), message)
}

func (s SStdOut) Warn(message string) {
	s.printMessage(WarnColor.Render("WARN"), message)
}

func (s SStdOut) RenderError(a ...interface{}) string {
	return ErrorColor.Render(a)
}

func (s SStdOut) RenderInfo(a ...interface{}) string {
	return InfoColor.Render(a)
}

func (s SStdOut) RenderSuccess(a ...interface{}) string {
	return SuccessColor.Render(a)
}

func (s SStdOut) RenderWarn(a ...interface{}) string {
	return WarnColor.Render(a)
}

func (s SStdOut) printMessage(status, message interface{}) {
	fmt.Fprintf(s.Out, "%-16s %s\n", status, message)
}
