mod an_elephant_named_joseph {
    fn get_lucky_elf_from_list(elves: Vec<usize>) -> usize {
        let elve_count = elves.len(); 
        if elve_count == 1 {
            return elves[0];
        }
        
        get_lucky_elf_from_list(elves.into_iter().step_by(2).skip(elve_count % 2).collect())
    }

    fn remove_single_elf(elves: &mut Vec<usize>) {
        elves.remove(elves.len() / 2);
        let first_elf = elves.remove(0);
        elves.push(first_elf);
    }

    fn get_lucky_elf_from_list_improved_brute_force(mut elves: Vec<usize>) -> usize {
        while elves.len() != 1 {
            remove_single_elf(&mut elves);
        }
        elves[0]
    }

    fn get_lucky_elf_from_list_improved(elves: Vec<usize>) -> usize {
        let elf_count = elves.len();

        if elf_count % 2 == 1 {
            let mut elves = elves;
            remove_single_elf(&mut elves);
            return get_lucky_elf_from_list_improved(elves);
        } 

        let half = elf_count / 2;
        let remove_counter = half / 3;
        if remove_counter == 0 {
            return get_lucky_elf_from_list_improved_brute_force(elves);
        }

        let first = elves.iter().skip(remove_counter * 2).take(half - remove_counter * 2 - 1);
        let middle_removed = elves.iter().skip(half - 1).step_by(3).take(remove_counter + 2);
        let middle_remainder = elves.iter().skip(elf_count - (half % 3));
        let last = elves.iter().take(remove_counter * 2);

        let elves = first.chain(middle_removed).chain(middle_remainder).chain(last);
        return get_lucky_elf_from_list_improved(elves.into_iter().map(|x|*x).collect());
    }

    pub(crate) fn get_lucky_elf(elf_count: usize) -> usize {
        get_lucky_elf_from_list((1..elf_count + 1).collect())
    }

    pub(crate) fn get_lucky_elf_improved(elf_count: usize) -> usize {
        get_lucky_elf_from_list_improved((1..elf_count + 1).collect())
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::an_elephant_named_joseph::{get_lucky_elf, get_lucky_elf_improved};

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    let elve_count = content.parse::<usize>().unwrap();

    println!("part1: {}", get_lucky_elf(elve_count));
    println!("part2: {}", get_lucky_elf_improved(elve_count));
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_example() {
        use crate::an_elephant_named_joseph::{get_lucky_elf, get_lucky_elf_improved};

        assert_eq!(3, get_lucky_elf(5));
        assert_eq!(2, get_lucky_elf_improved(5));
        assert_eq!(9, get_lucky_elf_improved(9));
        assert_eq!(1, get_lucky_elf_improved(10));
        assert_eq!(17, get_lucky_elf_improved(98));
        assert_eq!(18, get_lucky_elf_improved(99));
        assert_eq!(19, get_lucky_elf_improved(100));
        assert_eq!(20, get_lucky_elf_improved(101));
    }
}
