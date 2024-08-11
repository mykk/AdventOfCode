mod dragon_checksum {
    use std::iter::zip;

    fn generate_initial_data(starting_data: &str, limit: usize) -> String {
        if starting_data.len() >= limit {
            starting_data.chars().take(limit).collect()
        }
        else {
            let inverted_and_reversed =  starting_data.chars().rev().map(|c| { if c == '1' { '0' } else { '1' } }).collect::<String>();
            generate_initial_data(&format!("{}0{}", starting_data, inverted_and_reversed), limit)
        }
    }

    fn create_dragon_checksum_from_generated_data(generated_data: &str) -> String {
        if generated_data.len() % 2 == 0 {
            let pairs = zip(generated_data.chars().step_by(2), generated_data.chars().skip(1).step_by(2));
            create_dragon_checksum_from_generated_data(&pairs.map(|(a, b)|{ if a == b { '1' } else { '0' } }).collect::<String>())
        }
        else {
            generated_data.chars().collect()
        }
    }

    pub(crate) fn create_dragon_checksum(starting_data: &str, limit: usize) -> String {
        create_dragon_checksum_from_generated_data(&generate_initial_data(starting_data, limit))
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::dragon_checksum::create_dragon_checksum;

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();

    println!("part1: {}", create_dragon_checksum(&content, 272));
    println!("part2: {}", create_dragon_checksum(&content, 35651584));
}

#[cfg(test)]
mod tests {
    use crate::dragon_checksum::create_dragon_checksum;

    #[test]
    fn test_example() {
        assert_eq!("01100", create_dragon_checksum("10000", 20));
    }
}
