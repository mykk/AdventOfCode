use std::convert::AsRef;
use std::num::ParseIntError;
use thiserror::Error;
use std::thread;
use std::sync::Arc;
use rayon::prelude::*;

#[path = "../utils/aoc_file.rs"] mod aoc_file;

#[derive(Debug, Clone, Copy, PartialEq)]
struct Triangle(u32, u32, u32);

#[derive(Error, Debug)]
pub enum TriangleError {
    #[error("Failed to parse side as u32")]
    IntParseError,
    #[error("Failed to parse file")]
    ParsingError,
    #[error("Internal error")]
    InternalError,
}

fn possible_triangle(triangle: &Triangle) -> bool {
    let largest_side = triangle.0.max(triangle.1).max(triangle.2);
    largest_side <  triangle.0 + triangle.1 + triangle.2 - largest_side  
}

fn parse_sides(line: &str) -> Result<Vec<u32>, ParseIntError> {
    line.split_whitespace()
        .map(|side|side.parse::<u32>())
        .collect()
}

fn parse_data_part1<T>(lines: &[T]) -> Result<Vec<Triangle>, TriangleError> 
where
    T: AsRef<str> + Sync
{
    lines.par_iter()
        .map(|line| {
            let sides = parse_sides(line.as_ref()).map_err(|_| TriangleError::IntParseError)?;
            if let [a, b, c] = sides.as_slice() {
                Ok(Triangle(*a, *b, *c))
            } else {
                Err(TriangleError::ParsingError)
            }
        })
        .collect()
}

fn parse_data_part2<T>(lines: &[T]) -> Result<Vec<Triangle>, TriangleError> 
where T: AsRef<str>
{
    let mut triangles = Vec::new();
    
    let mut index = 0;
    let mut current_triangles = Vec::new();
    for line in lines {
        let sides = parse_sides(line.as_ref()).map_err(|_| TriangleError::IntParseError)?;

        match index {
            0 => current_triangles = sides.iter().map(|side| Triangle{0: *side, 1: 0, 2: 0}).collect(),
            1 => sides.iter().enumerate().for_each(|(i, side)| current_triangles[i].1 = *side),
            2 => {
                sides.iter().enumerate().for_each(|(i, side)| current_triangles[i].2 = *side);
                triangles.append(&mut current_triangles);
            }
            _ => return Err(TriangleError::InternalError)
        }
        index = (index + 1) % 3;
    }
    
    if index != 0 {
        return Err(TriangleError::ParsingError);
    }

    Ok(triangles)
}

fn count_possible_triangles(triangles: &[Triangle]) -> u32 {
    triangles.par_iter().filter(|triangle| possible_triangle(triangle)).count() as u32
}

fn parse_and_count_possible_triangles<F, T>(lines: &[T], parse: F) -> Result<u32, TriangleError> 
where 
    F: Fn(&[T]) -> Result<Vec<Triangle>, TriangleError>,
    T: AsRef<str>
{
    Ok(count_possible_triangles(&parse(lines)?))
}

fn main() {
    let possible_triangles = match aoc_file::open_and_read_file(&mut std::env::args()) {
        Ok(data) => 
        {
            let lines: Arc<Vec<String>> = Arc::new(data.split('\n')
                .map(|line| line.strip_prefix('\r')
                .unwrap_or(line).to_string())
                .collect());

            let lines_clone = Arc::clone(&lines);    
            let thread1 = thread::spawn(move|| {
                parse_and_count_possible_triangles(&lines_clone, parse_data_part1)
            });

            let thread2 = thread::spawn(move|| {
                parse_and_count_possible_triangles(&lines, parse_data_part2)
            });
            (thread1.join().unwrap(), thread2.join().unwrap())
        }
        Err(_) => {
            eprintln!("Error reading file");
            std::process::exit(1);
        }
    };
    match possible_triangles.0 {
        Ok(possible_triangles) => println!("Part 1: {}", possible_triangles),
        Err(err) => println!("Failed to get result for part 1: {}", err)
    }

    match possible_triangles.1 {
        Ok(possible_triangles) => println!("Part 2: {}", possible_triangles),
        Err(err) => println!("Failed to get result for part 2: {}", err)
    }
    
}

#[cfg(test)]
mod tests {
    use crate::possible_triangle;
    use crate::Triangle;

    #[test]
    fn test_impossible_triangle() {
        assert!(!possible_triangle(&Triangle{0: 5, 1: 10, 2: 25}));
    }

    #[test]
    fn test_possible_triangle() {
        assert!(possible_triangle(&Triangle{0: 10, 1: 16, 2: 25}));
    }

    #[test]
    fn test_impossible_triangle_() {
        assert!(!possible_triangle(&Triangle{0: 10, 1: 15, 2: 25}));
    }
}