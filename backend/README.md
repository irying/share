# 设计文档

### 目录结构

```
backend
└─── api  #gin api
│   |___ controllers
│   |___ middleware	
│   |___ proxy
│   │___ routers
│   │___ runtime # errors.log 
│   api.go
└─── conf 
└─── pb		 # protobul file
└─── runtime # errors.log,every day's request day
└─── scripts # wrk test lua scripts
└─── service # rpc service
└─── static	 # frontend js,css
└─── storage # db,redis
└─── system  # system attribute
|	|___ exception
│   |___ upload
└─── test  
└─── utils   
└─── vendor
└─── views   # html
│   .gitgnore
│   glide.lock 
│   glide.yaml
│   main.go    
│	README.MD
```

### 流程介绍

1. 根目录下，go run main.go，启动rpc服务

2. api目录下，go run api.go，启动api服务

3. 后台登陆：http://localhost:8080/backend/login 账号/密码：admin/111111

4. 请求先经过gin的接口处理，接口rpc调用内部service，处理后返回结果

   Request——> Gin Api ——> Service ——> Db/Redis ——> Response

5. Api接口定义在api/router.go文件，RPC接口定义参考pb/user/user.proto文件

6. 日志有api的日志及rpc服务日志，考虑到线上部署，web server前面会有nginx等LB 作负载均衡，觉得api的正常请求日志打到nginx的access.log即可，无需重复打(若需要记录，可以启动时重定向go run api.go >> request.log)。

   rpc日志每次调用都会记录到runtime下面的每天请求日志，error,fatal,panic会记录到errors.log

7. 单元测试，分为endpointTests和serviceTests，serviceTests正常，endpointTests有点问题，先用postman替代

   cd test/serviceTests，执行go test

8. 部署在vagrant，公司的vagrant包

### 设计细节

1. **JWT**

   JWT是一种Token的编码算法，服务器端负责根据一个密码和算法生成Token，然后发给客户端，客户端只负责后面每次请求都在请求里面带上这个Token，服务器负责验证这个Token是不是合法的，有没有过期等，并可以解析出subject和claim里面的数据。

   使用 JWT 的目的不是为了以任何方式隐藏或者混淆数据。使用 JWT 是为了保证发送的数据是由可信的来源创建的(**来源可靠性**)。

   JWT的**优点是无需存储**，只通过算法实现对Token合法性的验证，不依赖数据库，nosql等存储系统，因此可以做到跨服务器验证，只要密钥和算法相同，不同服务器程序生成的Token可以互相验证。

   JWT的**缺点是Token到期刷新问题**，Token设置了到期时间，如果用户频繁使用应用程序，就不希望用户每X分钟后需要重新登录。老是要用户登录，从用户体验的角度来说，肯定是非常糟糕的。

   为了解决这个刷新问题，我想到以下两种思路，最终选择第二种。

   ​	1）延长过期时间，比如一周，用户只要登录过一次，一周后再重新登录；

   ​	2）自动刷新，服务器端或者客户端都可以做，服务器端在生成token之后，存一份到redis中，并设置过期时间。用户在使用系统期间，如果JWT解密出来的token已过期，但redis中的未过期，用户照常使用系统，同时会重新刷新过期时间，相当于延长过期时间。用户退出系统时在redis中销毁这个token，用户长时间未使用系统，token过期自动销毁。

2. **Convey**

   GoConvey是一款针对Golang的测试框架，可以管理和运行测试用例，同时提供了丰富的断言函数，并支持很多 Web 界面特性。
   它对比其他框架的优势在于能够在Web页面进行单元测试，对于集成自动化测试非常方便，官网地址https://github.com/smartystreets/goconvey

3. **Wrk** 

   目前最常见的压力测试工具是ab，但ab单线程及不能构建动态请求的缺点不满足性能测试要求，所以我选用wrk，一款功能较为强大的测试工具，可以使用lua脚本来支持更为复杂的测试场景。


### 性能测试报告

1.wrk -t4 -c200 -s fixed_profile.lua http://localhost:8080 **QPS:3703**

> 200并发（固定用户）情况下，HTTP API QPS大于3000
>
> - 200个client，每个client模拟一个用户（因此需要200个不同的固定用户账号）

![image-20181008170038807](https://wx1.sinaimg.cn/mw690/e5b38bb8ly1fw19wnu9dxj20rl07075u.jpg)



2.wrk -t4 -c2000 -s fixed_profile.lua http://localhost:8080 **QPS:4680**

> 2000并发（固定用户）情况下，HTTP API QPS大于1500

![image-20181008170445209](https://wx4.sinaimg.cn/mw690/e5b38bb8ly1fw19wnwk1xj20re070abm.jpg)



3.wrk -t4 -c200 -s profile.lua http://localhost:8080 **QPS:4235**

> 200并发（随机用户）情况下，HTTP API QPS大于1000
>
> - 200个client，每个client每次随机从10,000,000条记录中选取一个用户，发起请求（如果涉及到鉴权，可以使用一个测试用的token）

![image-20181008170211454](https://wx2.sinaimg.cn/mw690/e5b38bb8ly1fw19wnvhbgj20re070dhe.jpg)



4.wrk -t4 -c2000 -s profile.lua http://localhost:8080  **QPS:2690**

> 2000并发（随机用户）情况下，HTTP API QPS大于800

![image-20181008200533091](https://wx2.sinaimg.cn/mw690/e5b38bb8ly1fw19wnwgdkj20j207kmyx.jpg)



测试显示：固定用户在200个并发下表现还不如随机用户，但并发上到2000以后，**固定用户的QPS差点是随机用户的2倍。**