mod like_a_rogue {
    use once_cell::sync::Lazy;

    #[derive(Clone, Copy, Debug, PartialEq, Eq)]
    enum Tile {
        Safe,
        Trap
    }

    fn parse_first_row(first_row: &str) -> Vec<Tile> {
        first_row.chars().map(|x| {
            match x {
                '.' => Tile::Safe,
                '^' => Tile::Trap,
                _ => panic!("unexpected character in input")
            }            
        }).collect()
    }
    
    type TrapFunctions = Vec<fn(Tile, Tile, Tile) -> bool>;
    fn get_trap_functions() -> TrapFunctions {
        vec![
            |left: Tile, center: Tile, right: Tile| { left == Tile::Trap && center == Tile::Trap && right == Tile::Safe },
            |left: Tile, center: Tile, right: Tile| { left == Tile::Safe && center == Tile::Trap && right == Tile::Trap },
            |left: Tile, center: Tile, right: Tile| { left == Tile::Trap && center == Tile::Safe && right == Tile::Safe },
            |left: Tile, center: Tile, right: Tile| { left == Tile::Safe && center == Tile::Safe && right == Tile::Trap }
        ]
    } 

    static TRAP_FUNCTIONS: Lazy<TrapFunctions> = Lazy::new(get_trap_functions);

    fn get_next_row(current_row: &[Tile]) -> Vec<Tile> {
        current_row.iter().enumerate().map(|(i, center)| {
            let left = if i > 0 { current_row.get(i - 1).unwrap() } else { &Tile::Safe };
            let right = current_row.get(i + 1).unwrap_or(&Tile::Safe);

            if TRAP_FUNCTIONS.iter().any(|f| f(*left, *center, *right)) {
                Tile::Trap
            }
            else {
                Tile::Safe
            }
        }).collect()
    }

    pub(crate) fn count_tiles(first_row: &str, row_count: usize) -> usize {
        (0..row_count - 1).fold(vec![parse_first_row(first_row)], |mut v, _| {
            v.push(get_next_row(v.last().unwrap()));
            v
        }).iter().flatten().filter(|x| **x == Tile::Safe).count()
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::like_a_rogue::count_tiles;

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();

    println!("part1: {}", count_tiles(&content, 40));
    println!("part1: {}", count_tiles(&content, 400000));
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_example() {
        use crate::like_a_rogue::count_tiles;
        assert_eq!(38, count_tiles(".^^.^.^^^^", 10));
    }
}
