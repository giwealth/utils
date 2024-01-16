package permission

import "testing"

func TestPermission(t *testing.T)  {
	if err := Generate("./httpapi", "./"); err != nil {
		panic(err)
	}
}