pub mod door_hacking {
    use std::{collections::HashSet, sync::{atomic::{AtomicBool, AtomicUsize, Ordering}, Mutex}};

    use md5;

    pub struct HackTheDoor {
        hack_algorithm: Mutex<Box<dyn HackAlgorithm>>
    }

    impl HackTheDoor {
        pub fn new(hack_algorithm: Box<dyn HackAlgorithm>) -> Self {
            Self{hack_algorithm: Mutex::new(hack_algorithm)}
        }

        pub fn hack_the_door(&self, door_id: &str, cache: &mut Vec<(String, usize)>) -> String {
            const THREAD_COUNT: usize = 8;
        
            if self.hack_algorithm.lock().unwrap().hacked(cache) {
                return self.hack_algorithm.lock().unwrap().extract(cache);
            }

            let index = AtomicUsize::new(cache.last().map(|(_, pos)| *pos + 1).unwrap_or(0));
            let terminate = AtomicBool::new(false);
            let thread_cache: Mutex<&mut Vec<(String, usize)>> = Mutex::new(cache);

            std::thread::scope(|s| {
                for _ in 0..THREAD_COUNT {
                    s.spawn(|| {
                        loop {
                            if terminate.load(Ordering::Relaxed) {
                                return;
                            }
    
                            let i = index.fetch_add(1, Ordering::Relaxed);
        
                            let digest = format!("{:x}", md5::compute(door_id.to_string() + i.to_string().as_ref()));
                            if digest.starts_with("00000") {
                                let mut thread_cache = thread_cache.lock().unwrap();
                                thread_cache.push((digest, i));
                                thread_cache.sort_by_key(|(_, pos)|{*pos});
    
                                if self.hack_algorithm.lock().unwrap().hacked(&thread_cache) {
                                    terminate.store(true, Ordering::Relaxed);
                                    return;
                                }
                            }
                        }
                    });
                }
            });
            return self.hack_algorithm.lock().unwrap().extract(cache)
        }
        
    }
    
    pub trait HackAlgorithm : Send {
        fn hacked(&self, cache: &[(String, usize)]) -> bool;
        fn extract(&self, cache: &[(String, usize)]) -> String;
    }

    pub struct HackFirstDoor;
    impl HackAlgorithm for HackFirstDoor {
        fn hacked(&self, cache: &[(String, usize)]) -> bool {
            cache.len() == 8
        }

        fn extract(&self, cache: &[(String, usize)]) -> String {
            let mut password = String::new();
            for i in 0..8 {
                password.push(cache[i].0.chars().nth(5).unwrap());
            }
            password
        }
    }

    pub struct HackSecondDoor;
    impl HackAlgorithm for HackSecondDoor {

        fn hacked(&self, cache: &[(String, usize)]) -> bool {
            cache.iter().fold(HashSet::new(), |mut counter, (digest, _)|{
                if let Ok(i) = digest.chars().nth(5).unwrap().to_string().parse::<usize>() {
                    if i < 8 { counter.insert(i); }
                }
                counter
            }).len() >= 8
        }

        fn extract(&self, cache: &[(String, usize)]) -> String {
            let password = cache.iter().fold([None; 8], |mut password, (digest, _) | {
                if let Ok(i) = digest.chars().nth(5).unwrap().to_string().parse::<usize>() {
                    if i < 8 && password[i].is_none() {
                        password[i] = Some(digest.chars().nth(6).unwrap()); //unwrap to crash if not found
                    }
                }    
                password
            });
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