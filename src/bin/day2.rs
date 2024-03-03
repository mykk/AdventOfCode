#[path = "../utils/aoc_file.rs"] mod aoc_file;

pub mod keypad_solver {
    #[derive(Debug, Clone, Copy, PartialEq)]
    pub enum Direction {
        Up,
        Down,
        Right,
        Left
    }

    pub trait KeypadSolverInstructions {
        fn get_start_position(&self) -> (i8, i8);
        fn get_key(&self, position: (i8, i8)) -> char;
        fn get_new_position(&self, current: (i8, i8), delta: (i8, i8)) -> (i8, i8);
    }

    pub struct KeypadSolver {
        instructions: Box<dyn KeypadSolverInstructions>
    }

    impl KeypadSolver {
        pub fn new(instructions: Box<dyn KeypadSolverInstructions>) -> Self {
            KeypadSolver{instructions}
        }

        fn get_direction(dir: Direction) -> (i8, i8) {
            match dir {
                Direction::Up => (-1, 0),
                Direction::Down => (1, 0),
                Direction::Left => (0, -1),
                Direction::Right => (0, 1)
            }
        }

        pub fn solve(&self, directions: &Vec<Vec<Direction>>) -> String {
            let (sequence, _) = directions.iter().fold(
                (String::new(), self.instructions.get_start_position()), 
                |(mut sequence, position), direction| 
                {
                    let position = direction.iter().
                        fold(position, |pos, dir| self.instructions.get_new_position(pos, Self::get_direction(*dir)));

                    sequence.push(self.instructions.get_key(position));
                    (sequence, position)
                });
            sequence
        }

    }
    
    pub struct AocPart1Solver;

    impl AocPart1Solver {
        pub fn new() -> Self {
            Self{}
        }

        const KEY_PAD_SIZE: usize = 3;
        const START: (i8, i8) = (1, 1);
        const KEY_PAD: [[char; Self::KEY_PAD_SIZE]; Self::KEY_PAD_SIZE]  = [
            ['1', '2', '3'], 
            ['4', '5', '6'], 
            ['7', '8', '9']];
    }

    impl KeypadSolverInstructions for AocPart1Solver {
        fn get_start_position(&self) -> (i8, i8) {
            AocPart1Solver::START
        }

        fn get_key(&self, position: (i8, i8)) -> char {
            AocPart1Solver::KEY_PAD[position.0 as usize][position.1 as usize]
        }

        fn get_new_position(&self, current: (i8, i8), delta: (i8, i8)) -> (i8, i8) {
            ((current.0 + delta.0).max(0).min(Self::KEY_PAD_SIZE as i8 - 1), 
             (current.1 + delta.1).max(0).min(Self::KEY_PAD_SIZE as i8 - 1))
        }
    }

    pub struct AocPart2Solver;

    impl AocPart2Solver {
        pub fn new() -> Self {
            Self{}
        }

        const KEY_PAD_SIZE: usize = 5;
        const START: (i8, i8) = (2, 0);
        const KEY_PAD: [[char; Self::KEY_PAD_SIZE]; Self::KEY_PAD_SIZE]  = [
            [' ', ' ', '1', ' ', ' '],
            [' ', '2', '3', '4', ' '],
            ['5', '6', '7', '8', '9'],
            [' ', 'A', 'B', 'C', ' '],
            [' ', ' ', 'D', ' ', ' ']];
    }

    impl KeypadSolverInstructions for AocPart2Solver {
        fn get_start_position(&self) -> (i8, i8) {
            AocPart2Solver::START
        }

        fn get_key(&self, position: (i8, i8)) -> char {
            AocPart2Solver::KEY_PAD[position.0 as usize][position.1 as usize]
        }

        fn get_new_position(&self, current: (i8, i8), delta: (i8, i8)) -> (i8, i8) {
            let new_pos = ((current.0 + delta.0).max(0).min(Self::KEY_PAD_SIZE as i8 - 1), 
                                     (current.1 + delta.1).max(0).min(Self::KEY_PAD_SIZE as i8 - 1));
            match self.get_key(new_pos) {
                ' ' => current,
                _ => new_pos
            }
        }
    }

    fn char_to_direction(c: char) -> Option<Direction> {
        match c {
            'U' => Some(Direction::Up),
            'D' => Some(Direction::Down),
            'L' => Some(Direction::Left),
            'R' => Some(Direction::Right),
            _ => None
        }
    }

    pub fn parse_directions(content: &str) -> Option<Vec<Vec<Direction>>> {
        content.split_whitespace()
            .map(|word| word.chars().map(|c| char_to_direction(c)).collect())
            .collect()
    }
}

fn main() {    
    use crate::keypad_solver::{AocPart1Solver, AocPart2Solver, KeypadSolver};

    let directions = match aoc_file::open_and_read_file(&mut std::env::args()) {
        Ok(contents) => keypad_solver::parse_directions(&contents),
        Err(_) => {
            eprintln!("Error reading file");
            std::process::exit(1);
        }
    }.expect("Failed to parse the file");

    let solver = KeypadSolver::new(Box::new(AocPart1Solver::new()));
    println!("The answer to part 1 is {}", solver.solve(&directions));

    let solver = KeypadSolver::new(Box::new(AocPart2Solver::new()));
    println!("The answer to part 2 is {}", solver.solve(&directions));
}

#[cfg(test)]
mod tests {
    use crate::keypad_solver::{self, KeypadSolver, AocPart1Solver, AocPart2Solver};

    #[test]
    fn test_solver1() {
        let directions = keypad_solver::parse_directions("ULL\nRRDDD\nLURDL\nUUUUD").unwrap();
        let solver = KeypadSolver::new(Box::new(AocPart1Solver::new()));

        assert_eq!("1985", solver.solve(&directions));
    }

    #[test]
    fn test_solver2() {
        let directions = keypad_solver::parse_directions("ULL\nRRDDD\nLURDL\nUUUUD").unwrap();
        let solver = KeypadSolver::new(Box::new(AocPart2Solver::new()));

        assert_eq!("5DB3", solver.solve(&directions));
    }
}