use pre_hf;

fn main() {
  // foo
  pre_hf::log(&"🚀🚀🚀".to_owned());

}


fn print_hello(name: &String) {
  pre_hf::log(&["👋 hello ", name].concat());
}

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
