package emailverify

import (
	"ChildrenMath/pkg/settings"
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
)

const ExpireTime = 5

var (
	ErrEqual = errors.New("验证码错误")
	ErrSend  = errors.New("该邮箱5分钟内已获取过验证码")
)

// MailboxConf 邮箱配置
type MailboxConf struct {
	// 邮件标题
	Title string
	// 邮件内容
	Body string
	// 收件人列表
	RecipientList []string
	// 发件人账号
	Sender string
	// 发件人密码，QQ邮箱这里配置授权码
	SPassword string
	// SMTP 服务器地址， QQ邮箱是smtp.qq.com
	SMTPAddr string
	// SMTP端口 QQ邮箱是25
	SMTPPort int
}

func SendRegisterEmail(recipient, vCode string) error {
	var mailConf MailboxConf
	mailConf.Title = settings.EmailTitle + "注册认证"

	//这里支持群发，只需填写多个人的邮箱即可，我这里发送人使用的是QQ邮箱，所以接收人也必须都要是
	//QQ邮箱
	mailConf.RecipientList = []string{recipient}
	mailConf.Sender = settings.EmailSender

	//这里QQ邮箱要填写授权码，网易邮箱则直接填写自己的邮箱密码，授权码获得方法在下面
	mailConf.SPassword = settings.EmailSPassword

	//下面是官方邮箱提供的SMTP服务地址和端口
	// QQ邮箱：SMTP服务器地址：smtp.qq.com（端口：587）
	// 雅虎邮箱: SMTP服务器地址：smtp.yahoo.com（端口：587）
	// 163邮箱：SMTP服务器地址：smtp.163.com（端口：25）
	// 126邮箱: SMTP服务器地址：smtp.126.com（端口：25）
	// 新浪邮箱: SMTP服务器地址：smtp.sina.com（端口：25）

	mailConf.SMTPAddr = settings.EmailSMTPAddr
	mailConf.SMTPPort = settings.EmailSMTPPort

	//发送的内容
	html := fmt.Sprintf(`<div>
        <div>
            尊敬的用户，您好！
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>你本次的验证码为%s,为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
        </div>
        <div>
            <p>此邮箱为系统邮箱，请勿回复。</p>
        </div>
    </div>`, vCode)

	m := gomail.NewMessage()

	// 第三个参数是我们发送者的名称，但是如果对方有发送者的好友，优先显示对方好友备注名
	//m.SetHeader(`From`, mailConf.Sender, )
	m.SetHeader(`From`, m.FormatAddress(mailConf.Sender, "少儿数学计算竞技平台"))
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件

	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}

func SendForgetPasswordEmail(recipient, vCode string) error {
	var mailConf MailboxConf
	mailConf.Title = settings.EmailTitle + "重置密码"

	//这里支持群发，只需填写多个人的邮箱即可，我这里发送人使用的是QQ邮箱，所以接收人也必须都要是
	//QQ邮箱
	mailConf.RecipientList = []string{recipient}
	mailConf.Sender = settings.EmailSender

	//这里QQ邮箱要填写授权码，网易邮箱则直接填写自己的邮箱密码，授权码获得方法在下面
	mailConf.SPassword = settings.EmailSPassword

	//下面是官方邮箱提供的SMTP服务地址和端口
	// QQ邮箱：SMTP服务器地址：smtp.qq.com（端口：587）
	// 雅虎邮箱: SMTP服务器地址：smtp.yahoo.com（端口：587）
	// 163邮箱：SMTP服务器地址：smtp.163.com（端口：25）
	// 126邮箱: SMTP服务器地址：smtp.126.com（端口：25）
	// 新浪邮箱: SMTP服务器地址：smtp.sina.com（端口：25）

	mailConf.SMTPAddr = settings.EmailSMTPAddr
	mailConf.SMTPPort = settings.EmailSMTPPort

	//发送的内容
	html := fmt.Sprintf(`<div>
        <div>
            尊敬的用户，您好！
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>您正在通过“忘记密码”选项进行密码重置，重置的验证码为%s,为了保证账号安全，验证码有效期为5分钟。</p>
			<p>如果这不是您的操作，切勿向他人泄露并及时登录您的账户查看，感谢您的理解与使用。</p>
        </div>
        <div>
            <p>此邮箱为系统邮箱，请勿回复。</p>
        </div>
    </div>`, vCode)

	m := gomail.NewMessage()

	// 第三个参数是我们发送者的名称，但是如果对方有发送者的好友，优先显示对方好友备注名
	//m.SetHeader(`From`, mailConf.Sender, )
	m.SetHeader(`From`, m.FormatAddress(mailConf.Sender, "少儿数学计算竞技平台"))
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件

	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}
