package core

import (
	"os"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func Env(key string, defaultVal string) string {
	godotenv.Load(".env")
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultVal
	}
	return value
}

func gcdExtended(x int, m int) (int,int,int){
	if (x == 0){
		return m, 0, 1
	}

	gcd, a, b := gcdExtended(m%x, x)

	x1 := b - (m / (x * 1.0)) * a
	y1 := a
	
	return gcd, x1, y1
}

func MultiplicativeInverse(x int, m int) int{
	if ((x % m) != 1){
		return -1;
	}

	_, a, _ := gcdExtended(x, m)

	return a % m;
}

func UrlValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		url, ok := fl.Field().Interface().(string)
		if ok {
			var re = regexp.MustCompile(`((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+(:[0-9]+)?|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)`)
			return re.MatchString(url)
		}
		return false
	}
}
