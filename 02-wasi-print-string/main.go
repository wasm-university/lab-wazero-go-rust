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
		NewFunctionBuilder().WithFunc(logString).Export("log").
		Instantiate(ctx)
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

	mod, err := wasmRuntime.Instantiate(ctx, helloWasm)
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

}

func logString(ctx context.Context, m api.Module, offset, byteCount uint32) {
	buf, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	fmt.Println(string(buf))
}
