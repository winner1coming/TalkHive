package global

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Db         *gorm.DB
	RedisDB    *redis.Client
	SmsCodeKey = "sms_code:" // 短信验证码的 Redis 键前缀
)

func GetAvatarPath(base64Str string, userID string, Path string) (string, error) {
	// 提取 Base64 数据和文件类型
	base64Data, fileType, err := ExtractBase64Data(base64Str)
	if err != nil {
		return "", fmt.Errorf("Base64 数据解析失败: %v", err)
	}

	// 解码 Base64 数据
	avatarData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("Base64 解码失败: %v", err)
	}

	// 定义保存路径
	avatarDir := "D:/TalkHive/" + Path
	avatarPath := filepath.Join(avatarDir, fmt.Sprintf("%s.%s", userID, fileType)) // 文件名为 userID.<fileType>

	// 检查目录是否存在，如果不存在则创建
	if _, err := os.Stat(avatarDir); os.IsNotExist(err) {
		if err := os.MkdirAll(avatarDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("创建目录失败: %v", err)
		}
	}

	// 保存图片到指定路径
	if err := ioutil.WriteFile(avatarPath, avatarData, 0644); err != nil {
		return "", fmt.Errorf("保存图片失败: %v", err)
	}

	// 返回图片路径
	return avatarPath, nil
}

func ExtractBase64Data(base64Str string) (string, string, error) {
	parts := strings.Split(base64Str, ",")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("无效的 Base64 数据")
	}

	// 提取文件类型
	header := parts[0]
	fileType := ""
	if strings.HasPrefix(header, "data:image/") && strings.Contains(header, ";base64") {
		// 从头部提取文件类型，例如 "data:image/png;base64" -> "png"
		fileType = strings.TrimSuffix(strings.TrimPrefix(header, "data:image/"), ";base64")
	} else {
		return "", "", fmt.Errorf("不支持的文件类型或无效的 Base64 头部")
	}

	// 返回纯 Base64 数据和文件类型
	return parts[1], fileType, nil
}

// GetFileContentAndType 读取文件内容并获取文件类型
func GetFileContentAndType(filePath string) (string, string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	fileContents, err := io.ReadAll(file)
	if err != nil {
		return "", "", fmt.Errorf("读取文件内容失败: %v", err)
	}

	// 将文件内容编码为 Base64
	fileBase64 := base64.StdEncoding.EncodeToString(fileContents)

	// 获取文件类型
	fileType := mime.TypeByExtension(filepath.Ext(filePath))
	if fileType == "" {
		fileType = "image/jpg" // 默认类型
	}

	// 返回文件内容和文件类型
	return fileBase64, fileType, nil
}
