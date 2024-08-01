
use once_cell::sync::Lazy;
use regex::Regex;
use std::collections::HashMap;

#[derive(Debug, PartialEq, Clone, Eq)]
pub enum Instructions {
    Cpy(String, String),
    Inc(String),
    Dec(String),
    Jnz(String, String),
    Tgl(String),
    Mul(String, String, String),
    Out(String),
    Pass(),
}

fn extract_capture(captures: &regex::Captures, index: usize) -> String {
    captures.get(index).map_or_else(
        || panic!("Capture group {} not found", index),
        |capture| capture.as_str().to_string())
}

fn parse_instruction(line: &str) -> Option::<Instructions> {
    static PATTERNS: Lazy<Vec<Regex>> = Lazy::new(|| vec![
        Regex::new(r"cpy (-*\w+) (\w+)").unwrap(),
        Regex::new(r"inc (\w+)").unwrap(),
        Regex::new(r"dec (\w+)").unwrap(),
        Regex::new(r"jnz (-*\w+) (-*\w+)").unwrap(),
        Regex::new(r"tgl (-*\w+)").unwrap(),
        Regex::new(r"out (-*\w+)").unwrap(),
    ]);

    PATTERNS.iter().enumerate().find_map(|(i, pattern)| {
        if let Some(captures) = pattern.captures(line) {
            let instruction = match i {
                0 => Instructions::Cpy(extract_capture(&captures, 1), extract_capture(&captures, 2)),
                1 => Instructions::Inc(extract_capture(&captures, 1)),
                2 => Instructions::Dec(extract_capture(&captures, 1)),
                3 => Instructions::Jnz(extract_capture(&captures, 1), extract_capture(&captures, 2)),
                4 => Instructions::Tgl(extract_capture(&captures, 1)),
                5 => Instructions::Out(extract_capture(&captures, 1)),
                _ => unreachable!(),
            };
            return Some(instruction);
        }
        None
    })
}

pub fn parse_assembly(lines: &[&str]) -> Vec<Instructions> {    
    lines.iter().fold(Vec::new(), |mut vec, line| {
        if let Some(instruction) = parse_instruction(line) {
            vec.push(instruction);
        }
        else {
            panic!("unrecognized instruction!");
        }
        vec
    })
}

fn get_value(values: &HashMap<String, i32>, val: &str) -> i32 {
    if let Some(val) = val.parse().ok() {
        val
    }
    else {
        *values.get(val).unwrap_or(&0)
    }
}

fn get_multiplication_result_reg<'a>(instructions: &[&'a Instructions]) -> Option<&'a String> {
    instructions.iter().skip(3).take(2).find_map(|instruction|
        if let Instructions::Inc(reg) = instruction {
            Some(reg)
        } else {
            None
        }
    )
}

fn get_optimized_instructions(instructions: &[&Instructions]) -> Option<Vec<Instructions>> {
    let iter = instructions.iter();
    let init_vals = iter.clone().take(3).map(|instruction| match instruction {
        Instructions::Cpy(val, reg) => Some((val, reg)),
        _ => None
    }).collect::<Option<Vec<_>>>()?;

    let result_reg = get_multiplication_result_reg(instructions)?;

    let mut multipliers = init_vals.iter().filter(|(_, reg)|*reg != result_reg);
    let multiplier1 = multipliers.next()?.1;
    let multiplier2 = multipliers.next()?.1;
    if *multiplier1 == *multiplier2 {
        return None;
    }

    if iter.clone().skip(3).take(2).all(|instruction| matches!(instruction, Instructions::Inc(reg) if result_reg != reg)) {
        return None;
    }

    let current_multiplier = iter.clone().skip(3).take(2).find(|instruction| matches!(instruction, Instructions::Dec(multiplier) if multiplier == multiplier1 || multiplier == multiplier2))?;
    let (multiplier1, multiplier2) = if matches!(current_multiplier, Instructions::Dec(multiplier) if multiplier == multiplier1) {(multiplier1, multiplier2)} else { (multiplier2, multiplier1) };

    
    let mut iter = iter.skip(5);
    if matches!(iter.next()?, Instructions::Jnz(reg, value) if reg != multiplier1 || value != "-2") { return None; }
    if matches!(iter.next()?, Instructions::Dec(reg) if reg != multiplier2) { return None; }
    if matches!(iter.next()?, Instructions::Jnz(reg, value) if reg != multiplier2 || value != "-5") { return None; }

    let cpy_instructions = 
        instructions.iter().take(3).map(|x|(**x).clone()).
        chain([Instructions::Mul(multiplier1.clone(), multiplier2.clone(), result_reg.clone())]).
        chain([Instructions::Pass(), Instructions::Pass(), Instructions::Pass(), Instructions::Pass()]);
    Some(cpy_instructions.collect())
}

