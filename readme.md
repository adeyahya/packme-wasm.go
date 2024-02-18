# packme-wasm.go

go binding for (packme-wasm)[https://github.com/adeyahya/packme-wasm] A 3D bin packing library in Rust/WebAssembly.

### Usage example

```go
func TestPacking(t *testing.T) {
	packme := Packme{}.New("./packme.wasm")
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
```
