pub(crate) mod decode {
    use std::collections::HashMap;

    pub(crate) fn decode<T>(message: &[T]) -> Vec<HashMap<char, i32>>
    where T: AsRef<str> 
    {
        message.iter().fold(Vec::new(), |mut counter, line| {
            line.as_ref().chars().enumerate().for_each(|(i, c)| {
                if counter.get(i).is_none() { counter.push(HashMap::new()); }
                *counter[i].entry(c).or_insert(0) += 1;
            });
            counter
        })
    }

    pub(crate) fn decode_as_most_common(counter: &[HashMap<char, i32>]) -> String {
        counter.iter().map(|char_counter| {
            char_counter.iter().max_by_key(|(_, i)| *i).unwrap().0
        }).collect()
    }

    pub(crate) fn decode_as_least_common(counter: &[HashMap<char, i32>]) -> String {
        counter.iter().map(|char_counter| {
            char_counter.iter().min_by_key(|(_, i)| *i).unwrap().0
        }).collect()
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    
    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    let counter = decode::decode(&content.lines().collect::<Vec<_>>());

    println!("part1 {}", decode::decode_as_most_common(&counter));
    println!("part2 {}", decode::decode_as_least_common(&counter));
}

#[cfg(test)]
mod tests {
    use crate::decode;

    #[test]
    fn decode_most_common_single_line() {
        let test = decode::decode_as_most_common(&decode::decode(&["test"]));
        assert_eq!("test", test);
    }

    #[test]
    fn decode_most_common_two_lines() {
        let test = decode::decode_as_most_common(&decode::decode(&["test", "test"]));
        assert_eq!("test", test);
    }

    #[test]
    fn decode_most_common_three_lines() {
        let test = decode::decode_as_most_common(&decode::decode(&["text", "taso", "rast"]));
        assert_eq!("tast", test);
    }

    #[test]
    fn decode_least_common_single_line() {
        let test = decode::decode_as_least_common(&decode::decode(&["test"]));
        assert_eq!("test", test);
    }

    #[test]
    fn decode_most_least_two_lines() {
        let test = decode::decode_as_least_common(&decode::decode(&["test", "test"]));
        assert_eq!("test", test);
    }

    #[test]
    fn decode_least_common_three_lines() {
        let test = decode::decode_as_least_common(&decode::decode(&["text", "taso", "rast"]));
        assert_eq!("rexo", test);
    }
}
