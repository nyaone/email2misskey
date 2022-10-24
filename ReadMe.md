# eMail to Misskey

## 流程

1. 作为 SMTP 服务器接收 25 端口入站邮件
2. 将邮件压缩成 zip 格式，上传到 Misskey
3. 通过 Misskey 发送私信给目标用户

## 部署

请先确认您目标服务器上的 TCP 协议入站 25 端口开启，并先安装好 docker 与 docker-compose 。

### 获得一个 API 访问令牌

1. 前往您实例的 `设置 - 其他设置 - API` (/settings/api) 
2. 单击 `生成访问令牌` ，请启用 `管理网盘文件` 和 `撰写或删除消息` 两项，其他的保持禁用状态
3. 给它取一个名（例如 eMail to Misskey ），单击右上角的钩子确认，在弹出窗口中复制您的 API 令牌

### 部署服务

1. 新建一个目录并进入
2. 复制 `docker-compose.yml` ，放置在您的目录中
3. 复制 `config.json.example` ，重命名为 `config.json` 并放置在您的目录中
4. 编辑 `config.json`
   1. 将您的 Misskey 实例域名（例如 nya.one）写入 `misskey.instance` 字段
   2. 将上一步得到的 API 令牌写入 `misskey.token` 字段
   3. 如果您需要设置一个目录用于管理邮件（推荐），请将目录 ID 写入 `misskey.folderId` 字段
   4. 在 `email.host` 中设置您需要启用的收件地址域名（可以为多个），请注意非该域名的邮件都会被拒收
5. 使用 `docker-compose pull` 拉取容器镜像
6. 使用 `docker-compose up -d` 启动服务

### 设置域名 MX 解析

MX 解析表示使用指定的服务器处理发送至该域名的邮件，请于您的 DNS 解析服务商处填写您使用的服务器地址即可。

## 二次开发

您可以 fork 此项目以进行二次开发使用，请自行调整相关的设置内容。

## 调试

如有疑问，请检查上述步骤中是否存在疏漏的部分，如果依然无法解决，可以尝试将 `config.json` 中的 `system.production
 设置为 false ，以启用调试模式查看详细的日志。

如果您还是有疑问，您可以开启一个 issue ，附上您的完整复现流程，以获取更多帮助。
