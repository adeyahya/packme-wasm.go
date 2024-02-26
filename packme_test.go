package packme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestPacking(t *testing.T) {
	packme := New()
	input := PackmeInput{
		Containers: []ItemInput{
			{Id: "container 1", Qty: 1, Dim: [3]int{20, 20, 30}},
		},
		Items: []ItemInput{
			{Id: "item 1", Qty: 5, Dim: [3]int{10, 10, 30}},
		},
	}
	packed := packme.Pack(input)
	assert.Equal(t, len(packed.UnpackedItems), 1)
	assert.Equal(t, len(packed.Containers[0].Items), 4)
}

func TestVersion(t *testing.T) {
	packme := New()
	version := packme.Version()
	assert.GreaterOrEqual(t, len(version), 1)
}
