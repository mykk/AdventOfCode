mod maze_of_twisty_little_cubicles {
    use std::{cmp::Ordering, collections::{BinaryHeap, HashMap}};
    use once_cell::sync::Lazy;

    #[derive(Debug, PartialEq, Clone, Eq, Hash)]
    struct Position {
        x: i32,
        y: i32,
        count: usize
    }
    impl PartialOrd for Position {
        fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
            Some(self.cmp(other))
        }
    }
    impl Ord for Position {
        fn cmp(&self, other: &Self) -> Ordering {
            other.count.cmp(&self.count)
        }
    }

    fn is_wall(x: i32, y: i32, favorite_number: i32) -> bool {
        let checksum = x*x + 3*x + 2*x*y + y + y*y + favorite_number;
        checksum.count_ones() % 2 == 1
    }

    fn moveable_position(x: i32, y: i32, favorite_number: i32, count: usize, visited: &HashMap<(i32, i32), usize>) -> bool {
        x >= 0 && y >= 0 && !is_wall(x, y, favorite_number) && (visited.get(&(x, y)).is_none() || visited.get(&(x, y)).is_some_and(|value| *value > count)) 
    }

    fn search_maze<FinishF, ResultF>(favorite_number: i32, is_finished: FinishF, get_result: ResultF) -> usize 
    where FinishF: Fn(&Position) -> bool, 
          ResultF: Fn(&Position, &HashMap<(i32, i32), usize>) -> usize
    {
        static MOVE_DIRECTIONS: Lazy<Vec<(i32, i32)>> = Lazy::new(||vec![(-1, 0), (1, 0), (0, -1), (0, 1)]);

        let mut positions = BinaryHeap::new();
        positions.push(Position{x: 1, y: 1, count: 0});

        let mut visited: HashMap<(i32, i32), usize> = HashMap::new();
        visited.insert((1, 1), 0);

        while let Some(current) = positions.pop() {
            if is_finished(&current) {
                return get_result(&current, &visited);
            }

            for move_direction in MOVE_DIRECTIONS.iter() {
                let new_pos = (current.x + move_direction.0, current.y + move_direction.1); 
                let current_count = current.count + 1;
                if moveable_position(new_pos.0, new_pos.1, favorite_number, current_count, &visited) {
                    positions.push(Position{x: new_pos.0, y: new_pos.1, count: current_count});
                    visited.insert((new_pos.0, new_pos.1), current_count);
                }    
            }
        }
        unreachable!();
    }

    pub(crate) fn find_shortest_path(x: i32, y: i32, favorite_number: i32) -> usize {
        search_maze(favorite_number, |position|position.x == x && position.y == y, |position, _|position.count)
    }

    pub(crate) fn find_visited_locations(favorite_number: i32) -> usize {
        search_maze(favorite_number, |position|position.count == 50, |_, visited|visited.len())
    }

}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::maze_of_twisty_little_cubicles::{find_shortest_path, find_visited_locations};

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    let favorite_number = content.parse::<i32>().unwrap();
    println!("part1: {}", find_shortest_path(31, 39, favorite_number));
    println!("part2: {}", find_visited_locations(favorite_number));
}

#[cfg(test)]
mod tests {
    use crate::maze_of_twisty_little_cubicles::find_shortest_path;

    #[test]
    fn test_example() {
        assert_eq!(11, find_shortest_path(7, 4, 10))
    }
}
