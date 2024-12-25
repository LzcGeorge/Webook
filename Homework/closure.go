package Grammar

import "fmt"

func DeferClosureLoop1() {
	for i := 0; i < 10; i++ {
		defer func() {
			fmt.Printf("i 的地址是 %p，值是 %d\n", &i, i)
		}()
	}
}

func DeferClosureLoop2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			fmt.Printf("val 的地址是 %p，值是 %d\n", &val, val)
		}(i)
	}
}
func DeferClosureLoop3() {
	var j int
	for i := 0; i < 10; i++ {
		j = i
		defer func() {
			fmt.Printf("j 的地址是 %p，值是 %d\n", &j, j)
		}()
	}
}
