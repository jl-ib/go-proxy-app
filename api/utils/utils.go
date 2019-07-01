package utils

import (
	"fmt"
	"os"
	// "path"
	env "github.com/joho/godotenv"
)
//
func LoadEnv() {
	// base, _ := filepath.Abs("../../")
	// envPath := path.Join(base, ".env")
	env.Load("/Users/jorge/go/src/github.com/jl-ib/proxy-app/.env")
	fmt.Println(os.Getenv("PORT"))
}