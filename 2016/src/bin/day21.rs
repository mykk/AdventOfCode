mod scrambled_letters_and_hash {
    

    #[derive(Clone, Copy, Debug, PartialEq, Eq)]
    pub(crate) enum Rotate {
        Left,
        Right
    }

    impl Rotate {
        pub(crate) fn inverse(&self) -> Rotate {
            match &self {
                Rotate::Left => Rotate::Right,
                Rotate::Right => Rotate::Left
            }
        }
    }

    #[derive(Clone, Copy, Debug, PartialEq, Eq)]
    pub(crate) enum Scrambling {
        SwapPosition(usize, usize),
        SwapLetter(char, char),
        Rotate(Rotate, usize),
        RotateLetterBased(char),
        Reverse(usize, usize),
        Move(usize, usize)
    }

    impl Scrambling {
        fn rotate_letter_based(seed: &str, x: char) -> String {
            let x = seed.chars().position(|c|{ c == x}).unwrap_or(0);
            let x = if x >= 4 { x + 2 } else { x + 1 };
            Self::rotate_sized(seed, Rotate::Right, x)
        }
    
        fn swap_position(seed: &str, x: usize, y: usize) -> String {
            seed.char_indices().map(|(i, c)| {
                if i == x { seed.chars().nth(y).unwrap() } 
                else if i == y { seed.chars().nth(x).unwrap() } 
                else { c }
            }).collect()
        }
    
        fn swap_letter(seed: &str, x: char, y: char) -> String {
            seed.chars().map(|c| {
                if c == x { y } else if c == y { x } else { c }
            }).collect()
        }
    
        fn rotate_sized(seed: &str, rotate: Rotate, x: usize) -> String {
            let x = x % seed.len();
            match rotate {
                Rotate::Left => format!("{}{}", &seed[x..], &seed[..x]),
                Rotate::Right => format!("{}{}", &seed[seed.len()- x..], &seed[..seed.len() - x])
            }
        }
    
        fn reverse(seed: &str, x: usize, y: usize) -> String {
            format!("{}{}{}", &seed[..x], &seed[x..=y].chars().rev().collect::<String>(), &seed[y + 1..])
        }
    
        fn move_from_to_position(seed: &str, x: usize, y: usize) -> String {
            let temp = format!("{}{}", &seed[..x], &seed[x + 1..]);
            format!("{}{}{}", &temp[..y], seed.chars().nth(x).unwrap(), &temp[y..])
        }

        fn reverse_rotate_letter_based(seed: &str, x: char) -> String {
            let mut reversed = seed.to_string();
    
            loop {
                if seed == Self::rotate_letter_based(&reversed, x) {
                    return reversed;
                }
                reversed = Self::rotate_sized(&reversed, Rotate::Left, 1);
            }
        }
    
        pub(crate) fn scramble(&self, seed: &str) -> String {
            match self {
                Scrambling::SwapPosition(x, y) => Self::swap_position(seed, *x, *y),
                Scrambling::SwapLetter(x, y) => Self::swap_letter(seed, *x, *y),
                Scrambling::Rotate(rotate, x) => Self::rotate_sized(seed, *rotate, *x),
                Scrambling::RotateLetterBased(x) => Self::rotate_letter_based(seed, *x),
                Scrambling::Reverse(x, y) => Self::reverse(seed, *x, *y),
                Scrambling::Move(x, y) => Self::move_from_to_position(seed, *x, *y)
            }
        }

        pub(crate) fn unscramble(&self, seed: &str) -> String {
            match self {
                Scrambling::SwapPosition(x, y) => Self::swap_position(seed, *x, *y),
                Scrambling::SwapLetter(x, y) => Self::swap_letter(seed, *x, *y),
                Scrambling::Rotate(rotate, x) => Self::rotate_sized(seed, rotate.inverse(), *x),
                Scrambling::RotateLetterBased(x) => Self::reverse_rotate_letter_based(seed, *x),
                Scrambling::Reverse(x, y) => Self::reverse(seed, *x, *y),
                Scrambling::Move(x, y) => Self::move_from_to_position(seed, *y, *x)
            }
        }
    }

    type ScramblingFunc = Box<dyn Fn(regex::Captures) -> Scrambling>;
    pub(crate) fn parse(lines: &[&str]) -> Option<Vec<Scrambling>> {
        let scramble_regex: Vec<(regex::Regex, ScramblingFunc)> = vec![
            (regex::Regex::new(r"swap position (\d+) with position (\d+)").unwrap(), Box::new(|capture|Scrambling::SwapPosition(capture[1].parse().unwrap(), capture[2].parse().unwrap()))),
            (regex::Regex::new(r"swap letter (\w) with letter (\w)").unwrap(), Box::new(|capture|Scrambling::SwapLetter(capture[1].chars().next().unwrap(), capture[2].chars().next().unwrap()))), 
            (regex::Regex::new(r"rotate left (\d+) step").unwrap(), Box::new(|capture|Scrambling::Rotate(Rotate::Left, capture[1].parse().unwrap()))),
            (regex::Regex::new(r"rotate right (\d+) step").unwrap(), Box::new(|capture|Scrambling::Rotate(Rotate::Right, capture[1].parse().unwrap()))),
            (regex::Regex::new(r"rotate based on position of letter (\w)").unwrap(), Box::new(|capture|Scrambling::RotateLetterBased(capture[1].chars().next().unwrap()))),
            (regex::Regex::new(r"reverse positions (\d+) through (\d+)").unwrap(), Box::new(|capture|Scrambling::Reverse(capture[1].parse().unwrap(), capture[2].parse().unwrap()))),
            (regex::Regex::new(r"move position (\d+) to position (\d+)").unwrap(), Box::new(|capture|Scrambling::Move(capture[1].parse().unwrap(), capture[2].parse().unwrap())))
        ];

        lines.iter().map(|line| {
            scramble_regex.iter().find_map(|(reg, factory)|{ Some(factory(reg.captures(line)?)) })
        }).collect()
    }

    pub(crate) fn scramble(seed: &str, scramblings: &[Scrambling]) -> String {
        scramblings.iter().fold(seed.to_string(), |seed, scramble| scramble.scramble(&seed))
    }

    pub(crate) fn unscramble(seed: &str, scramblings: &[Scrambling]) -> String {
        scramblings.iter().rev().fold(seed.to_string(), |seed, scramble| scramble.unscramble(&seed))
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::scrambled_letters_and_hash::{parse, scramble, unscramble};

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    let lines: Vec<_> = content.lines().collect();
    let scramblings = parse(&lines).expect("Failed to parse the file");
    println!("part1: {}", scramble("abcdefgh", &scramblings));
    println!("part1: {}", unscramble("fbgdceah", &scramblings));
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_example() {
        use crate::scrambled_letters_and_hash::{parse, scramble, unscramble};
        let scramblings = parse(
            &["swap position 4 with position 0",
                    "swap letter d with letter b",
                    "reverse positions 0 through 4",
                    "rotate left 1 step",
                    "move position 1 to position 4",
                    "move position 3 to position 0",
                    "rotate based on position of letter b",
                    "rotate based on position of letter d"]).unwrap();

        assert_eq!("decab", scramble("abcde", &scramblings));
        assert_eq!("abcde", unscramble("decab", &scramblings));
    }
}
