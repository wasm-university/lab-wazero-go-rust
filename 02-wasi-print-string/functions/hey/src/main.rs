extern crate alloc;
extern crate core;
extern crate wee_alloc;


/// Logs a message to the console using [`_log`].
fn log(message: &String) {
  unsafe {
      let (ptr, len) = string_to_ptr(message);
      _log(ptr, len);
  }
}

#[link(wasm_import_module = "env")]
extern "C" {
  /// WebAssembly import which prints a string (linear memory offset,
  /// byteCount) to the console.
  ///
  /// Note: This is not an ownership transfer: Rust still owns the pointer
  /// and ensures it isn't deallocated during this call.
  #[link_name = "log"]
  fn _log(ptr: u32, size: u32);
}

/// Returns a pointer and size pair for the given string in a way compatible
/// with WebAssembly numeric types.
///
/// Note: This doesn't change the ownership of the String. To intentionally
/// leak it, use [`std::mem::forget`] on the input after calling this.
unsafe fn string_to_ptr(s: &String) -> (u32, u32) {
  return (s.as_ptr() as u32, s.len() as u32);
}


fn main() {
    // foo
}


#[no_mangle]
pub fn add(a: i32, b: i32) -> i32 {
  log(&"👋 hello world 🌍".to_owned());
  return a + b;
}
