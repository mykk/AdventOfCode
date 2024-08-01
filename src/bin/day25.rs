mod safe_cracking {
    use std::collections::HashMap;
    use aoc_2016::utils::assembunny;

    pub(crate) fn find_transmision_index(assembly: &[assembunny::Instructions]) -> i32 {
        (0..).find(|a| {
            let mut previous_output = 1;
            let mut seen_states = Vec::new();
            let mut found_signal = false;
            assembunny::execute_assembly(assembly, [("a".to_string(), *a)].into(), |values, output_val|{
                if (output_val != 0 && output_val != 1) || (output_val == previous_output) {
                    return true;
                }

                previous_output = output_val;

                if seen_states.iter().any(|seen_state: &HashMap<String, i32>| { seen_state["a"] == values["a"] && seen_state["b"] == values["b"] && seen_state["c"] == values["c"] && seen_state["d"] == values["d"] }) {
                    found_signal = true;
                    true
                }
                else {
                    seen_states.push(values.clone());
                    false    
                }
            });

            found_signal
        }).unwrap()
    } 
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use aoc_2016::utils::assembunny;

    use crate::safe_cracking::find_transmision_index;

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();

    let lines = content.lines().collect::<Vec<_>>();
    let assembly = assembunny::parse_assembly(&lines);

    let value = find_transmision_index(&assembly);
    println!("part1: {}", value);
}
