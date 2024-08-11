mod timing_is_everything {
    use once_cell::sync::Lazy;
    use regex;
    use num_integer;

    pub(crate) struct Disc {
        pub(crate) level: usize,
        pub(crate) position_at_zero: usize,
        pub(crate) available_positions: usize
    }

    pub(crate) fn parse(lines: &[&str]) -> Vec<Disc> {
        static DISC_REGEX: Lazy<regex::Regex> = Lazy::new(||regex::Regex::new(r"Disc .+ has (\d+) positions; at time=0, it is at position (\d+).").unwrap());

        lines.iter().enumerate().map(|(level, line)| {
            let captures = DISC_REGEX.captures(line).unwrap();
            Disc{level: level + 1, available_positions: captures[1].parse::<usize>().unwrap(), position_at_zero : captures[2].parse::<usize>().unwrap()}
        }).collect()
    }
    
    pub(crate) fn find_best_timing_single(current_time: usize, increment: usize, disc: &Disc) -> (usize, usize) {
        let new_time = (0..).find(|count| {
            let new_time = count * increment + current_time;
            (disc.position_at_zero + new_time + disc.level) % disc.available_positions == 0
        }).unwrap() * increment + current_time;

        (new_time, num_integer::lcm(increment, disc.available_positions))
    }
    
    pub(crate) fn find_best_timing(discs: &[Disc]) -> (usize, usize) {
        discs.iter().fold((0, 1), |(current_time, increment), disc|{
            find_best_timing_single(current_time, increment, disc)
        })
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::timing_is_everything::{parse, find_best_timing, find_best_timing_single, Disc};

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    let lines: Vec<_> = content.lines().collect();

    let best_timing = find_best_timing(&parse(&lines));
    println!("part1: {}", best_timing.0);
    let new_disc = Disc{level : lines.len() + 1, available_positions : 11, position_at_zero : 0};
    println!("part2: {}", find_best_timing_single(best_timing.0, best_timing.1, &new_disc).0);
}

#[cfg(test)]
mod tests {
    use crate::timing_is_everything::{find_best_timing, Disc};

    #[test]
    fn test_example() {
        assert_eq!(5, find_best_timing(&vec![Disc{level : 1, available_positions: 5, position_at_zero: 4}, 
            Disc{level : 2, available_positions : 2, position_at_zero : 1}]).0);
    }
}