fn optimize_assembly(assembly: &[Instructions]) -> Vec<Instructions> {
    let mut index = 0;

    let mut optimized_assembly = vec![];
    while index < assembly.len() {
        if let Some(instructions) = get_optimized_instructions(&assembly.iter().skip(index).collect::<Vec<_>>()) {
            index += instructions.len();
            optimized_assembly.extend(instructions.into_iter());    
        }
        else {
            optimized_assembly.push(assembly[index].clone());
            index += 1;
        }
    }

    optimized_assembly 
}

pub fn execute_assembly<OutputF>(assembly: &[Instructions], initial_values: HashMap<String, i32>, mut output: OutputF) -> HashMap<String, i32> 
where OutputF: FnMut(&HashMap<String, i32>, i32) -> bool 
{
    let mut index: usize = 0;
    let mut values = initial_values;

    let mut assembly = optimize_assembly(assembly); 
    while index < assembly.len() {
        index = match &assembly[index] {
            Instructions::Cpy(value, registry) => {
                if registry.parse::<i32>().is_err() {
                    values.insert(registry.clone(), get_value(&values, value));
                }
                index + 1    
            }
            Instructions::Inc(registry) => {
                values.entry(registry.clone()).and_modify(|value| *value += 1);
                index + 1
            }
            Instructions::Dec(registry) => {
                values.entry(registry.clone()).and_modify(|value| *value -= 1);
                index + 1
            }
            Instructions::Jnz(value, jump) => {
                if get_value(&values, value) != 0 {
                    let jump = get_value(&values, jump);
                    if jump < 0 { index - (-jump) as usize } else { index + jump as usize }
                }
                else {
                    index + 1
                }
            }
            Instructions::Tgl(value) => {
                let value = get_value(&values, value);
                let tgl_index = if value < 0 { index - (-value) as usize } else { index + value as usize };
                if let Some(instruction) = assembly.get(tgl_index) {
                    assembly[tgl_index] = match instruction {
                        Instructions::Inc(val) => Instructions::Dec(val.clone()),
                        Instructions::Dec(val) => Instructions::Inc(val.clone()),
                        Instructions::Tgl(val) => Instructions::Inc(val.clone()),
                        Instructions::Jnz(val1, val2) => Instructions::Cpy(val1.clone(), val2.clone()),
                        Instructions::Cpy(val1, val2) => Instructions::Jnz(val1.clone(), val2.clone()),
                        _ => panic!("optimized instruction toggled")
                    }
                }
                index + 1
            }
            Instructions::Pass() => panic!("multiplication should jump over these intstructions and no other instruction should jump here"),
            Instructions::Mul(val1, val2, registry) => {
                let original_value = get_value(&values, registry); 
                let val1 = get_value(&values, val1);
                let val2 = get_value(&values, val2);
                values.insert(registry.clone(), val1 * val2 + original_value);
                index + 5
            }
            Instructions::Out(val) => {
                let terminate = output(&values, get_value(&values, val));
                if terminate {
                    return values;
                }
                index + 1
            }
        };
    }
    values
}
