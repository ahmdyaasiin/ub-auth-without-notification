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
