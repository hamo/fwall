** client **
+ support UDP protocol

** compression **

** encrypt **
+ implement none encrypt method
  - Done.

** protocol **
+ support UDP protocol
+ pass userdb into protocol-server
+ Move write/read master/user to protocol since we would have a fixed way to encrypt header

** server **
+ support UDP protocol

** tunnel **

** userdb **
+ refactoring userdb as interface
+ support more backend, such as redis and mongodb
+ support more content in userdb, such as traffic statistics
