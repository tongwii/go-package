package common

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

/**
显示banner
*/
func Showflag() {
	if !HiddenBanner {
		fmt.Printf(BANNER, Product, Version, Build, Website)
	}

	// 显示端口信息
	fmt.Printf("%s started on %s\n", Product, Port)
}

/*
获取程序运行路径
*/
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	//return strings.Replace(dir, "\\", "/", -1)
	return dir
}

/**
返回程序文件名
*/
func GetExeName() string {
	fullFilename := os.Args[0]
	var filenameWithSuffix string
	filenameWithSuffix = filepath.Base(fullFilename) // 获取文件名带后缀
	var fileSuffix string
	fileSuffix = filepath.Ext(filenameWithSuffix) // 获取文件后缀
	var filenameOnly string
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix) // 获取文件名
	return filenameOnly
}

/**
如果返回的错误为nil,说明文件或文件夹存在
如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
如果返回的错误为其它类型,则不确定是否在存在
*/
func IsPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
ip转byte
*/
func StringIPToByte(ipstring string) []byte {
	ipaddbyte := []byte{0, 0, 0, 0}
	ipSegs := strings.Split(ipstring, ".")
	for i, ipSeg := range ipSegs {
		//fmt.Printf(" %d ipSegs: %s\n",i,ipSeg);
		ipInt, _ := strconv.Atoi(ipSeg)
		ipaddbyte[i] = byte(ipInt)
	}
	return ipaddbyte
}

/**
转码到UTF-8
*/
func ConvertToUTF8(src []byte) (dst []byte, err error) {
	dst, err = ioutil.ReadAll(transform.NewReader(bytes.NewReader(src), simplifiedchinese.GB18030.NewDecoder()))
	if err == nil {
		return dst, nil
	}
	return nil, err
}

/**
MD5
*/
func EncryptionMd5(password string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password))
	cipherStr := md5Ctx.Sum(nil)
	passString := hex.EncodeToString(cipherStr)
	return passString
}

/**
字符串是否为空
*/
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

/**
字符串是否纯数字
*/
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

/**
复制文件
src 源
dst 目的
buffer 缓冲区大小
*/
func CopyFile(src, dst string, buffer int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return errors.New(fmt.Sprintf("%s is not a regular file.", src))
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = source.Close() }()

	_, err = os.Stat(dst)
	if err == nil {
		err = os.Remove(dst)
		if err != nil {
			return errors.New(fmt.Sprintf("File %s already exists.", dst))
		}
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = destination.Close() }()

	if buffer <= 0 {
		buffer = 1000000
	}
	buf := make([]byte, buffer)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

// 删除下载文件及临时文件
func DeleteFile(parent, path string) error {
	if IsEmpty(path) {
		Logger.Debugf("delete file path empty")
		return nil
	}
	filePath := filepath.Join(parent, path)
	err := os.Remove(filePath) // 删除文件

	tempPath := fmt.Sprintf("%s.temp", filePath)
	errs := os.Remove(tempPath) // 删除临时文件

	if err != nil && errs != nil {
		return err
	}
	return nil
}

// 批量删除文件
func DeleteFiles(filePaths ...string) error {
	var err error
	for _, filePath := range filePaths {
		if IsEmpty(filePath) {
			Logger.Debugf("delete file path empty")
			continue
		}
		err = os.Remove(filePath) // 删除文件
		if err != nil {
			Logger.Error(err.Error())
		}
	}
	return err
}

/*
获取文件大小
src 文件路径
*/
func GetFileSize(src string) (int64, error) {
	fi, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	// get the size
	return fi.Size(), nil
}

/**
检查IP是否在一个范围内
start 开始IP 192.168.0.1
end 结束IP 192.168.0.254
ip 需检测IP 192.168.200
*/
func IsIPInRange(start, end, ip string) bool {
	ip1 := net.ParseIP(start)
	ip2 := net.ParseIP(end)
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false
	}
	if bytes.Compare(trial, ip1) >= 0 && bytes.Compare(trial, ip2) <= 0 {
		return true
	}
	return false
}

/**
查询数组中是否已包含字符串
*/
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// 判断系统位数，64位返回64，32位返回32
func GetBits() (bit int) {
	bit = 32 << (^uint(0) >> 63)
	return
}

// 批量删除map的key
func DeleteMap(value *map[string]interface{}, attrs ...string) {
	for _, v := range attrs {
		delete(*value, v)
	}
}

// 三元运算
// 例子：max := If(a > b, a, b).(int)
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
