pub mod door_hacking {
    use md5;

    pub struct HackTheDoor<T> {
        hack_algorithm: Box<dyn HackAlgorithm<T>>
    }

    impl<T> HackTheDoor<T> {
        pub fn new(hack_algorithm: Box<dyn HackAlgorithm<T>>) -> Self {
            Self{hack_algorithm}
        }

        fn get_initial_state(&self) -> T {
            self.hack_algorithm.get_initial_state()
        } 

        fn hacked(&self, password: &T) -> bool {
            self.hack_algorithm.hacked(password)
        }

        fn hack(&self, password: &mut T, digest: &str) {
            self.hack_algorithm.hack(password, digest);
        }

        fn extract(&self, password: &T) -> String {
            self.hack_algorithm.extract(password)
        }

        fn hack_from_cache(&self, cache: &Vec<(String, usize)>) -> (T, usize) {
            cache.iter().try_fold((self.get_initial_state(), 0), |(mut pass, _), (digest, position)| {
                if self.hacked(&pass) {
                    return None;
                }
                self.hack(&mut pass, digest);
                Some((pass, position + 1))
            }).unwrap_or((self.get_initial_state(), 0))
        }
        
        pub fn hack_the_door(&self, door_id: &str, cache: &mut Vec<(String, usize)>) -> String {
            let (mut password, starting_pos) = self.hack_from_cache(cache);
            if self.hacked(&password) {
                return self.extract(&password);
            }

            for i in starting_pos.. {
                let digest = format!("{:x}", md5::compute(door_id.to_string() + i.to_string().as_ref()));
                if !digest.starts_with("00000") {
                    continue;
                }
                self.hack(&mut password, &digest);
                if !cache.iter().any(|(_, pos)| *pos > i) {
                    cache.push((digest, i));
                }

                if self.hacked(&password) {
                    return self.extract(&password)
                }
            }
            unreachable!();
        }
    }

    pub trait HackAlgorithm<T> {
        fn get_initial_state(&self) -> T;
        fn hacked(&self, password: &T) -> bool;
        fn hack(&self, password: &mut T, digest: &str);
        fn extract(&self, password: &T) -> String;
    }

    pub struct HackFirstDoor;
    impl HackAlgorithm<String> for HackFirstDoor {
        fn get_initial_state(&self) -> String {
            String::new()
        }

        fn hacked(&self, password: &String) -> bool {
            password.len() == 8
        }

        fn hack(&self, password: &mut String, digest: &str) {
            password.push(digest.chars().nth(5).unwrap());
        }

        fn extract(&self, password: &String) -> String {
            password.clone()
        }
    }

    pub struct HackSecondDoor;
    impl HackAlgorithm<[Option<char>; 8]> for HackSecondDoor {
        fn get_initial_state(&self) -> [Option<char>; 8] {
            [None; 8]
        }

        fn hacked(&self, password: &[Option<char>; 8]) -> bool {
            password.iter().find(|c| c.is_none()).is_none()
        }

        fn hack(&self, password: &mut [Option<char>; 8], digest: &str) {
            if let Ok(i) = digest.chars().nth(5).unwrap().to_string().parse::<usize>() {
                if i < 8 && password[i].is_none() {
                    password[i] = Some(digest.chars().nth(6).unwrap()); //unwrap to crash if not found
                }
            }
        }

        fn extract(&self, password: &[Option<char>; 8]) -> String {
            password.iter().map(|c| c.unwrap()).collect()
        }
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::door_hacking::HackTheDoor;
    
    let door_id = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    let mut cache = Vec::new();

    let hack = HackTheDoor::new(Box::new(door_hacking::HackFirstDoor{}));
    println!("part 1: {}", hack.hack_the_door(&door_id, &mut cache));

    let hack = HackTheDoor::new(Box::new(door_hacking::HackSecondDoor{}));
    println!("part 2: {}", hack.hack_the_door(&door_id, &mut cache));
}

#[cfg(test)]
mod tests {
    use crate::door_hacking;
    use std::time::Instant;

    #[test]
    fn hack_first_door() {
        let mut cache = Vec::new();
        let hack = door_hacking::HackTheDoor::new(Box::new(door_hacking::HackFirstDoor{}));
        assert_eq!("18f47a30", hack.hack_the_door("abc", &mut cache));
    }

    #[test]
    fn hack_second_door_no_cache() {
        let mut cache = Vec::new();
        let hack = door_hacking::HackTheDoor::new(Box::new(door_hacking::HackSecondDoor{}));
        assert_eq!("05ace8e3", hack.hack_the_door("abc", &mut cache));
    }

    #[test]
    fn hack_second_door_with_first_door_cache() {
        let mut cache = Vec::new();
        let hack = door_hacking::HackTheDoor::new(Box::new(door_hacking::HackFirstDoor{}));
        assert_eq!("18f47a30", hack.hack_the_door("abc", &mut cache));

        let hack = door_hacking::HackTheDoor::new(Box::new(door_hacking::HackSecondDoor{}));
        assert_eq!("05ace8e3", hack.hack_the_door("abc", &mut cache));
    }

    #[test]
    fn hack_second_door_with_first_door_cache_twice() {
        let mut cache = Vec::new();
        let hack = door_hacking::HackTheDoor::new(Box::new(door_hacking::HackFirstDoor{}));
        assert_eq!("18f47a30", hack.hack_the_door("abc", &mut cache));

        let hack = door_hacking::HackTheDoor::new(Box::new(door_hacking::HackFirstDoor{}));
        assert_eq!("18f47a30", hack.hack_the_door("abc", &mut cache));

        let hack = door_hacking::HackTheDoor::new(Box::new(door_hacking::HackSecondDoor{}));
        assert_eq!("05ace8e3", hack.hack_the_door("abc", &mut cache));
    }

    #[test]
    fn benchmark_both_doors() {
        let mut cache = Vec::new();

        let start = Instant::now();
        let hack = door_hacking::HackTheDoor::new(Box::new(door_hacking::HackFirstDoor{}));
        assert_eq!("18f47a30", hack.hack_the_door("abc", &mut cache));
        let initial_execution_time = start.elapsed();

        let start = Instant::now();
        let hack = door_hacking::HackTheDoor::new(Box::new(door_hacking::HackFirstDoor{}));
        assert_eq!("18f47a30", hack.hack_the_door("abc", &mut cache));
        let cached_execution_time = start.elapsed() / 10000; //at least 10000 times faster
        assert!(cached_execution_time < initial_execution_time);

        let start = Instant::now();
        let hack = door_hacking::HackTheDoor::new(Box::new(door_hacking::HackSecondDoor{}));
        assert_eq!("05ace8e3", hack.hack_the_door("abc", &mut cache));
        let initial_execution_time = start.elapsed();

        let start = Instant::now();
        let hack = door_hacking::HackTheDoor::new(Box::new(door_hacking::HackSecondDoor{}));
        assert_eq!("05ace8e3", hack.hack_the_door("abc", &mut cache));
        let cached_execution_time = start.elapsed() / 10000; //at least 10000 times faster
        assert!(cached_execution_time < initial_execution_time);
    }
}