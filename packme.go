package packme

import (
	"encoding/json"
	"io/ioutil"
	"unsafe"

	"github.com/bytecodealliance/wasmtime-go/v17"
)

type Packme struct {
	instance            *wasmtime.Instance
	store               *wasmtime.Store
	wasm_vector_len     int
	cached_uint8_memory []byte
}

func (p Packme) New(file_path string) Packme {
	data, _ := ioutil.ReadFile(file_path)
	store := wasmtime.NewStore(wasmtime.NewEngine())
	wasm_module, _ := wasmtime.NewModule(store.Engine, data)
	instance, _ := wasmtime.NewInstance(store, wasm_module, nil)
	return Packme{instance, store, 0, []byte{}}
}

func (p *Packme) getUint8Memory() []byte {
	if len(p.cached_uint8_memory) == 0 {
		memory := p.instance.GetExport(p.store, "memory").Memory()
		p.cached_uint8_memory = memory.UnsafeData(p.store)
	}
	return p.cached_uint8_memory
}

func (p *Packme) getInt32Memory() []int32 {
	memory := p.instance.GetExport(p.store, "memory").Memory()
	data := memory.UnsafeData(p.store)
	int32Data := *(*[]int32)(unsafe.Pointer(&data))
	return int32Data
}

func (p *Packme) encodeString(str string, view []byte) int {
	for i := 0; i < len(str); i++ {
		view[i] = str[i]
	}
	return len(str)
}

func (p *Packme) passStringToWasm(str string) int32 {
	str_len := len(str)
	malloc := p.instance.GetFunc(p.store, "__wbindgen_malloc")
	realloc := p.instance.GetFunc(p.store, "__wbindgen_realloc")
	ptr, err := malloc.Call(p.store, str_len, 1)
	if err != nil {
		panic(err)
	}
	mem := p.getUint8Memory()
	offset := 0
	for ; offset < str_len; offset += 1 {
		if int(str[offset]) == 0x7f {
			break
		}
		mem[ptr.(int32)+int32(offset)] = byte(str[offset])
	}

	if offset != str_len {
		if offset != 0 {
			str = str[:offset]
		}
		align := offset + len(str)*3
		ptr, _ = realloc.Call(p.store, ptr, align, 1)
		str_len = align
		view := mem[ptr.(int32)+int32(offset) : ptr.(int32)+int32(str_len)]
		ret := p.encodeString(str, view)
		offset += ret
		ptr, _ = realloc.Call(p.store, ptr, str_len, offset, 1)
	}

	p.wasm_vector_len = offset
	return ptr.(int32)
}

func (p *Packme) getStringFromWasm(ptr int32, len int32) string {
	mem := p.getUint8Memory()
	return string(mem[ptr : ptr+len])
}

func (p *Packme) Pack(input PackmeInput) PackmeOutput {
	str, _ := json.Marshal(input)
	addToStackPointer := p.instance.GetFunc(p.store, "__wbindgen_add_to_stack_pointer")
	freeWasm := p.instance.GetFunc(p.store, "__wbindgen_free")
	packWasm := p.instance.GetFunc(p.store, "pack")
	retptr, _ := addToStackPointer.Call(p.store, -16)
	ptr0 := p.passStringToWasm(string(str))
	packWasm.Call(p.store, retptr, ptr0, p.wasm_vector_len)
	r0 := p.getInt32Memory()[retptr.(int32)/4+0]
	r1 := p.getInt32Memory()[retptr.(int32)/4+1]
	defer freeWasm.Call(p.store, r0, r1, 1)
	defer addToStackPointer.Call(p.store, 16)
	strResult := p.getStringFromWasm(r0, r1)
	var result PackmeOutput
	json.Unmarshal([]byte(strResult), &result)
	return result
}
