package common

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/store_server/logger"
	"gopkg.in/yaml.v2"
)

func CloseFile(file *os.File) {
	if file == nil {
		return
	}
	err := file.Close()
	if err != nil {
		logger.Entry().Errorf("close file[%s] err: %v", file.Name(), err)
	}
	return
}

func GetSize(path string) int64 { //获取单个文件大小
	fileInfo, err := os.Stat(path)
	if err != nil {
		logger.Entry().Errorf("get file[%s] err: %v", path, err)
		return int64(0)
	}
	fileSize := fileInfo.Size()
	return fileSize
}

func CalcMediaFileDuration(fileurl string) string { //计算媒体文件播放时长
	args := []string{
		"--Inform=General;%Duration%",
		fileurl,
	}
	if output, err := exec.Command("mediainfo", args...).CombinedOutput(); err == nil {
		/*args := []string{
			"-v",
			"quiet",
			"-print_format",
			"compact=print_section=0:nokey=1:escape=csv",
			"-show_entries",
			"format=duration",
			"-i",
			fileurl,
		}
		if output, err := exec.Command("./ffprobe", args...).CombinedOutput(); err == nil {*/
		dur, _ := strconv.ParseFloat(strings.TrimSpace(string(output)), 32)
		return strconv.Itoa(int(dur / 1000)) //unit to second
	} else {
		logger.Entry().Errorf("calculate media file duration err: %v|fileurl: %v", err, fileurl)
	}
	return "0"
}

func CheckFileMime(fileurl, mime string) (bool, string) { //校验文件mime类型
	args := []string{
		"--mime-type",
		fileurl,
	}
	ok := false
	var suffix string
	if output, err := exec.Command("file", args...).CombinedOutput(); err == nil {
		mimeOutput := strings.TrimSpace(strings.Split(string(output), ":")[1])
		if mimeOutput == mime {
			ok = true
			suffix = strings.Split(mimeOutput, "/")[1]
		}
	}
	return ok, suffix
}

func CheckExt(filename string, ext []string) (bool, string) { //校验文件后缀
	fileExt := filepath.Ext(filename)
	for _, v := range ext {
		if v == fileExt {
			return true, fileExt
		}
	}
	return false, fileExt
}

func GetFileNameFromPath(path string) (filename string) { //从路径中获取文件名(去掉后缀)
	ext := filepath.Ext(path)
	filename = filepath.Base(path)
	if len(ext) != 0 {
		filename = strings.TrimSuffix(filename, ext)
	}
	return
}

func Md5SumForStr(data string) string { //对字符串进行MD5哈希
	md5Ctx := md5.New()
	io.WriteString(md5Ctx, data)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func Md5Sum(data string) string { //对字符串进行MD5哈希
	md5Ctx := md5.New()
	io.WriteString(md5Ctx, data)
	return fmt.Sprintf("%x", md5Ctx.Sum(nil))
}

func Md5File(path string) (string, error) { //计算文件md5
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	r := bufio.NewReader(f)
	h := md5.New()
	if _, err = io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func SHA1File(path string) (string, error) { //计算文件sha1
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	r := bufio.NewReader(f)
	h := sha1.New()
	if _, err = io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func CheckDirEmpty(dir string) (bool, []os.FileInfo) { //检查给定目录是否为空, 不为空则返回文件列表
	fis, _ := ioutil.ReadDir(dir)
	if len(fis) == 0 {
		return true, fis
	}
	return false, fis
}

func CheckFileExist(path string) bool { //检查给定路径文件是否存在
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func IsDir(path string) bool { //检查是否为目录
	return isFileOrDir(path, true)
}

func IsFile(path string) bool { //检查是否为文件
	return isFileOrDir(path, false)
}

func isFileOrDir(path string, judgeDir bool) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	isDir := fileInfo.IsDir()
	if judgeDir {
		return isDir
	}
	return !isDir
}

func DeleteDir(dir string) error { //删除指定目录
	if !IsDir(dir) {
		return fmt.Errorf("given path is not a dir.")
	}
	if empty, _ := CheckDirEmpty(dir); !empty { //目录非空
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}
	err := os.Remove(dir) //目录为空
	return err
}

// ParseYamlConfigData 解析二进制格式的yaml数据
func ParseYamlConfigData(data []byte, out interface{}) error {
	err := yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}

// ParseYamlConfigFile 解析yaml文件
func ParseYamlConfigFile(filePath string, out interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return ParseYamlConfigData(data, out)
}
