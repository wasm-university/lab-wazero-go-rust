package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
  "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	wasmRuntime := wazero.NewRuntime(ctx)

	defer wasmRuntime.Close(ctx) // This closes everything this Runtime created.

	// host functions
	_, err := wasmRuntime.NewHostModuleBuilder("env").
  ExportFunction("log", logString).
		Instantiate(ctx, wasmRuntime)
	if err != nil {
		log.Panicln(err)
	}

	_, err = wasi_snapshot_preview1.Instantiate(ctx, wasmRuntime)
	if err != nil {
		log.Panicln(err)
	}

	// Load then Instantiate a WebAssembly module
	wasmPath1 := "./functions/hey/target/wasm32-wasi/debug/hey.wasm"
	//wasmPath2 := "./functions/hello/hello.wasm"
	helloWasm, err := os.ReadFile(wasmPath1)

	if err != nil {
		log.Panicln(err)
	}

	mod, err := wasmRuntime.InstantiateModuleFromBinary(ctx, helloWasm)
	if err != nil {
		log.Panicln(err)
	}

	// Get references to WebAssembly function: "add"
	addWasmModuleFunction := mod.ExportedFunction("add")

	// Now, we can call "add", which reads the data we wrote to memory!
	// result []uint64
	result, err := addWasmModuleFunction.Call(ctx, 20, 22)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println("result:", result[0])

  funcPrintHello := mod.ExportedFunction("print_hello")
	allocate := mod.ExportedFunction("allocate")
	deallocate := mod.ExportedFunction("deallocate")

  name := "Bob Morane"
	nameSize := uint64(len(name))

	// Instead of an arbitrary memory offset, use Rust's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	results, err := allocate.Call(ctx, nameSize)
	if err != nil {
		log.Panicln(err)
	}
	namePtr := results[0]
	// This pointer was allocated by Rust, but owned by Go, So, we have to
	// deallocate it when finished
	defer deallocate.Call(ctx, namePtr, nameSize)

	// The pointer is a linear memory offset, which is where we write the name.
	if !mod.Memory().Write(ctx, uint32(namePtr), []byte(name)) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
			namePtr, nameSize, mod.Memory().Size(ctx))
	}

  nickName := "Bob Morane"
	nickNameSize := uint64(len(nickName))
	resultsBis, err := allocate.Call(ctx, nickNameSize)
	if err != nil {
		log.Panicln(err)
	}
  nickNamePtr := resultsBis[0]
	defer deallocate.Call(ctx, nickNamePtr, nickNameSize)

	if !mod.Memory().Write(ctx, uint32(nickNamePtr), []byte(nickName)) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
    nickNamePtr, nickNameSize, mod.Memory().Size(ctx))
	}


	// Now, we can call "greet", which reads the string we wrote to memory!
	_, err = funcPrintHello.Call(ctx, nickNamePtr, nickNameSize)
	if err != nil {
		log.Panicln(err)
	}


}

func logString(ctx context.Context, m api.Module, offset, byteCount uint32) {
	buf, ok := m.Memory().Read(ctx, offset, byteCount)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	fmt.Println(string(buf))
}
