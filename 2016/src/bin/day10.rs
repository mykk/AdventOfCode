use std::collections::HashMap;

mod balance_bots {
    use std::collections::HashMap;
    use regex::Regex;
    use once_cell::sync::Lazy;

    struct Bot {
        chip: Option<u32>,

        give_low_to: Option<u32>,
        give_high_to: Option<u32>,

        discard_low_to: Option<u32>,
        discard_high_to: Option<u32>,
    }

    impl Bot {
        fn new(give_low_to: Option<u32>, give_high_to: Option<u32>, discard_low_to: Option<u32>, discard_high_to: Option<u32>) -> Self {
            Bot{ chip: None, give_low_to, give_high_to, discard_low_to, discard_high_to}
        }

        //bot releases both chips as soon as it gets 2nd.
        fn receive_chip(&mut self, chip: u32) -> Option<(u32, u32)> {
            match self.chip {
                Some(val) => {
                    let low = val.min(chip);
                    let high = val.max(chip);
                    self.chip = None;
        
                    Some((low, high))
                }
                None => {
                    self.chip = Some(chip);
                    None
                }
            } 
        }
    }

    fn add_chips_to_discard(bot: &mut Bot, low_chip: u32, high_chip: u32, discarded: &mut HashMap<u32, u32>) {
        if let Some(discard_low_to) = bot.discard_low_to {
            discarded.insert(discard_low_to, low_chip);
        }
        
        if let Some(discard_high_to) = bot.discard_high_to {
            discarded.insert(discard_high_to, high_chip);
        }
    }

    fn balance_bot(bot_number: u32, chip: u32, bots: &mut HashMap<u32, Bot>, discarded: &mut HashMap<u32, u32>, search_val1: u32, search_val2: u32) -> Option<u32>{
        let bot = bots.get_mut(&bot_number).unwrap();

        let mut bot_found = None;
        if let Some((low_chip, high_chip)) = bot.receive_chip(chip) {
            add_chips_to_discard(bot, low_chip, high_chip, discarded);

            if low_chip == search_val1.min(search_val2) && high_chip == search_val1.max(search_val2) {
                bot_found = Some(bot_number);
            }

            let give_low_to = bot.give_low_to;
            let give_high_to = bot.give_high_to;

            if let Some(bot_number) = give_low_to.and_then(|bot|balance_bot(bot, low_chip, bots, discarded, search_val1, search_val2)) {
                bot_found =  Some(bot_number);
            }

            if let Some(bot_number) = give_high_to.and_then(|bot|balance_bot(bot, high_chip, bots, discarded, search_val1, search_val2)) {
                bot_found =  Some(bot_number);
            }
        }
        bot_found
    }

    fn balance_bots(bots: &mut HashMap<u32, Bot>, initial_values: &[(u32, u32)], discarded: &mut HashMap<u32, u32>, search_val1: u32, search_val2: u32) -> Option<u32> {
        initial_values.iter().fold(None, |mut bot_found, (bot_number, chip)|{
            if let Some(bot) = balance_bot(*bot_number, *chip, bots, discarded, search_val1, search_val2) {
                bot_found = Some(bot)
            }
            bot_found
        })
    }

    fn parse_bot_action(caps: &regex::Captures) -> (u32, Bot) {
        let bot_number = caps[1].parse::<u32>().unwrap();

        let give_low_to = if caps[2] == *"bot" { caps[3].parse::<u32>().ok() } else { None };
        let give_high_to = if caps[4] == *"bot" { caps[5].parse::<u32>().ok() } else { None };

        let discard_low_to = if caps[2] == *"output" { caps[3].parse::<u32>().ok() } else { None };
        let discard_high_to = if caps[4] == *"output" { caps[5].parse::<u32>().ok() } else { None };

        (bot_number, Bot::new(give_low_to, give_high_to, discard_low_to, discard_high_to))
    }
    
    fn parse_instructions<T>(instructions: &[T]) -> (HashMap<u32, Bot>, Vec<(u32, u32)>)
    where T: AsRef<str> {
        static BOT_ACTION_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"bot (\d+) gives low to (bot|output) (\d+) and high to (bot|output) (\d+)").unwrap());
        static BOT_INIT_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"value (\d+) goes to bot (\d+)").unwrap());

        instructions.iter()
            .fold((HashMap::new(), Vec::new()), |(mut bots, mut initial_vals), line| {
                if let Some(caps) = BOT_INIT_REGEX.captures(line.as_ref()) {
                    initial_vals.push((caps[2].parse::<u32>().unwrap(), caps[1].parse::<u32>().unwrap()));
                }
                else if let Some(caps) = BOT_ACTION_REGEX.captures(line.as_ref()) {
                    let (bot_number, bot) = parse_bot_action(&caps);
                    bots.insert(bot_number, bot);
                }

                (bots, initial_vals)    
            })
    }

    pub(crate) fn find_bot_handling_values<T>(instructions: &[T], search_val1: u32, search_val2: u32, discarded: &mut HashMap<u32, u32>) -> Option<u32>
    where T: AsRef<str> {
        let (mut bots, initial_vals) = parse_instructions(instructions);
        balance_bots(&mut bots, &initial_vals, discarded, search_val1, search_val2)
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    let instructions = content.lines().collect::<Vec<_>>();

    let mut discarded = HashMap::new();
    let bot = balance_bots::find_bot_handling_values(&instructions, 61, 17, &mut discarded).unwrap();
    println!("part1: {}", bot);
    println!("part2: {}", discarded[&0] * discarded[&1] * discarded[&2]);
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_example() {
        let instructions = format!("{}\n{}\n{}\n{}\n{}\n{}", 
            "value 5 goes to bot 2",
            "bot 2 gives low to bot 1 and high to bot 0", 
            "value 3 goes to bot 1",
            "bot 1 gives low to output 1 and high to bot 0", 
            "bot 0 gives low to output 2 and high to output 0",
            "value 2 goes to bot 2");

        let mut discard_vals = HashMap::new();
        let bot = balance_bots::find_bot_handling_values(&instructions.lines().collect::<Vec<_>>(), 2, 5, &mut discard_vals).unwrap();
        assert_eq!(2, bot);
    }
}
