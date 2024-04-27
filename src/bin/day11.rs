mod radioisotope_thermoelectric_generators {
    use std::{cmp::Ordering, collections::HashSet, rc::Rc};

    use once_cell::sync::Lazy;
    use std::collections::BinaryHeap;
    use regex::bytes::Regex;

    #[derive(Debug, PartialEq, Clone, Eq, Hash)]
    pub(crate) enum ThermoelectricComponent {
        Microchip(Rc<String>),
        Generator(Rc<String>)
    }
    
    #[derive(Debug, PartialEq, Eq, Clone)]
    struct State {
        floor: usize,
        state: Vec<HashSet<ThermoelectricComponent>>,
        count: i32  
    }

    impl PartialOrd for State {
        fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
            Some(self.cmp(other))
        }
    }

    impl Ord for State {
        fn cmp(&self, other: &Self) -> Ordering {
            if self.count != other.count {
                return if self.count < other.count { Ordering::Greater } else { Ordering::Less } 
            }

            for (i, x) in self.state.iter().enumerate().rev() {
                if x.len() < other.state[i].len() {
                    return Ordering::Less;
                }

                if x.len() > other.state[i].len() {
                    return Ordering::Greater;
                }
            }
            Ordering::Equal
        }
    }

    fn parse_components<F>(line: &str, regex: &Regex, constructor: F) -> Vec<ThermoelectricComponent>
    where F: Fn(Rc<String>) -> ThermoelectricComponent {
        regex.captures_iter(line.as_ref())
            .filter_map(|cap| cap.get(1))
            .map(|group| {
                constructor(Rc::new(String::from_utf8(group.as_bytes().iter().map(|x| *x).collect()).unwrap()))
            }).collect()
    }
    
    pub(crate) fn parse_input(lines: &[&str]) -> Vec<HashSet<ThermoelectricComponent>>{
        static MICROCHIP_GROUP_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"(\w+)-compatible microchip").unwrap());
        static GENERATOR_GROUP_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"(\w+) generator").unwrap());

        lines.iter().map(|line| {
            let microchips = parse_components(line, &MICROCHIP_GROUP_REGEX, ThermoelectricComponent::Microchip);
            let generators = parse_components(line, &GENERATOR_GROUP_REGEX, ThermoelectricComponent::Generator);
            microchips.into_iter().chain(generators.into_iter()).collect()
        }).collect()
    }

    fn is_finished(current_state: &[HashSet<ThermoelectricComponent>]) -> bool {
        current_state.iter().take(current_state.len() - 1).all(|floor| floor.is_empty())
    }

    //microchip cannot be on the same floor with another generator if there's no matching generator
    fn is_valid(current_state: &[HashSet<ThermoelectricComponent>]) -> bool {
        !current_state.iter().any(|floor|floor.iter().any(|component|{
            if let ThermoelectricComponent::Microchip(microchip) = component {
                return !floor.iter().any(|component| matches!(component, ThermoelectricComponent::Generator(generator) if generator == microchip)) && 
                        floor.iter().any(|component| matches!(component, ThermoelectricComponent::Generator(_)))
            }
            false
        }))
    }

    fn is_used_state(state: &[HashSet<ThermoelectricComponent>], current_floor: usize, count: i32, previous_states: &BinaryHeap<State>) -> bool {
        previous_states.iter().any(|x| x.floor == current_floor && 
            x.count <= count && 
            x.state.iter().enumerate().all(|(i, components)| {
                components.len() == state[i].len() && 
                components.iter().filter(|x|matches!(x, ThermoelectricComponent::Generator(_))).count() == 
                state[i].iter().filter(|x|matches!(x, ThermoelectricComponent::Generator(_))).count()
            }))
    }

    fn get_new_state_template(current_state: &[HashSet<ThermoelectricComponent>], current_floor: usize) -> Vec<HashSet<ThermoelectricComponent>> {
        current_state.iter().enumerate().fold(Vec::new(), |mut vec, (floor, components)| {
            if floor == current_floor {
                vec.push(HashSet::new());
            }
            else {
                vec.push(components.clone());
            }
            vec
        })
    }

    fn move_components(current_state: &[HashSet<ThermoelectricComponent>], current_floor: usize, new_floor: usize) -> Vec<Vec<HashSet<ThermoelectricComponent>>> {
        let new_state_template = get_new_state_template(current_state, current_floor);

        current_state[current_floor].iter().enumerate().fold(Vec::new(), |mut vec, (index, component)| {
            //move pairs
            for other_component in current_state[current_floor].iter().skip(index + 1) {
                let mut new_state = new_state_template.clone();
                new_state[current_floor] = current_state[current_floor].iter().filter(|x| *x != component && *x != other_component).map(|x|x.clone()).collect();
                new_state[new_floor].insert(component.clone());
                new_state[new_floor].insert(other_component.clone());                
                vec.push(new_state)    
            }

            //move singles
            let mut new_state = new_state_template.clone();
            new_state[current_floor] = current_state[current_floor].iter().filter(|x| *x != component).map(|x|x.clone()).collect();
            new_state[new_floor].insert(component.clone());
            vec.push(new_state);

            vec
        })
    }

    fn move_and_filter_components(current_state: &State, new_floor: usize, used_states: &BinaryHeap<State>, state_stack: &BinaryHeap<State>) -> BinaryHeap<State> {
        move_components(&current_state.state, current_state.floor, new_floor).into_iter()
        .filter(|new_state| 
            is_valid(new_state) && 
            !is_used_state(new_state, new_floor, current_state.count + 1, used_states) &&
            !is_used_state(new_state, new_floor, current_state.count + 1, state_stack))
        .map(|state|State{state, floor: new_floor, count: current_state.count + 1})
        .collect::<BinaryHeap<_>>()
    }

    fn move_components_up(current_state: &State, used_states: &BinaryHeap<State>, state_stack: &BinaryHeap<State>) -> BinaryHeap<State> {
        move_and_filter_components(current_state, current_state.floor + 1, used_states, state_stack)
    }

    fn move_components_down(current_state: &State, used_states: &BinaryHeap<State>, state_stack: &BinaryHeap<State>) -> BinaryHeap<State> {
        move_and_filter_components(current_state, current_state.floor - 1, used_states, state_stack)
    }

    fn move_components_bfs(starting_state: &[HashSet<ThermoelectricComponent>]) -> i32 {
        let mut state_stack = BinaryHeap::from([State{state: starting_state.to_vec(), floor: 0, count : 0}]);
        let mut used_states = BinaryHeap::new();

        while let Some(current_state) = state_stack.pop() {
            if is_finished(&current_state.state) {
                return current_state.count;
            }

            if current_state.floor < current_state.state.len() - 1 {
                state_stack.append(&mut move_components_up(&current_state, &used_states, &state_stack));
            }

            if current_state.floor > 0 {
                state_stack.append(&mut move_components_down(&current_state, &used_states, &state_stack));
            }

            used_states.push(current_state);
        }
        unreachable!();
    }

    pub(crate) fn find_optimal_move_pattern(compoenents: &[HashSet<ThermoelectricComponent>]) -> i32 {
        move_components_bfs(&compoenents)
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();

    let lines = &content.lines().collect::<Vec<_>>();
    let components = radioisotope_thermoelectric_generators::parse_input(&lines);
    let moved_in = radioisotope_thermoelectric_generators::find_optimal_move_pattern(&components);

    println!("part1: {}", moved_in);
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_example() {
        let example_str = r"The first floor contains a hydrogen-compatible microchip and a lithium-compatible microchip.\n\
        The second floor contains a hydrogen generator.\n\ 
        The third floor contains a lithium generator.\n\
        The fourth floor contains nothing relevant.\n";

        let lines = example_str.lines().collect::<Vec<_>>();
        let components = radioisotope_thermoelectric_generators::parse_input(&lines);

        let moved_in = radioisotope_thermoelectric_generators::find_optimal_move_pattern(&components);
        assert_eq!(moved_in, 11);
    }
}
