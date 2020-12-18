# go-12306

#### 介绍
使用golang开发的12306抢票程序

#### 已有功能

- [x] 自动打码
- [x] 自动滑块
- [x] 自动登录
- [x] 自动获取Cookie
- [x] 准点预售捡漏
- [x] 自动提交订单
- [x] 邮件通知
- [ ] 微信通知

#### 依赖库

- 验证码识别使用的是第三方

  ```
  http://littlebigluo.qicp.net:47720/
  ```

- 火狐浏览器

  ```
  http://www.firefox.com.cn/
  ```

- go安装教程

  ```
  https://www.runoob.com/go/go-environment.html
  ```

#### 项目使用说明

- 修改配置文件

  ```
  配置文件格式是yaml不会的百度，邮箱地址先配置发件邮箱
  ```

  

- 启动服务

  ```
  使用的是go mod包管理
  ```

  ```
  1. go mod init go-12306
  2. go mod vendor
  3. go run main.go
  ```

   

  ```
  打包
  go build -a
  ```

  

#### 建议

```
先运行 login.exe 他是获取cookie，之后会每次查是否登录没有会自动弹出获取
```

#### 项目声明：

- 本软件只供学习交流使用，勿作为商业用途，交流群号
  - 781286902
- 进群先看公告！！！

