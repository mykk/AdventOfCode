pub fn open_and_read_file(args: &mut std::env::Args) -> std::io::Result<String> {
    match &args.nth(1) {
        Some(path) => std::fs::read_to_string(path),
        _ => std::fs::read_to_string(ask_file())
    }
}

fn ask_file() -> String {
    println!("Provide file to parse");
    text_io::read!()
}
