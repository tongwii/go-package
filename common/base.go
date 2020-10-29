/**
在编译参数中加入，位置：go tool arguments
windows:
go build -ldflags "-X ccms/common.Build=%date:~0,10%.%time:~0,8%"
linux:
go build -ldflags "-X ccms/common.Build=`date -u +%Y%m%d.%H%M%S`"
  用于修改build
*/
package common

const (
	BANNER = `
___________                   __      __.___.___
\__    ___/___   ____    ____/  \    /  \   |   |
  |    | /  _ \ /    \  / ___\   \/\/   /   |   |
  |    |(  <_> )   |  \/ /_/  >        /|   |   |
  |____| \____/|___|  /\___  / \__/\  / |___|___|
                    \//_____/       \/
%s %s Build:%s
%s
`
)

var (
	h, v            bool      // 用于传入参数 -h help -v version
	Build           = ""      // 编译时间
	GitBranch       = ""      // git 分支
	GitCommit       = ""      // git 提交
	CompilerVersion = ""      // go 编译版本
	Version         = "2.0.0" // 软件版本
	Website         = "http://www.tongwii.com"
	Product         = ""     // 产品名称
	HiddenBanner    = false  // 是否显示banner
	Port            = ":"    // 端口号
	LogHoldDays     = 30     // 日志保留天数
	LogPath         = "logs" // 日志路径
	ServiceFlag     = ""     // 服务控制
	Debug           = false  // 调试模式
	Trace           = false  // 日志模式
)
