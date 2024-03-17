use std::collections::HashMap;

pub mod door_auth {
    use rayon::iter::{IntoParallelIterator, IntoParallelRefIterator, ParallelIterator};
    use regex::Regex;

    #[derive(PartialEq, Debug)]
    pub(crate) enum Instruction {
        Rect(usize, usize),
        RotateRow(usize, usize),
        RotateColumn(usize, usize)
    }

    pub(crate) fn parse_input(lines: &[&str]) -> Option<Vec<Instruction>> {
        let rect_regex = Regex::new(r"rect (\d+)x(\d+)").unwrap();
        let row_regex = Regex::new(r"rotate row y=(\d+) by (\d+)").unwrap();
        let column_regex = Regex::new(r"rotate column x=(\d+) by (\d+)").unwrap();

        let unpack_captures = |caps: &regex::Captures| -> Option<(usize, usize)> { 
            Some((caps[1].parse().ok()?, caps[2].parse().ok()?))
        };

        lines.iter().map(|line| {            
            if let Some(caps) = rect_regex.captures(line) {
                let (a, b) = unpack_captures(&caps)?;
                return Some(Instruction::Rect(a, b));
            }

            if let Some(caps) = row_regex.captures(line) {
                let (a, b) = unpack_captures(&caps)?;
                return Some(Instruction::RotateRow(a, b));
            }

            if let Some(caps) = column_regex.captures(line) {
                let (a, b) = unpack_captures(&caps)?;
                return Some(Instruction::RotateColumn(a, b));
            }

            None
        }).collect::<Option<Vec<_>>>()
    }

    fn get_initial_lock(wide: usize, tall: usize) -> Vec<Vec<bool>> {
        let mut the_lock = Vec::new();
        the_lock.resize_with(tall, ||{
            let mut row = Vec::new();
            row.resize(wide, false);
            row
        });
        the_lock
    }

    pub(crate) fn rotate_row(row: &[bool], right_shift: usize) -> Vec<bool> {
        let width = row.len();
        let right_shift = right_shift % width;
        (0..width).into_par_iter().map(|i| {
            let remap_to = if i >= right_shift {
                i - right_shift
            }
            else {
                width - (right_shift - i)
            };
            row[remap_to]
        }).collect()
    }

    pub(crate) fn get_display(wide: usize, tall: usize, instructions: &[Instruction]) -> Vec<Vec<bool>> {
        instructions.iter().fold(get_initial_lock(wide, tall), |mut the_lock, instruction| {
            match instruction {
                Instruction::Rect(width, height) => {
                    (0..*height).for_each(|i| { (0..*width).for_each(|j| { the_lock[i][j] = true; }); });
                }
                Instruction::RotateRow(row, right_shift) => { 
                    the_lock[*row] = rotate_row(&the_lock[*row], *right_shift);
                }
                Instruction::RotateColumn(column, down_shift) => {
                    let column_as_row = the_lock.iter().fold(Vec::new(), |mut column_as_row, row| {
                        column_as_row.push(row[*column]);
                        column_as_row
                    });

                    rotate_row(&column_as_row, *down_shift).iter().enumerate().for_each(|(i, val)| {
                        the_lock[i][*column] = *val;
                    });
                }
            }
            the_lock
        })        
    }

    pub(crate) fn get_door_sum(door_display: &[Vec<bool>]) -> usize {
        door_display.par_iter()
            .map(|row| { row.par_iter().fold(||0, |count, p| count + 1 * *p as usize).sum::<usize>() })
            .sum()
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    let instructions = door_auth::parse_input(&content.lines().collect::<Vec<&str>>()).expect("invalid file");

    let display = door_auth::get_display(50, 6, &instructions);
    println!("part1: {}", door_auth::get_door_sum(&display));

    let mut remap = HashMap::new();
    remap.insert(true, '*');
    remap.insert(false, ' ');

    let display = display.iter().map(|row|{
        row.iter().map(|val| remap[val].to_string()).collect::<String>()
    }).collect::<Vec<_>>();

    println!("part2:");
    for line in display {
        println!("{}", line);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_rect_parse() {
        let test = door_auth::parse_input(&["rect 5x2"]).unwrap();
        assert_eq!(1, test.len());
        assert_eq!(test[0], door_auth::Instruction::Rect(5, 2));
    }

    #[test]
    fn test_row_parse() {
        let test = door_auth::parse_input(&["rotate row y=8 by 2"]).unwrap();
        assert_eq!(1, test.len());
        assert_eq!(test[0], door_auth::Instruction::RotateRow(8, 2));
    }

    #[test]
    fn test_column_parse() {
        let test = door_auth::parse_input(&["rotate column x=8 by 2"]).unwrap();
        assert_eq!(1, test.len());
        assert_eq!(test[0], door_auth::Instruction::RotateColumn(8, 2));
    }

    #[test]
    fn test_parse() {
        let test = door_auth::parse_input(&["rect 5x2", "rotate row y=8 by 2", "rotate row y=5 by 4", "rotate column x=3 by 1"]).unwrap();
        assert_eq!(4, test.len());
        assert_eq!(test[0], door_auth::Instruction::Rect(5, 2));
        assert_eq!(test[1], door_auth::Instruction::RotateRow(8, 2));
        assert_eq!(test[2], door_auth::Instruction::RotateRow(5, 4));
        assert_eq!(test[3], door_auth::Instruction::RotateColumn(3, 1));
    }

    #[test]
    fn test_right_shift() {
        assert_eq!(vec![false, true, true, false], door_auth::rotate_row(&[true, true, false, false], 1));
        assert_eq!(vec![false, false, true, true], door_auth::rotate_row(&[true, true, false, false], 2));
        assert_eq!(vec![true, false, false, true], door_auth::rotate_row(&[true, true, false, false], 3));
        assert_eq!(vec![true, true, false, false], door_auth::rotate_row(&[true, true, false, false], 4));
        assert_eq!(vec![false, true, true, false], door_auth::rotate_row(&[true, true, false, false], 5));
    }

    #[test]
    fn test_get_door_sum() {
        let instructions = door_auth::parse_input(&["rect 3x2", 
            "rotate column x=1 by 1",
            "rotate row y=0 by 4",
            "rotate column x=1 by 1"]).unwrap();
        assert_eq!(6, door_auth::get_door_sum(&door_auth::get_display(3, 2, &instructions)));
    }
}
