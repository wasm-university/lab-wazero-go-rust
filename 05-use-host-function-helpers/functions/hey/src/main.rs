use pre_hf;


fn main() {
  // foo
}


fn print_hello(name: &String) {
  //pre_hf::log(&"👋👋👋 hello world 🌍".to_owned());

  pre_hf::log(&["👋 hello ", name].concat());
  //log(&["👋 hello ", name].concat());
}

/// WebAssembly export that accepts a string (linear memory offset, byteCount)
/// and calls [`greet`].
///
/// Note: The input parameters were returned by [`allocate`]. This is not an
/// ownership transfer, so the inputs can be reused after this call.
#[cfg_attr(all(target_arch = "wasm32"), export_name = "print_hello")]
#[no_mangle]
pub unsafe extern "C" fn _print_hello(ptr: u32, len: u32) {
  print_hello(&pre_hf::ptr_to_string(ptr, len));
}



#[no_mangle]
pub fn add(a: i32, b: i32) -> i32 {
  pre_hf::log(&"👋👋👋 hello world 🌍".to_owned());
  pre_hf::log(&"👋👋 hello world 🌍".to_owned());
  pre_hf::log(&"👋 hello world 🌍".to_owned());

  return a + b;
}
