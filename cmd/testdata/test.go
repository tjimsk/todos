/*
Some comment here
TODO: todo comment #1
FIXME: fixme comment #1
*/
package test

import (
	"testing"
)

// TODO: todo comment #2
func someFunc() int {
	a := 1
	b := 2 //FIXME: fixme comment #2

	return a + b // TODO: todo comment #3
}

func testThis(t *testing.T) {
	/*
		TODO: todo comment #4
		FIXME: fixme comment #3
	*/
	fmt.Println("hi // TODO: not a comment")
	fmt.Println(`hi // TODO: not a comment`)    // TODO: todo comment #5
	fmt.Println(`hi /* TODO: not a */ comment`) // TODO: todo comment #6

	fmt.Println(`
some string starts here
// TODO: not a comment
/* TODO: not a comment either */`) ////// TODO: todo comment #7

}

///// TODO: todo comment #8
