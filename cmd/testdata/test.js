/*
Some comment here
TODO: todo comment #A
FIXME: fixme comment #A
*/
import React from "react"
import {A, B, C} from "my/folder"

// TODO: todo comment #B
func someFunc() {
	var a = 1
	var b = 2 //FIXME: fixme comment #B

	return a + b // TODO: todo comment #C
}

func testThis(t) {
	/*
		TODO: todo comment #D
		FIXME: fixme comment #C
	*/
	fmt.Println("hi // TODO: not a comment")
	fmt.Println(`hi // TODO: not a comment`)    // TODO: todo comment #E
	fmt.Println(`hi /* TODO: not a */ comment`) // TODO: todo comment #F

	fmt.Println(`
some string starts here
// TODO: not a comment
/* TODO: not a comment either */`) ////// TODO: todo comment #G

	var o = {
		"a": "//some value",
		"b": '// TODO: this is a trap',
		"c": '/* TODO: This is another trap */'
	}

}

///// TODO: todo comment #G
