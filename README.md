# pandatv-highlights-downloader
下载熊猫TV的高能时刻

使用了 https://github.com/williamchanrico/m3u8-download 处理ffmpeg的调用

仓促写的，就不弄傻瓜模式了，需要一些技术基础使用

需要mac或者linux环境，要装有ffmpeg，需要go语言环境

要修改main.go中的

{TOKEN}与{HOSTID}

需要到主播直播间的首页，用chrome看network请求，搜索hostvideos关键字找到请求看。

需要在项目根目录创建一个叫video的文件夹

之后直接 go run ./main.go
