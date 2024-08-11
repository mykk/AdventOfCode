mod air_duct_spelunking {
    use std::collections::{BinaryHeap, HashMap, HashSet};
    use std::cmp::Reverse;

    fn get_next_state(line_number: usize, line_delta: isize, char_position: usize, char_delta: isize, lines: &[&str], visited_positions: &HashSet<(usize, usize)>) -> Option::<(char, usize, usize)> {
        let line_number = line_number.checked_add_signed(line_delta)?;
        let char_position = char_position.checked_add_signed(char_delta)?;

        if visited_positions.contains(&(line_number, char_position)) {
            return None;
        } 

        let c = lines.get(line_number)?.chars().nth(char_position).filter(|c| *c != '#')?;
        Some((c, line_number, char_position))
    }

    fn parse_node_distances(line_number: usize, char_position: usize, lines: &[&str]) -> HashMap::<char, u32> {
        const DELTAS: [(isize, isize); 4] = [(1, 0), (-1, 0), (0, 1), (0, -1)];

        let mut node_distances = HashMap::new();
        let mut states = BinaryHeap::new();
        let mut visited_positions = HashSet::new();

        states.push(Reverse((0, line_number, char_position)));
        visited_positions.insert((line_number, char_position));

        while let Some(Reverse((distance, line_number, char_position))) = states.pop() {
            for (line_delta, char_delta) in DELTAS {
                if let Some((next_c, next_line, next_char_pos)) = get_next_state(line_number, line_delta, char_position, char_delta, lines, &visited_positions) {
                    visited_positions.insert((next_line, next_char_pos));
                    states.push(Reverse((distance + 1, next_line, next_char_pos)));
    
                    if next_c != '.' && !node_distances.contains_key(&next_c) {
                        node_distances.insert(next_c, distance + 1);
                    }
                }
            }
        }
        
        node_distances
    }

    pub(crate) fn parse_map(lines: &[&str]) -> HashMap<char, HashMap<char, u32>> {
        lines.iter().enumerate().flat_map(|(line_number, line)| {
            line.chars().enumerate().filter_map(move |(char_position, c)| {
                if c != '.' && c != '#' {
                    Some((c, parse_node_distances(line_number, char_position, lines)))
                } else {
                    None
                }
            })
        }).collect()
    }
    
    fn on_final_node(node: char, graph: &HashMap::<char, HashMap::<char, u32>>, distance: u32, return_to_zero: bool) -> u32 {
        if !return_to_zero {
            distance
        } else {
            distance + graph[&node][&'0']
        }
    }

    fn find_shortest_route_inner(node: char, graph: &HashMap::<char, HashMap::<char, u32>>, visited: HashSet<char>, distance: u32, return_to_zero: bool) -> u32 {
        graph.get(&node).unwrap().iter().filter_map(|(current_node, current_distance)| {
            if visited.contains(current_node) {
                return None;
            }

            Some(find_shortest_route_inner(*current_node, graph, visited.iter().chain([current_node]).map(|c|*c).collect(), current_distance + distance, return_to_zero))
        }).min().unwrap_or_else(||on_final_node(node, graph, distance, return_to_zero))
    }

    pub(crate) fn find_shortest_route(starting_node: char, graph: &HashMap::<char, HashMap::<char, u32>>, return_to_zero: bool) -> u32 {
        find_shortest_route_inner(starting_node, graph, [starting_node].into(), 0, return_to_zero)
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::air_duct_spelunking::{parse_map, find_shortest_route};

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();

    let lines = content.lines().collect::<Vec<_>>();
    let graph = parse_map(&lines);
    let shortest_route = find_shortest_route('0', &graph, false);
    println!("part1: {}", shortest_route);

    let shortest_route = find_shortest_route('0', &graph, true);
    println!("part2: {}", shortest_route);
}

#[cfg(test)]
mod tests {
    use crate::air_duct_spelunking::{parse_map, find_shortest_route};

    #[test]
    fn test_example() {
        let duct_map = 
"###########
#0.1.....2#
#.#######.#
#4.......3#
###########";
    let lines = duct_map.lines().collect::<Vec<_>>();
    let graph = parse_map(&lines);
    assert_eq!(14, find_shortest_route('0', &graph, false));
    assert_eq!(20, find_shortest_route('0', &graph, true));
    }
}
