mod grid_computing {
    use once_cell::sync::Lazy;
    use std::{cmp::Ordering, collections::{BinaryHeap, HashSet}};
    

    #[derive(Clone, Debug, PartialEq, Eq, Hash)]
    pub(crate) struct Node {
        x: i32,
        y: i32,
        size: u32,
        used: u32
    }

    impl Node {
        fn available(&self) -> u32 {
            self.size - self.used
        }
    }

    pub(crate) fn parse(lines: &[&str]) -> Option<Vec<Node>> {
        static REG: Lazy<regex::Regex> = Lazy::new(||regex::Regex::new(r"x(\d+)-y(\d+)\s+(\d+)T\s+(\d+)T").unwrap());

        lines.iter().skip(2).map(|line| {
            let captures = REG.captures(line)?;
            Some(Node{x: captures[1].parse().ok()?, y: captures[2].parse().ok()?, size: captures[3].parse().ok()?, used: captures[4].parse().ok()?})
        }).collect()
    }

    pub(crate) fn viable_pairs(nodes: &[Node]) -> usize {
        nodes.iter().map(|node| {
            nodes.iter().filter(|other| node as *const _ != *other && node.used != 0 && node.used <= other.available()).count()
        }).sum()
    }

    #[derive(Debug, PartialEq, Clone, Eq, Hash)]
    struct NodeState {
        moves: i32,
        state: Vec<Node>,
        empty_node: Node,
        goal: Node,
        unmovable_node: Option<Node>
    }

