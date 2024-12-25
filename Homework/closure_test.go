package Grammar

import "testing"

func TestDeferClosureLoop1(t *testing.T) {
	DeferClosureLoop1()
	DeferClosureLoop2()
	DeferClosureLoop3()

}
