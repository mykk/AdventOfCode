mod scrambled_letters_and_hash {
    use regex;

    #[derive(Clone, Copy, Debug, PartialEq, Eq)]
    pub(crate) enum Rotate {
        Left,
        Right
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

    pub(crate) fn parse(lines: &[&str]) -> Option<Vec<Scrambling>> {
        let scramble_regex: Vec<(regex::Regex, Box<dyn Fn(regex::Captures) -> Scrambling>)> = vec![
            (regex::Regex::new(r"swap position (\d+) with position (\d+)").unwrap(), Box::new(|capture|Scrambling::SwapPosition(capture[1].parse().unwrap(), capture[2].parse().unwrap()))),
            (regex::Regex::new(r"swap letter (\w) with letter (\w)").unwrap(), Box::new(|capture|Scrambling::SwapLetter(capture[1].chars().nth(0).unwrap(), capture[2].chars().nth(0).unwrap()))), 
            (regex::Regex::new(r"rotate left (\d+) step").unwrap(), Box::new(|capture|Scrambling::Rotate(Rotate::Left, capture[1].parse().unwrap()))),
            (regex::Regex::new(r"rotate right (\d+) step").unwrap(), Box::new(|capture|Scrambling::Rotate(Rotate::Right, capture[1].parse().unwrap()))),
            (regex::Regex::new(r"rotate based on position of letter (\w)").unwrap(), Box::new(|capture|Scrambling::RotateLetterBased(capture[1].chars().nth(0).unwrap()))),
            (regex::Regex::new(r"reverse positions (\d+) through (\d+)").unwrap(), Box::new(|capture|Scrambling::Reverse(capture[1].parse().unwrap(), capture[2].parse().unwrap()))),
            (regex::Regex::new(r"move position (\d+) to position (\d+)").unwrap(), Box::new(|capture|Scrambling::Move(capture[1].parse().unwrap(), capture[2].parse().unwrap())))
        ];

        lines.iter().map(|line| {
            scramble_regex.iter().find_map(|(reg, factory)|{ Some(factory(reg.captures(&line)?)) })
        }).collect()
    }

    fn rotate_letter_based(seed: &str, x: char) -> String {
        let x = seed.chars().position(|c|{ c == x}).unwrap_or(0);
        let x = if x >= 4 { x + 2 } else { x + 1 };
        let x = x % seed.chars().count();
        seed.chars().skip(seed.chars().count() - x).chain(seed.chars().take(seed.chars().count() - x)).collect()
    }

    fn swap_position(seed: &str, x: usize, y: usize) -> String {
        seed.chars().enumerate().map(|(i, c)| {
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
        let x = x % seed.chars().count();
        match rotate {
            Rotate::Left => seed.chars().skip(x).chain(seed.chars().take(x)).collect(),
            Rotate::Right => seed.chars().skip(seed.chars().count() - x).chain(seed.chars().take(seed.chars().count() - x)).collect()
        }
    }

    fn reverse(seed: &str, x: usize, y: usize) -> String {
        let chunk = seed.chars().skip(x).take(y - x + 1).collect::<Vec<_>>();
        let reversed = chunk.into_iter().rev().collect::<Vec<_>>();
        seed.chars().take(x).chain(reversed.into_iter()).chain(seed.chars().skip(y + 1)).collect()
    }

    fn move_from_to_position(seed: &str, x: usize, y: usize) -> String {
        let temp = seed.chars().take(x).chain(seed.chars().skip(x + 1)).collect::<String>();
        temp.chars().take(y).chain([seed.chars().nth(x).unwrap()]).chain(temp.chars().skip(y)).collect()
    }

    pub(crate) fn scramble(seed: &str, scramblings: &[Scrambling]) -> String {
        scramblings.iter().fold(seed.to_string(), |seed, scramble| {
            match &scramble {
                Scrambling::SwapPosition(x, y) => swap_position(&seed, *x, *y),
                Scrambling::SwapLetter(x, y) => swap_letter(&seed, *x, *y),
                Scrambling::Rotate(rotate, x) => rotate_sized(&seed, *rotate, *x),
                Scrambling::RotateLetterBased(x) => rotate_letter_based(&seed, *x),
                Scrambling::Reverse(x, y) => reverse(&seed, *x, *y),
                Scrambling::Move(x, y) => move_from_to_position(&seed, *x, *y)
            }
        })
    }

    fn reverse_rotate_letter_based(seed: &str, x: char) -> String {
        let mut reversed = seed.to_string();

        loop {
            if seed == rotate_letter_based(&reversed, x) {
                return reversed;
            }
            reversed = format!("{}{}", reversed.chars().skip(1).collect::<String>(), reversed.chars().nth(0).unwrap());
        }
    }

    fn unscramble_reversed(seed: &str, scramblings: &[&Scrambling]) -> String {
        scramblings.iter().fold(seed.to_string(), |seed, scramble| {
            match &scramble {
                Scrambling::SwapPosition(x, y) => swap_position(&seed, *x, *y),
                Scrambling::SwapLetter(x, y) => swap_letter(&seed, *x, *y),
                Scrambling::Rotate(rotate, x) =>  {
                    let rotate = if *rotate == Rotate::Left { Rotate::Right } else { Rotate::Left }
                    rotate_sized(&seed, rotate, *x)
                }
                Scrambling::RotateLetterBased(x) => reverse_rotate_letter_based(&seed, *x),
                Scrambling::Reverse(x, y) => reverse(&seed, *x, *y),
                Scrambling::Move(x, y) => move_from_to_position(&seed, *y, *x)
            }
        })
    }

    pub(crate) fn unscramble(seed: &str, scramblings: &[Scrambling]) -> String {
        let scramblings = scramblings.iter().rev().collect::<Vec<_>>();
        unscramble_reversed(seed, &scramblings)
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
