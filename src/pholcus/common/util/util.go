package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"io"
	r "math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html/charset"

	"github.com/henrylee2cn/pholcus/logs"
)

// JsonpToJson modify jsonp string to json string
// Example: forbar({a:"1",b:2}) to {"a":"1","b":2}
func JsonpToJson(json string) string {
	start := strings.Index(json, "{")
	end := strings.LastIndex(json, "}")
	start1 := strings.Index(json, "[")
	if start1 > 0 && start > start1 {
		start = start1
		end = strings.LastIndex(json, "]")
	}
	if end > start && end != -1 && start != -1 {
		json = json[start : end+1]
	}
	json = strings.Replace(json, "\\'", "", -1)
	regDetail, _ := regexp.Compile("([^\\s\\:\\{\\,\\d\"]+|[a-z][a-z\\d]*)\\s*\\:")
	return regDetail.ReplaceAllString(json, "\"$1\":")
}

// 创建目录
func Mkdir(Path string) {
	p, _ := path.Split(Path)
	if p == "" {
		return
	}
	d, err := os.Stat(p)
	if err != nil || !d.IsDir() {
		if err = os.MkdirAll(p, 0777); err != nil {
			logs.Log.Error("创建路径失败[%v]: %v\n", Path, err)
		}
	}
}

// The GetWDPath gets the work directory path.
func GetWDPath() string {
	wd := os.Getenv("GOPATH")
	if wd == "" {
		panic("GOPATH is not setted in env.")
	}
	return wd
}

// The IsDirExists judges path is directory or not.
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	panic("util isDirExists not reached")
}

// The IsFileExists judges path is file or not.
func IsFileExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return !fi.IsDir()
	}

	panic("util isFileExists not reached")
}

// 遍历并返回指定类型范围的文件名列表
// 默认返回所有文件
func WalkFiles(path string, suffixes ...string) (filelist []string) {
	path, _ = filepath.Abs(path)

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if len(suffixes) == 0 {
			filelist = append(filelist, path)
			return nil
		}
		for _, suffix := range suffixes {
			if strings.HasSuffix(path, suffix) {
				filelist = append(filelist, path)
			}
		}
		return nil
	})

	if err != nil {
		logs.Log.Error("filepath.Walk() returned %v\n", err)
	}

	return
}

// The IsNum judges string is number or not.
func IsNum(a string) bool {
	reg, _ := regexp.Compile("^\\d+$")
	return reg.MatchString(a)
}

// simple xml to string  support utf8
func XML2mapstr(xmldoc string) map[string]string {
	var t xml.Token
	var err error
	inputReader := strings.NewReader(xmldoc)
	decoder := xml.NewDecoder(inputReader)
	decoder.CharsetReader = func(s string, r io.Reader) (io.Reader, error) {
		return charset.NewReader(r, s)
	}
	m := make(map[string]string, 32)
	key := ""
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			key = token.Name.Local
		case xml.CharData:
			content := string([]byte(token))
			m[key] = content
		default:
			// ...
		}
	}

	return m
}

//string to hash
func MakeHash(s string) string {
	const IEEE = 0xedb88320
	var IEEETable = crc32.MakeTable(IEEE)
	hash := fmt.Sprintf("%x", crc32.Checksum([]byte(s), IEEETable))
	return hash
}

func HashString(encode string) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(encode))
	return hash.Sum64()
}

// 制作特征值方法一
func MakeUnique(obj interface{}) string {
	baseString, _ := json.Marshal(obj)
	return strconv.FormatUint(HashString(string(baseString)), 10)
}

// 制作特征值方法二
func MakeMd5(obj interface{}, length int) string {
	if length > 32 {
		length = 32
	}
	h := md5.New()
	baseString, _ := json.Marshal(obj)
	h.Write([]byte(baseString))
	s := hex.EncodeToString(h.Sum(nil))
	return s[:length]
}

// 将对象转为json字符串
func JsonString(obj interface{}) string {
	b, _ := json.Marshal(obj)
	s := fmt.Sprintf("%+v", string(b))
	r := strings.Replace(s, `\u003c`, "<", -1)
	r = strings.Replace(r, `\u003e`, ">", -1)
	return r
}

//检查并打印错误
func CheckErr(err error) {
	if err != nil {
		logs.Log.Error("%v", err)
	}
}

const (
	// Spider中启用Keyword的初始值
	USE_KEYWORD = "\r\t\n"
)

// 将文件名非法字符替换为相似字符
func FileNameReplace(fileName string) (rfn string) {
	// 替换`""`为`“”`
	if strings.Count(fileName, `"`) > 0 {
		var i = 1
	label:
		for k, v := range []byte(fileName) {
			if string(v) != `"` {
				continue
			}
			if i%2 == 1 {
				fileName = string(fileName[:k]) + `“` + string(fileName[k+1:])
			} else {
				fileName = string(fileName[:k]) + `”` + string(fileName[k+1:])
			}
			i++
			goto label
		}
	}

	rfn = strings.Replace(fileName, `:`, `：`, -1)
	rfn = strings.Replace(rfn, `*`, `ж`, -1)
	// rfn = strings.Replace(rfn, `*`, `×`, -1)
	rfn = strings.Replace(rfn, `<`, `＜`, -1)
	rfn = strings.Replace(rfn, `>`, `＞`, -1)
	rfn = strings.Replace(rfn, `?`, `？`, -1)
	rfn = strings.Replace(rfn, `/`, `／`, -1)
	rfn = strings.Replace(rfn, `|`, `∣`, -1)
	rfn = strings.Replace(rfn, `\`, `╲`, -1)
	rfn = strings.Replace(rfn, USE_KEYWORD, ``, -1)
	return
}

// 将Excel工作表名中非法字符替换为下划线
func ExcelSheetNameReplace(fileName string) (rfn string) {
	rfn = strings.Replace(fileName, `:`, `_`, -1)
	rfn = strings.Replace(rfn, `：`, `_`, -1)
	rfn = strings.Replace(rfn, `*`, `_`, -1)
	rfn = strings.Replace(rfn, `?`, `_`, -1)
	rfn = strings.Replace(rfn, `？`, `_`, -1)
	rfn = strings.Replace(rfn, `/`, `_`, -1)
	rfn = strings.Replace(rfn, `／`, `_`, -1)
	rfn = strings.Replace(rfn, `\`, `_`, -1)
	rfn = strings.Replace(rfn, `╲`, `_`, -1)
	rfn = strings.Replace(rfn, `]`, `_`, -1)
	rfn = strings.Replace(rfn, `[`, `_`, -1)
	rfn = strings.Replace(rfn, USE_KEYWORD, ``, -1)

	return
}

func Atoa(str interface{}) string {
	if str == nil {
		return ""
	}
	return strings.Trim(str.(string), " ")
}

func Atoi(str interface{}) int {
	if str == nil {
		return 0
	}
	i, _ := strconv.Atoi(strings.Trim(str.(string), " "))
	return i
}

func Atoui(str interface{}) uint {
	if str == nil {
		return 0
	}
	u, _ := strconv.Atoi(strings.Trim(str.(string), " "))
	return uint(u)
}

// RandomCreateBytes generate random []byte by specify chars.
func RandomCreateBytes(n int, alphabets ...byte) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return bytes
}