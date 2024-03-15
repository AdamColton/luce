package linject

import "testing"

// private seals FuncType so it can only be implemented here.
func TestFnTypesFulfillsFuncType(t *testing.T) {
	var ft FuncType = &fnTypes{}
	ft.private()
}
