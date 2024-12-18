package utils

import (
	"TalkHive/global"
	"TalkHive/models"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
	"regexp"
	"time"
)

// HashPassword 哈希化密码
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
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

// CheckPassword 检查密码
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

// SendSms 发送短信验证码（授权码: nnqfhioshmxndehi）
func SendSms(email string, code string) interface{} {
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
