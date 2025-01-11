package utils

import (
	"TalkHive/global"
	"TalkHive/models"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// HashPassword 哈希化密码
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
}

// CheckPassword 检查密码
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT 生成JWT
func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	signedToken, err := token.SignedString([]byte("secret"))
	return "Bearer " + signedToken, err
}

// ParseJWT 解析JWT
func ParseJWT(tokenString string) (string, error) {
	// 正确的Token带有Bearer前缀
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected Signing Method")
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("username claim is not a string")
		}
		return username, nil
	}

	return "", err
}

// RandomCode 产生6位随机验证码
func RandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}

// ValidatePhone 格式校验：手机号必须是以 1 开头的 11 位数字
func ValidatePhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^1[0-9]{10}$`)
	return phoneRegex.MatchString(phone)
}

// ValidateEmail 验证电子邮件格式
func ValidateEmail(email string) bool {
	// 定义电子邮件格式的正则表达式
	emailRegex := `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// CheckEmailRegistered 检查邮箱是否已经注册
func CheckEmailRegistered(email string) bool {
	var user models.AccountInfo
	err := global.Db.Where("email = ?", email).First(&user).Error
	return err == nil
}

// SendSms 发送短信验证码
func SendSms(email string, code string) error {
	// QQ邮箱的SMTP服务器地址
	smtpHost := "smtp.qq.com"
	smtpPort := 587 // QQ邮箱SMTP端口

	// 发件人信息
	fromEmail := "15886421754@qq.com" // 发送邮箱
	authCode := "nnqfhioshmxndehi"    // QQ邮箱的SMTP授权码

	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)  // 设置发件人
	m.SetHeader("To", email)        // 设置收件人（即目标邮箱）
	m.SetHeader("Subject", "您的验证码") // 邮件主题
	m.SetBody("text/html", fmt.Sprintf("<h2>您的验证码是：<strong>%s</strong></h2><p>该验证码有效期为5分钟，请尽快使用。</p>", code))

	// 设置SMTP客户端
	d := gomail.NewDialer(smtpHost, smtpPort, fromEmail, authCode)

	// 发送邮件
	err := d.DialAndSend(m)
	if err != nil {
		log.Printf("邮件发送失败: %v", err)
		return fmt.Errorf("邮件发送失败: %v", err)
	}

	log.Printf("验证码邮件已发送至: %s", email)
	return nil
}

// GetAvatarPath 获取图片路径
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

	// 保存在数据库的路径
	return avatarPath, nil
}

// ExtractBase64Data 被调用
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
