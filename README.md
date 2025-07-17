# mysql_install_v2




```
cd existing_repo
git remote add origin http://10.10.9.201:91/root/mysqldeploy.git
git branch -M main
git push -uf origin main
```

下载代码
```
git pull origin main
```

上传代码
```
git add . 
git commit -m "$(date)"
git push -u origin main -f
```


构建命令
```
go clean -cache -modcache -i -r
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o mysqldeploy -ldflags="-s -w" .
```