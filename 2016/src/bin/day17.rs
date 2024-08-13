mod two_steps_forward {
    use std::{cmp::Ordering, collections::{BinaryHeap, HashMap}};
    
    use once_cell::sync::Lazy;

    #[derive(Debug, PartialEq, Clone, Eq, Hash)]
    struct State {
        location: (i8, i8),
        path: String
    }

    impl PartialOrd for State {
        fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
            Some(self.cmp(other))
        }
    }

    impl Ord for State {
        fn cmp(&self, other: &Self) -> Ordering {
            other.path.len().cmp(&self.path.len())
        }
    }
    type DirFunctions = HashMap<char, fn((i8, i8)) -> (i8, i8)>;

    fn get_direction_operations() -> DirFunctions {
        let mut direction_operations: DirFunctions = HashMap::new();
        direction_operations.insert('U', |xy: (i8, i8)| (xy.0, xy.1 - 1));
        direction_operations.insert('D', |xy: (i8, i8)| (xy.0, xy.1 + 1));
        direction_operations.insert('L', |xy: (i8, i8)| (xy.0 - 1, xy.1));
        direction_operations.insert('R', |xy: (i8, i8)| (xy.0 + 1, xy.1));
        direction_operations
    } 

    const GOOD_CHARS: &[char] = &['b', 'c', 'd', 'e', 'f'];
    const DIRECTIONS: &[char] = &['U', 'D', 'L', 'R'];
    static DIRECTION_OPERATIONS: Lazy<DirFunctions> = Lazy::new(get_direction_operations);

    fn create_new_state(current_state: &State, digest: &str, dir: char) -> Option<State> {
        if !GOOD_CHARS.contains(&digest.chars().nth(DIRECTIONS.iter().position(|x|*x == dir).unwrap()).unwrap()) {
            return None;
        }

        let location = DIRECTION_OPERATIONS[&dir](current_state.location);
        if location.0 < 0 || location.1 < 0 || location.0 > 3 || location.1 > 3 {
            return None;
        }
        
        Some(State{location, path : format!("{}{}", current_state.path, dir)})
    }

    pub(crate) fn find_exits(passcode: &str) -> Vec<String> {
        let mut vault_paths = vec![];

        let mut states = BinaryHeap::new();
        states.push(State{location: (0, 0), path: String::new()});

        while let Some(state) = states.pop() {
            let digest = format!("{:x}", md5::compute(format!("{}{}", passcode, state.path)));
            for dir in DIRECTIONS.iter() {
                if let Some(new_state) = create_new_state(&state, &digest, *dir) {
                    if new_state.location.0 == 3 && new_state.location.1 == 3 {
                        vault_paths.push(new_state.path);
                    }
                    else {
                        states.push(new_state);
                    }
                }
            }
        }
        vault_paths
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::two_steps_forward::find_exits;

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();

    let exits = find_exits(&content);
    println!("part1: {}", exits.iter().min_by(|a, b|a.len().cmp(&b.len())).unwrap());
    println!("part2: {}", exits.iter().max_by(|a, b|a.len().cmp(&b.len())).unwrap().len());
}

#[cfg(test)]
mod tests {
    use crate::two_steps_forward::find_exits;

    #[test]
    fn test_example() {
        assert_eq!("DDRRRD", find_exits("ihgpwlah").iter().min_by(|a, b|a.len().cmp(&b.len())).unwrap());
        assert_eq!("DDUDRLRRUDRD", find_exits("kglvqrro").iter().min_by(|a, b|a.len().cmp(&b.len())).unwrap());
        assert_eq!("DRURDRUDDLLDLUURRDULRLDUUDDDRR", find_exits("ulqzkmiv").iter().min_by(|a, b|a.len().cmp(&b.len())).unwrap());

        assert_eq!(370, find_exits("ihgpwlah").iter().max_by(|a, b|a.len().cmp(&b.len())).unwrap().len());
        assert_eq!(492, find_exits("kglvqrro").iter().max_by(|a, b|a.len().cmp(&b.len())).unwrap().len());
        assert_eq!(830, find_exits("ulqzkmiv").iter().max_by(|a, b|a.len().cmp(&b.len())).unwrap().len());
    }
}
