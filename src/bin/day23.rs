fn main() {
    use aoc_2016::utils::aoc_file;
    use aoc_2016::utils::assembunny;

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();

    let lines = content.lines().collect::<Vec<_>>();
    let assembly = assembunny::parse_assembly(&lines);

    let values = assembunny::execute_assembly(&assembly, [('a'.to_string(), 7)].into(), |_, _|panic!());
    println!("part1: {}", values[&'a'.to_string()]);

    let values = assembunny::execute_assembly(&assembly, [('a'.to_string(), 12)].into(), |_, _|panic!());
    println!("part2: {}", values[&'a'.to_string()]);
}

#[cfg(test)]
mod tests {
    use std::collections::HashMap;
    use aoc_2016::utils::assembunny;

    #[test]
    fn test_example() {
        let example = "cpy 2 a
        tgl a
        tgl a
        tgl a
        cpy 1 a
        dec a
        dec a";

        let lines = example.lines().collect::<Vec<_>>();

        let assembly = assembunny::parse_assembly(&lines);
        let values = assembunny::execute_assembly(&assembly, HashMap::new(), |_,_|panic!());
        assert_eq!(values[&'a'.to_string()], 3);
    }
}
