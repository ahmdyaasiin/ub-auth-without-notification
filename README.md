# Brawijaya University Authentication Without Notification

## Overview
This library helps make it easier to check if someone is really a student or staff at Universitas Brawijaya. It lets developers quickly confirm someone's identity using their UB email or student ID number (NIM).

## Installation
Create a new directory for your project and change into it. Next, initialize your project using Go modules by running the following command in your terminal:
```
go mod init module-name
```
Install the library with the command below:
```
go get github.com/ahmdyaasiin/ub-auth-without-notification
```

## Quickstart
Create a new file and copy the code below into the file.
```go
package main

import (
    "fmt"
    u "github.com/ahmdyaasiin/ub-auth-without-notification"
    "log"
)

func main() {
    //
    resp, err := u.AuthUB("USERNAME_HERE", "PASSWORD_HERE")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(resp)
}
```
The output of the file will be as follows:
```
{NIM FullName Email Faculty StudyProgram}
```