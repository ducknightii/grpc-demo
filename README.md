//  protoc --go-grpc_out=paths=source_relative:. pb/helloworld.proto

//  protoc --go_out=paths=source_relative:.  pb/helloworld.proto



- dail 过程，服务中断 再恢复 后续请求怎么建立连接的
  - wireshark 抓包可以看到服务端恢复后，会自动建立新的连接处理请求
- cancel 捕获
  - 客户端超时以及主动断开连接 `DeadlineExceeded desc = context deadline exceeded=INFO`
- metadata 作用
  - 可以获取ip以及一些header信息