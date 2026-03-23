# go-imgbed

golang实现的极度轻量图床，可以在64M/5G的低配服务器上运行。单文件无依赖运行

## 使用方法

1.配置config.json

```json
{
	"save_file": "/mnt/data",
	"user": "admin",
	"pass": "123456789",
	"upload": "/upload",
	"listen": "0.0.0.0:80"
}
```

2.下载go-imgbed程序

3.启动程序

```bash
./go-imgbed config.json
```

4.浏览器访问/upload，输入上述设置的用户名与密码即可

<br>

---

## 配置说明

```json
{
	"save_file": "/mnt/data", //图片存储目录，确保有访问权限
	"user": "admin", //用户名
	"pass": "123456789", //密码
	"upload": "/upload", //上传目录，默认/upload。如果更改为"/upload1"。那么浏览器访问"/upload1"
	"listen": "0.0.0.0:443", //监听端口
	"crt": "/mnt/tls.crt", //证书crt，如果此项不为空，启动tls
	"key": "/mnt/tls.key" //证书key，如果此项不为空，启动tls
}
```


### 注意，如果是公网访问，务必启用tls。也可go-imgbed以http启动，上层反代实现https

<br>

---

### 其他事项：

- 上传目录，默认的"/upload"可以改一改，这样可以降低被扫到的概率。
- 密码至少8位以上
- 正式环境千万不要使用admin:123456789这种弱密码