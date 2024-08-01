fn main() {
    use aoc_2016::utils::aoc_file;
    use aoc_2016::utils::assembunny;

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();

    let lines = content.lines().collect::<Vec<_>>();
    let assembly = assembunny::parse_assembly(&lines);

    let registry_results = assembunny::execute_assembly(&assembly, [].into(), |_, _|panic!());
    println!("part1: {}", registry_results[&'a'.to_string()]);

    let registry_results = assembunny::execute_assembly(&assembly, [('c'.to_string(), 1)].into(), |_, _|panic!());
    println!("part2: {}", registry_results[&'a'.to_string()]);
}

#[cfg(test)]
mod tests {
    use aoc_2016::utils::assembunny;

    #[test]
    fn test_example() {
        let example = "cpy 41 a
            inc a
            inc a
            dec a
            jnz a 2
            dec a";

        let lines = example.lines().collect::<Vec<_>>();

        let assembly = assembunny::parse_assembly(&lines);
        let registry_results = assembunny::execute_assembly(&assembly, [].into(), |_,_|panic!());
        assert_eq!(registry_results[&'a'.to_string()], 42);
    }
}
