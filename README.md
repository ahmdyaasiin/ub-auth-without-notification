# Brawijaya University Authentication Without Notification

## Overview
This library helps make it easier to check if someone is really a student at Universitas Brawijaya. It lets developers quickly confirm someone's identity using their UB email or student ID number (NIM).

## Installation
Create a new directory for your project and change into it. Next, initialize your project using Go modules by running the following command in your terminal:
```
go mod init module-name
```
Install the library with the command below:
```
go get -u github.com/ahmdyaasiin/ub-auth-without-notification/v3
```

## Quickstart
Create a new file and copy the code below into the file.
```go
package main

import (
	"fmt"
	ub "github.com/ahmdyaasiin/ub-auth-without-notification/v3"
)

func main() {
	res, err := ub.Auth("EMAIL_OR_NIM_HERE", "PASSWORD_HERE")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", res)
}

```
If the email or NIM and password are correct, the output of the file should be as follows:
```
&{NIM FullName Email Faculty StudyProgram SIAKADPhotoURL FileFILKOMPhotoUrl}
```