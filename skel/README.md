# Go-Gin web application

This application was generated with [generator-gin](https://github.com/openware/generator-gin)

# How to run sonic with client

1. Put your frontend to client

2. Run makefile
```
make asset
```
It will run build in client and after that all build files will moved to public/assets.
#### **Be sure that you build your client (frontend) to the build folder, if it's another folder you can update your client (frontend) or makefile**

3. Run go server 
```
go run app.go serve
```

**if it's not work and you see white screen don't worry firstly check order of the import files in index.html**