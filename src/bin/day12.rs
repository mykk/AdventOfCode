use std::collections::HashMap;


mod leonardo_monorail {
    use std::collections::HashMap;

    use once_cell::sync::Lazy;
    use regex::Regex;

    pub(crate) enum Instructions {
        Cpy(String, String),
        Inc(String),
        Dec(String),
        Jnz(String, String)
    }

    fn extract_capture(captures: &regex::Captures, index: usize) -> String {
        captures.get(index).map_or_else(
            || panic!("Capture group {} not found", index),
            |capture| capture.as_str().to_string())
    }
    
    pub(crate) fn parse_assembly(lines: &[&str]) -> Vec<Instructions> {
        static PATTERNS: Lazy<Vec<Regex>> = Lazy::new(|| vec![
            Regex::new(r"cpy (-*\w+) (\w+)").unwrap(),
            Regex::new(r"inc (\w+)").unwrap(),
            Regex::new(r"dec (\w+)").unwrap(),
            Regex::new(r"jnz (-*\w+) (-*\w+)").unwrap(),
        ]);
    
        lines.iter().fold(Vec::new(), |mut vec, line| {
            for (i, pattern) in PATTERNS.iter().enumerate() {
                if let Some(captures) = pattern.captures(line) {
                    let instruction = match i {
                        0 => Instructions::Cpy(extract_capture(&captures, 1), extract_capture(&captures, 2)),
                        1 => Instructions::Inc(extract_capture(&captures, 1)),
                        2 => Instructions::Dec(extract_capture(&captures, 1)),
                        3 => Instructions::Jnz(extract_capture(&captures, 1), extract_capture(&captures, 2)),
                        _ => unreachable!(),
                    };
                    vec.push(instruction);
                    break;
                }
            }
            vec
        })
    }
    
    fn get_value(values: &HashMap<&str, i32>, val: &str) -> i32 {
        if let Some(val) = val.parse().ok() {
            val
        }
        else {
            *values.get(val).unwrap_or(&0)
        }
    }

    pub(crate) fn execute_assembly(assembly: &[Instructions], initial_vals: &HashMap<String, i32>) -> HashMap<String, i32> {
        let mut index: usize = 0;
        let mut values: HashMap<&str, i32> = initial_vals.iter().map(|(key, val)|((key.as_str(), *val))).collect();

        while index < assembly.len() {
            index = match &assembly[index] {
                Instructions::Cpy(value, registry) => {
                    values.insert(registry, get_value(&values, value));
                    index + 1
                }
                Instructions::Inc(registry) => {
                    values.entry(registry).and_modify(|value| *value += 1);
                    index + 1
                }
                Instructions::Dec(registry) => {
                    values.entry(registry).and_modify(|value| *value -= 1);
                    index + 1
                }
                Instructions::Jnz(value, jump) => {
                    if get_value(&values, value) != 0 {
                        let jump = get_value(&values, jump);
                        if jump < 0 { index - (-jump) as usize } else { index + jump as usize }
                    }
                    else {
                        index + 1
                    }
                }
            };
        }
        values.iter().map(|(key, val)| {(key.to_string(), *val)}).collect()
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::leonardo_monorail::{parse_assembly, execute_assembly};

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();

    let lines = content.lines().collect::<Vec<_>>();
    let assembly = parse_assembly(&lines);

    let registry_results = execute_assembly(&assembly, &HashMap::new());
    println!("part1: {}", registry_results[&'a'.to_string()]);

    let registry_results = execute_assembly(&assembly, &[('c'.to_string(), 1)].into_iter().collect());
    println!("part2: {}", registry_results[&'a'.to_string()]);
}

#[cfg(test)]
mod tests {
    use std::collections::HashMap;

    use crate::leonardo_monorail::{parse_assembly, execute_assembly};

    #[test]
    fn test_example() {
        let example = "cpy 41 a\n
            inc a\n
            inc a\n
            dec a\n
            jnz a 2\n
            dec a";

        let lines = example.lines().collect::<Vec<_>>();

        let assembly = parse_assembly(&lines);
        let registry_results = execute_assembly(&assembly, &HashMap::new());
        assert_eq!(registry_results[&'a'.to_string()], 42);
    }
}
