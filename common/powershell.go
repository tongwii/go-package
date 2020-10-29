package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
)

/*
   url参数处理
*/
func ParseArgs(params map[string][]string) string {
	var argsBuffer bytes.Buffer

	// 遍历query参数
	for key, value := range params {
		argsBuffer.WriteString(" " + key)
		argsBuffer.WriteString(" " + value[0])
	}
	// 定义错误处理
	//argsBuffer.WriteString(" -ErrorAction Stop")
	return argsBuffer.String()
}

/*
   实际执行powershell
*/
func ExecScript(script string, show bool) *Content {
	script = string(strings.TrimSpace(script))
	ret := &Content{Status: http.StatusOK}
	var commandBuffer bytes.Buffer
	// 组合指令开始
	// 添加错误处理
	//commandBuffer.WriteString("$ErrorView='CategoryView';")
	//commandBuffer.WriteString("try{")
	commandBuffer.WriteString("Trap {\"Error: $($_.Exception.Message)\"; Continue} ") // 捕获异常
	commandBuffer.WriteString(script)
	//commandBuffer.WriteString("}Catch{Write-Host \"Error\" $_.Exception.ToString()}")
	//Logger.Debug(commandBuffer.String())
	if show { // 有些脚本始终不输出
		Logger.Trace(commandBuffer.String()) // 日志中 trace 模式才显示
	}
	/**
	 * 执行powershell
	 * NoLogo 启动时隐藏版权标志。
	 * NonInteractive 不向用户显示交互式提示。
	 */
	cmd := exec.Command("powershell.exe", "-NoLogo", "-NonInteractive", "-Command", "&{", commandBuffer.String(), "}")
	//out, err := cmd.Output()          // 与CombinedOutput的区别是如果命令出错，这里不会有输出
	out, err := cmd.CombinedOutput() // 标准输出，与在控制台看到的一致

	// 这里的错误是指exec是否执行成功 而不是powershell执行的结果
	if err != nil {
		//log.WithFields(log.Fields{
		//    "command": string(sc),
		//}).Error(err.Error())
		Logger.Errorf("exec error: %s", err.Error())
	}

	// 无返回
	if len(out) == 0 {
		ret.Message = "exec success"
		ret.Status = http.StatusOK
		return ret
	}

	// 转码到UTF-8
	out, _ = ConvertToUTF8(out)
	// 转换为json
	var data interface{}
	if err := json.Unmarshal(out, &data); err != nil {
		if strings.Contains(string(out), "Error") || strings.Contains(string(out), "error") {
			ret.Message = CleanString(string(out[7:]))
			ret.Status = http.StatusBadRequest
			Logger.Errorf("powershell error: %s", ret.Message)
		} else {
			ret.Data = CleanString(string(out))
			ret.Status = http.StatusOK
		}
		return ret
	}
	ret.Data = data
	return ret
}

/**
  清理字符串中无用符号
*/
func CleanString(src string) string {
	src = strings.Replace(src, "\r", "", -1)
	src = strings.Replace(src, "\n", "", -1)
	//src = strings.Replace(src, " ", "", -1)
	src = strings.Replace(src, "{", "", -1)
	src = strings.Replace(src, "}", "", -1)
	//src = strings.Replace(src, "[", "", -1)
	//src = strings.Replace(src, "]", "", -1)
	src = strings.Replace(src, "-", "", -1)
	src = strings.Replace(src, "_", "", -1)
	return src
}
