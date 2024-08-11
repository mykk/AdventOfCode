mod one_time_pad {
    use once_cell::sync::Lazy;
    use md5;
    use fancy_regex;
    use regex;

    fn get_digest(salt: &str, index: u32, stretching: u32) -> String {
        let digest = format!("{:x}", md5::compute(format!("{}{}", salt, index)));

        (0..stretching).fold(digest,|digest, _|{
            format!("{:x}", md5::compute(digest)
        )})
    }

    pub(crate) fn generate_new_keys(limit: usize, salt: &str, stretching: u32) -> u32 {
        static THREE_REPEATS_REGEX: Lazy<fancy_regex::Regex> = Lazy::new(||fancy_regex::Regex::new(r"(.)\1\1").unwrap());

        let mut potential_new_keys = Vec::new();
        let mut new_keys = Vec::new();
        let mut i = 0;
        loop {
            let digest = get_digest(salt, i, stretching);
            if let Some(captures) = THREE_REPEATS_REGEX.captures(&digest).unwrap() {
                let repeat_char = captures.get(1).unwrap().as_str();
                let five_repeat_regex = regex::Regex::new(&format!(r"{}{{5}}", repeat_char)).unwrap();
                potential_new_keys.push((i, five_repeat_regex));
            }
            potential_new_keys = potential_new_keys.into_iter().filter(|(index, _)| *index + 1000 >= i).collect();
            
            new_keys.append(&mut potential_new_keys.iter().filter(|(index, re)| *index != i && re.is_match(&digest)).map(|(index, _)|*index).collect());
            potential_new_keys = potential_new_keys.into_iter().filter(|(index, re)| *index == i || !re.is_match(&digest)).collect();

            if new_keys.len() >= limit {
                new_keys.sort();
                return new_keys[limit - 1];
            }
            i += 1;
        }
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::one_time_pad::generate_new_keys;

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    println!("part1: {}", generate_new_keys(64, &content, 0));
    println!("part2: {}", generate_new_keys(64, &content, 2016));
}

#[cfg(test)]
mod tests {
    use crate::one_time_pad::generate_new_keys;

    #[test]
    fn test_example() {
        assert_eq!(22728, generate_new_keys(64, "abc", 0));
        assert_eq!(22551, generate_new_keys(64, "abc", 2016));
    }
}
