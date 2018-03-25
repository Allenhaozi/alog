// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"fmt"
	"io"
	"os"

	"github.com/Allenhaozi/alog/colors"
	"github.com/Allenhaozi/alog/writers"
)

var consoleOutputMap = map[string]*os.File{
	"stderr": os.Stderr,
	"stdout": os.Stdout,
	"stdin":  os.Stdin,
}

var consoleColorMap = map[string]colors.Color{
	"default": colors.Default,
	"black":   colors.Black,
	"red":     colors.Red,
	"green":   colors.Green,
	"yellow":  colors.Yellow,
	"blue":    colors.Blue,
	"magenta": colors.Magenta,
	"cyan":    colors.Cyan,
	"white":   colors.White,
}

// Console 是 writers.Console 的初始化函数
func Console(args map[string]string) (io.Writer, error) {
	outputIndex, found := args["output"]
	if !found {
		outputIndex = "stderr"
	}

	output, found := consoleOutputMap[outputIndex]
	if !found {
		return nil, fmt.Errorf("[%v]不是一个有效的控制台输出项", outputIndex)
	}

	fcIndex, found := args["foreground"]
	if !found { // 默认用红色前景色
		fcIndex = "red"
	}
	fc, found := consoleColorMap[fcIndex]
	if !found {
		return nil, fmt.Errorf("无效的前景色[%v]", fcIndex)
	}

	bcIndex, found := args["background"]
	if !found {
		bcIndex = "default"
	}
	bc, found := consoleColorMap[bcIndex]
	if !found {
		return nil, fmt.Errorf("无效的背景色[%v]", bcIndex)
	}

	return writers.NewConsole(output, fc, bc), nil
}