    impl PartialOrd for NodeState {
        fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
            Some(self.cmp(other))
        }
    }

    impl Ord for NodeState {
        fn cmp(&self, other: &Self) -> Ordering {
            let other_distance = other.moves + (other.goal.x - other.empty_node.x).abs() + (other.goal.y - other.empty_node.y).abs() + other.goal.x + other.goal.y;
            let self_distance = self.moves + (self.goal.x - self.empty_node.x).abs() + (self.goal.y - self.empty_node.y).abs() + self.goal.x + self.goal.y;
            (other_distance).cmp(&self_distance)
        }
    }

    fn init_states(states: &mut BinaryHeap<NodeState>, nodes: &[Node], goal: &Node) {
        nodes.iter().for_each(|node| {
            if node.used == 0 {
                states.push(NodeState{moves: 0, goal: goal.clone(), state: nodes.to_vec(), empty_node: node.clone(), unmovable_node : None});
            }
        })
    }

    const NEIGHBOURS_OFFSET: [(i32, i32); 4] = [(0, 1), (1, 0), (0, -1), (-1, 0)];

    fn append_initial_states(nodes: &[Node], states: &mut BinaryHeap<NodeState>, goal: &Node) {
        for node in nodes{
            if node.used == 0 {
                continue;
            }

            //don't clean up goal node
            if node.x == goal.x && node.y == goal.y {
                continue;
            }

            for offset in NEIGHBOURS_OFFSET {
                if let Some(neighbour) = nodes.iter().find(|neighbour| node.used <= neighbour.available() && node.x + offset.0 == neighbour.x && node.y + offset.1 == neighbour.y) {
                    //don't put crap into goal node
                    if neighbour.x == goal.x && neighbour.y == goal.y {
                        continue;
                    }

                    let new_empty = Node{used : 0,.. *node};
                    let new_state = nodes.iter().map(|new_node| {
                        if new_node.x == node.x && new_node.y == node.y {
                            new_empty.clone()
                        }
                        else if new_node.x == neighbour.x && new_node.y == neighbour.y {
                            Node{used : neighbour.used + node.used,.. *neighbour}
                        }
                        else {
                            new_node.clone()
                        }
                    });
                    states.push(NodeState{moves : 1, goal : goal.clone(), state : new_state.collect(), empty_node : new_empty, unmovable_node : None});
                }
            }
        }
    }

    fn append_states_from_empty_node(nodes: &[Node], states: &mut BinaryHeap<NodeState>, used_states: &mut HashSet<(i32, i32, i32, i32)>, state: &NodeState) {
        for offset in NEIGHBOURS_OFFSET {
            if let Some(neighbour) = nodes.iter().find(|neighbour| neighbour.used <= state.empty_node.size && state.empty_node.x + offset.0 == neighbour.x && state.empty_node.y + offset.1 == neighbour.y) {
                if used_states.contains(&(state.goal.x, state.goal.y, neighbour.x, neighbour.y)) {
                    continue;
                }

                if neighbour.x == state.goal.x && neighbour.y == state.goal.y {
                    continue;
                }

                if let Some(unmovable_node) = &state.unmovable_node {
                    if (unmovable_node.x, unmovable_node.y) == (neighbour.x, neighbour.y) {
                        continue;
                    }
                }

                let new_state = nodes.iter().map(|new_node| {
                    if new_node.x == state.empty_node.x && new_node.y == state.empty_node.y {
                        Node{used : neighbour.used,.. state.empty_node}
                    }
                    else if new_node.x == neighbour.x && new_node.y == neighbour.y {
                        Node{used : 0,.. *neighbour}
                    }
                    else {
                        new_node.clone()
                    }
                });
                used_states.insert((state.goal.x, state.goal.y, neighbour.x, neighbour.y));
                states.push(NodeState{moves : state.moves + 1, state: new_state.collect(), empty_node : neighbour.clone(), .. state.clone()});
            }
        }
    }

    fn reached_goal(state: &NodeState) -> bool {
        NEIGHBOURS_OFFSET.iter().any(|offset| {
            state.empty_node.x + offset.0 == state.goal.x && state.empty_node.y + offset.1 == state.goal.y
        })
    }

    fn create_new_goal_state(state: NodeState) -> NodeState {
        let new_goal = Node{used : state.goal.used,.. state.empty_node};
        let new_empty = Node{used : 0,.. state.goal};
        let new_state = state.state.into_iter().map(|new_node| {
            if new_node.x == state.empty_node.x && new_node.y == state.empty_node.y {
                new_goal.clone()
            }
            else if new_node.x == state.goal.x && new_node.y == state.goal.y {
                new_empty.clone()
            }
            else {
                new_node
            }
        });
        NodeState{moves: state.moves + 1, state: new_state.collect(), goal : new_goal, empty_node : new_empty, unmovable_node: None}
    }

    fn append_new_goal_state(state: NodeState, states: &mut BinaryHeap<NodeState>, used_states: &mut HashSet<(i32, i32, i32, i32)>) {
        assert!(state.unmovable_node.is_none());

        if used_states.contains(&(state.empty_node.x, state.empty_node.y, state.goal.x, state.goal.y)) {
            return;
        }

        let state = create_new_goal_state(state);
        used_states.insert((state.goal.x, state.goal.y, state.empty_node.x, state.empty_node.y));
        states.push(state);    
    }

    fn create_access_node_goal_state(nodes: &[Node], state: NodeState) -> NodeState {
        let mut new_state = create_new_goal_state(state);
        new_state.unmovable_node = Some(new_state.goal);
        new_state.goal = Node{ ..*nodes.iter().find(|node|node.x == 0 && node.y == 0).unwrap() };
        new_state
    }

    pub(crate) fn get_node(nodes: &[Node], initial_goal: &Node) -> i32 {
        let mut states = BinaryHeap::new();
        init_states(&mut states, nodes, initial_goal);
        append_initial_states(nodes, &mut states, initial_goal);

        let mut used_states = states.iter().map(|state| (state.goal.x, state.goal.y, state.empty_node.x, state.empty_node.y)).collect::<HashSet<_>>();

        while let Some(state) = states.pop() {
            if reached_goal(&state) {
                if (state.empty_node.x == 1 && state.empty_node.y == 0) || (state.empty_node.x == 0 && state.empty_node.y == 1) {
                    if (state.goal.x, state.goal.y) == (0, 0) {
                        return state.moves + 2;
                    }
                    else {
                        let access_point_state = create_access_node_goal_state(nodes, state);
                        used_states.insert((access_point_state.goal.x, access_point_state.goal.y, access_point_state.empty_node.x, access_point_state.empty_node.y));
                        states.push(access_point_state);
                    }
                }
                else {
                    append_states_from_empty_node(nodes, &mut states, &mut used_states, &state);
                    append_new_goal_state(state, &mut states, &mut used_states);
                }
            }
            else {
                append_states_from_empty_node(nodes, &mut states, &mut used_states, &state);
            }
        }
        unreachable!("grid is in an invalid state for this algorithm")
    }

    pub(crate) fn get_top_right(nodes: &[Node]) -> i32 {
        get_node(nodes, nodes.iter().max_by_key(|node|(node.x, -node.y)).unwrap())
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::grid_computing::{parse, viable_pairs, get_top_right};

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    let lines: Vec<_> = content.lines().collect();

    let nodes = parse(&lines).expect("failed to parse input file");
    println!("part1: {}", viable_pairs(&nodes));
    println!("part2: {}", get_top_right(&nodes));
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_example() {
        use crate::grid_computing::{parse, get_top_right};

        let lines = ["root@ebhq-gridcenter# df -h",
            "Filesystem            Size  Used  Avail  Use%",
            "/dev/grid/node-x0-y0   10T    8T     2T   80%",
            "/dev/grid/node-x0-y1   11T    6T     5T   54%",
            "/dev/grid/node-x0-y2   32T   28T     4T   87%",
            "/dev/grid/node-x1-y0    9T    7T     2T   77%",
            "/dev/grid/node-x1-y1    8T    0T     8T    0%",
            "/dev/grid/node-x1-y2   11T    7T     4T   63%",
            "/dev/grid/node-x2-y0   10T    6T     4T   60%",
            "/dev/grid/node-x2-y1    9T    8T     1T   88%",
            "/dev/grid/node-x2-y2    9T    6T     3T   66%"];

            let nodes = parse(&lines).unwrap();

            assert_eq!(7, get_top_right(&nodes));
        }

    #[test]
    fn test_example_no_initial_state() {
        use crate::grid_computing::{parse, get_top_right};

        let lines = ["root@ebhq-gridcenter# df -h",
            "Filesystem            Size  Used  Avail  Use%",
            "/dev/grid/node-x0-y0   10T    8T     2T   80%",
            "/dev/grid/node-x0-y1   11T    6T     5T   54%",
            "/dev/grid/node-x0-y2   32T   28T     4T   87%",
            "/dev/grid/node-x1-y0    9T    7T     2T   77%",
            "/dev/grid/node-x1-y1    8T    1T     8T    0%",
            "/dev/grid/node-x1-y2   11T    7T     4T   63%",
            "/dev/grid/node-x2-y0   10T    6T     4T   60%",
            "/dev/grid/node-x2-y1    9T    8T     1T   88%",
            "/dev/grid/node-x2-y2    9T    6T     3T   66%"];

            let nodes = parse(&lines).unwrap();

            assert_eq!(7, get_top_right(&nodes));
        }

    #[test]
    fn test_example_empty_node_low() {
        use crate::grid_computing::{parse, get_top_right};

        let lines = ["root@ebhq-gridcenter# df -h",
            "Filesystem            Size  Used  Avail  Use%",
            "/dev/grid/node-x0-y0   10T    8T     2T   80%",
            "/dev/grid/node-x0-y1   11T    6T     5T   54%",
            "/dev/grid/node-x0-y2   32T   28T     4T   87%",
            "/dev/grid/node-x1-y0    9T    7T     2T   77%",
            "/dev/grid/node-x1-y1    8T    6T     8T    0%",
            "/dev/grid/node-x1-y2   11T    0T     4T   63%",
            "/dev/grid/node-x2-y0   10T    6T     4T   60%",
            "/dev/grid/node-x2-y1    9T    8T     1T   88%",
            "/dev/grid/node-x2-y2    9T    6T     3T   66%"];

            let nodes = parse(&lines).unwrap();

            assert_eq!(8, get_top_right(&nodes));
        }

    #[test]
    fn test_example_extended() {
        use crate::grid_computing::{parse, get_top_right};

        let lines = ["root@ebhq-gridcenter# df -h",
            "Filesystem            Size  Used  Avail  Use%",
            "/dev/grid/node-x0-y0   10T    8T     2T   80%",
            "/dev/grid/node-x0-y1   11T    6T     5T   54%",
            "/dev/grid/node-x0-y2   32T   28T     4T   87%",
            "/dev/grid/node-x1-y0    9T    7T     2T   77%",
            "/dev/grid/node-x1-y1    8T    6T     8T    0%",
            "/dev/grid/node-x1-y2   11T    0T     4T   63%",
            "/dev/grid/node-x2-y0   10T    6T     4T   60%",
            "/dev/grid/node-x2-y1    9T    8T     1T   88%",
            "/dev/grid/node-x2-y2    9T    6T     3T   66%",
            "/dev/grid/node-x3-y0   10T    6T     4T   60%",
            "/dev/grid/node-x3-y1    9T    8T     1T   88%",
            "/dev/grid/node-x3-y2    9T    6T     3T   66%"
            ];

            let nodes = parse(&lines).unwrap();

            assert_eq!(14, get_top_right(&nodes));
        }

        #[test]
        fn test_example_extended_with_goal_blocker() {
            use crate::grid_computing::{parse, get_top_right};
    
            let lines = ["root@ebhq-gridcenter# df -h",
                "Filesystem            Size  Used  Avail  Use%",
                "/dev/grid/node-x0-y0   10T    8T     2T   80%",
                "/dev/grid/node-x0-y1   11T    6T     5T   54%",
                "/dev/grid/node-x0-y2    8T    6T     4T   87%",
                "/dev/grid/node-x1-y0    9T    7T     2T   77%",
                "/dev/grid/node-x1-y1   99T   99T     8T  100%",
                "/dev/grid/node-x1-y2   11T    6T     4T   63%",
                "/dev/grid/node-x2-y0   10T    6T     4T   60%",
                "/dev/grid/node-x2-y1    9T    8T     1T   88%",
                "/dev/grid/node-x2-y2    9T    6T     3T   66%",
                "/dev/grid/node-x3-y0   10T    6T     4T   60%",
                "/dev/grid/node-x3-y1    9T    8T     1T   88%",
                "/dev/grid/node-x3-y2    9T    0T     9T   0%"
                ];
    
                let nodes = parse(&lines).unwrap();
    
                assert_eq!(20, get_top_right(&nodes));
            }
}
