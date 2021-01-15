# Go-Gin web application

This application was generated with [sonic](https://github.com/openware/sonic)

# How to run sonic with client

1. Put your frontend application to the client folder

2. Run makefile
```
make asset
```
It will run build in the client folder and then all build files will be moved to public/assets.
#### **Be sure that you build your client (frontend) to the build folder, if it's another folder you can update your client (frontend) or makefile**

3. Run go server 
```
go run app.go serve
```

# Troubleshooting
**If it doesn't work and you see the white screen, check the order of the import files in index.html**