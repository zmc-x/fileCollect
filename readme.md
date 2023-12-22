## FileCollect
This project is a personal project. It can help me collect files like coursework. This makes me really relaxed, because nobady will send some files to my QQ or Wechat.

## Features
Nowadays, this project has the following features. I will add features or fix bugs in the future.

+ Database table Automigration, don't execute sql file
+ Repositories are able to set deadlines
+ Yaml configuration files are supported
+ Use Redis as a server cache to improve access speed
+ Able to download and Upload (of course, this is a basic feature)
+ .....

## Install
If you want to use this project, you should clone this project in github. 

1. Modify related configurations in `config.yaml`
2. Build the project and run (If there is an error, you need to change the go version information in go.mod to your go version information, if there is still an error, it is recommended to upgrade the go version)
   ```shell
   go build main.go && ./main
   ```