#[path = "../utils/aoc_file.rs"] mod aoc_file;

pub mod distance_solver {
    use std::collections::HashSet;

    #[derive(Debug, PartialEq, Copy, Clone)]
    pub enum Direction {
        North = 0,
        East,
        South,
        West
    }

    #[derive(Debug, PartialEq, Copy, Clone)]
    pub enum TurnDirection{
        Left = -1,
        Right = 1
    }

    #[derive(Debug, Clone)]
    pub struct Position {
        pos:(i32, i32),
        facing: Direction
    }

    #[derive(Debug)]
    pub struct MoveInstruction {
        turn_dir: TurnDirection,
        val: i32
    }

    impl Position {
        fn get_new_direction(direction: Direction, turn_direction: TurnDirection) -> Direction {            
            match ((turn_direction as i32) + (direction as i32) + 4) % 4 {
                0 => Direction::North,
                1 => Direction::East,
                2 => Direction::South,
                3 => Direction::West,
                _ => unreachable!(),
            }
        }

        fn get_new_position(position: &Position, mv: &MoveInstruction) -> Position {
            let new_direction = Position::get_new_direction(position.facing, mv.turn_dir);

            match new_direction {
                Direction::North => Position{pos: (position.pos.0, position.pos.1 + mv.val), facing: new_direction},
                Direction::South => Position{pos: (position.pos.0, position.pos.1 - mv.val), facing: new_direction},
                Direction::East => Position{pos: (position.pos.0 + mv.val, position.pos.1), facing: new_direction},
                Direction::West => Position{pos: (position.pos.0 - mv.val, position.pos.1), facing: new_direction},
            }
        }
        
        fn get_starting_position() -> Self {
            Position{pos: (0, 0), facing: Direction::North}
        }
    }

    pub fn get_full_distance(mv_instructions: &[MoveInstruction]) -> i32 {
        let position = mv_instructions.iter()
            .fold(Position::get_starting_position(), 
            |pos, mv| Position::get_new_position(&pos, mv));

        position.pos.0.abs() + position.pos.1.abs()
    }

    fn intersect_with_visited(visited: &mut HashSet::<(i32, i32)>, start: (i32, i32), end: (i32, i32)) -> Option<(i32, i32)> {        
        let (x_increment, y_increment) = 
        if end.0 != start.0 {
            if end.0 > start.0 { (1, 0) } else { (-1, 0) }
        } 
        else {
            if end.1 > start.1 { (0, 1) } else { (0, -1) } 
        };
        
        let mut current = start;
        while current != end {
            current = (current.0 + x_increment, current.1 + y_increment);
            if !visited.insert(current) {
                return Some(current);
            }
        }
        None
    }

    pub fn get_first_duplicate_position_dist(mv_instructions: &[MoveInstruction]) -> Option<i32> {
        let mut visited = HashSet::new();

        let mut position = Position::get_starting_position();
        visited.insert(position.pos);

        for mv in mv_instructions {
            let new_position = Position::get_new_position(&position, mv);
            if let Some(pos) = intersect_with_visited(&mut visited, position.pos, new_position.pos) {
                return Some(pos.0.abs() + pos.1.abs());
            }

            position = new_position;
        }
        None
    }

    pub fn parse_input(content: &str) -> Option<Vec<MoveInstruction>> {
        content.split(", ")
        .map(|instruction| {
            let (turn_dir, val_str) = instruction.split_at(1);

            match turn_dir {
                "R" => Some(MoveInstruction{turn_dir: TurnDirection::Right, val: val_str.parse::<i32>().ok()?}),
                "L" => Some(MoveInstruction{turn_dir: TurnDirection::Left, val: val_str.parse::<i32>().ok()?}),
                _ => None
            }
        })
        .collect()
    }
}

fn main() {    
    let mv_instructions = match aoc_file::open_and_read_file(&mut std::env::args()) {
        Ok(content) => distance_solver::parse_input(&content),
        Err(_) => {
            eprintln!("Error reading file");
            std::process::exit(1);
        }
    }.expect("Failed to parse the file");

    println!("The answer to part 1 is {}", distance_solver::get_full_distance(&mv_instructions));

    match distance_solver::get_first_duplicate_position_dist(&mv_instructions) {
        Some(val) => println!("The answer to part 2 is {}", val),
        None  => println!("No duplicate position found.")
    };
}

#[cfg(test)]
mod tests {
    use super::distance_solver;

    #[test]
    fn part1_1() {
        let mv_instructions = distance_solver::parse_input("R2, L3").unwrap();
        assert_eq!(5, distance_solver::get_full_distance(&mv_instructions))
    }

    #[test]
    fn part1_2() {
        let mv_instructions = distance_solver::parse_input("R2, R2, R2").unwrap();
        assert_eq!(2, distance_solver::get_full_distance(&mv_instructions))
    }

    #[test]
    fn part1_3() {
        let mv_instructions = distance_solver::parse_input("R5, L5, R5, R3").unwrap();
        assert_eq!(12, distance_solver::get_full_distance(&mv_instructions))
    }

    #[test]
    fn part2() {
        let mv_instructions = distance_solver::parse_input("R8, R4, R4, R8").unwrap();
        assert_eq!(4, distance_solver::get_first_duplicate_position_dist(&mv_instructions).unwrap())
    }
}
