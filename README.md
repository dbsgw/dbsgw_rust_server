# dbsgw_rust_server

#### 介绍
{**以下是 Gitee 平台说明，您可以替换此简介**
Gitee 是 OSCHINA 推出的基于 Git 的代码托管平台（同时支持 SVN）。专为开发者提供稳定、高效、安全的云端软件开发协作平台
无论是个人、团队、或是企业，都能够用 Gitee 实现代码托管、项目管理、协作开发。企业项目请看 [https://gitee.com/enterprises](https://gitee.com/enterprises)}

#### 软件架构
软件架构说明

## api说明
1. `controllers/v1` 是api路径，方便扩展二次开发 v1,v2.....
2. `admin` 是权限的后台文件api，用户管理项目的
3. `user` 是用户的权限的api，用户个人中心，管理文章的
4. `web` 是前台展示的api，无权限的api文件

#### 安装教程

1.  xxxx
2.  xxxx
3.  xxxx

#### 使用说明

1.  打包 `bee pack -be GOOS=linux`
2.  查询运行的进程 `ps -aux | grep dbs***`
3.  长期运行go程序 `nohup dbsgw****`

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Requesta


## 设计图

#### 用户登录

邮箱注册，第三方授权登录，
注册，邮箱注册直接生成user_base（自己生成的）和user_auth，第三方授权注册也生成user_base（信息资料用第三方的）和user_auth
邮箱注册后授权第三方，第三方授权后绑定邮箱


# 第三方授权绑定有bug
1. 不能改变 回调地址 回调地址不能传参
2. 更新了用户资料  没有更新 redis缓存  导致 新发帖的 帖子信息还是之前资料的