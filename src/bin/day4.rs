use rayon::prelude::*;

#[path = "../utils/aoc_file.rs"] mod aoc_file;

mod room {
    use std::collections::HashMap;
    use rayon::prelude::*;

    pub struct Room {
        checksum: String,
        encoding_parts: Vec<String>,
        num: u32
    }

    impl Room {
        fn parse_checksum(last_part: &str) -> Option<&str> {
            last_part.find('[').and_then(|i|last_part[i + 1 ..].find(']').and_then(|j| {
                Some(&last_part[i + 1 .. i + 1 + j])
            }))
        }

        fn parse_num(last_part: &str) -> Option<u32> {
            last_part.find('[').and_then(|i| last_part[..i].parse::<u32>().ok())
        }

        fn checksum_from_encoding(&self) -> String
        {   
            let char_counts = self.encoding_parts.par_iter()
                .flat_map(|part| part.par_chars())
                .fold(||HashMap::new(), |mut acc, c| { 
                    *acc.entry(c).or_insert(0) += 1; 
                    acc
                })
                .reduce(||HashMap::new(), |mut acc1, acc2|{
                    for (c, v) in acc2 {
                        *acc1.entry(c).or_insert(0) += v; 
                     }
                     acc1
                });
         
            let mut sorted_chars: Vec<char> = char_counts.par_iter().map(|item|{ *item.0 }).collect();
            sorted_chars.par_sort_by_key(|c| { ( -char_counts.get(c).unwrap_or(&0), (*c as i8)) });

            sorted_chars[..sorted_chars.len().min(self.checksum.len())].iter().collect()
        }
        
        const fn letter_count() -> u32 {
            (b'z' - b'a') as u32 + 1
        }
        
        fn decypher_char(c: char, room_num: u32) -> char {
            let c = ((c as u8) - b'a') as u32;
            let c = (c + room_num) % Room::letter_count() + (b'a' as u32);
            c as u8 as char
        }

        pub fn parse_from(room_code: &str) -> Option<Self> {
            let encoding_parts = room_code.split('-').collect::<Vec<&str>>();
            let (last, encoding_parts) = encoding_parts.split_last()?;

            Some(Self {
                checksum: Room::parse_checksum(last)?.to_string(), 
                num: Room::parse_num(last)?, 
                encoding_parts: encoding_parts.iter().map(|x| x.to_string()).collect()
            })
        }

        pub fn get_is_real(&self) -> bool {
            self.checksum_from_encoding() == self.checksum
        }

        pub fn get_num(&self) -> u32 {
            self.num
        }
        
        pub fn get_decrypted_name(&self) -> String {
            self.encoding_parts
                .iter()
                .map(|x| x.chars().map(|c|Room::decypher_char(c, self.num)).collect::<String>())
                .collect::<Vec<_>>()
                .join(" ")
        }

    }

}

fn main() {
    let rooms = match aoc_file::open_and_read_file(&mut std::env::args()) {
        Ok(data) => 
        {
            data.split('\n')
            .collect::<Vec<_>>()
            .par_iter()
            .map(|line| room::Room::parse_from(line.strip_prefix('\r').unwrap_or(line)))
            .collect::<Option<Vec<_>>>()
        }
        Err(_) => {
            eprintln!("Error while reading the file");
            std::process::exit(1);
        }
    };

    let rooms = match &rooms {
        Some(rooms) => rooms.par_iter().filter(|room| room.get_is_real()).collect::<Vec<_>>(),
        None => {
            eprintln!("File parsing error");
            std::process::exit(1);
        }
    };
    
    println!("part 1: {}", rooms.par_iter().map(|room|room.get_num()).sum::<u32>());
    match rooms.par_iter().find_any(|room| room.get_decrypted_name() == "northpole object storage") {
        Some(room) => println!("part 2: {}", room.get_num()),
        None => println!("did not find the northpole object storage!!")
    }
}


#[cfg(test)]
mod tests {
    use crate::room::Room;

    #[test]
    fn test_room_parse() {
        let room = Room::parse_from("aaaaa-bbb-z-y-x-123[abxyz]");
        assert!(room.is_some());

        let room = room.unwrap();
        assert!(room.get_is_real());
        assert!(room.get_num() == 123);
    }

    #[test]
    fn test_room_parse2() {
        let room = Room::parse_from("a-b-c-d-e-f-g-h-987[abcde]");
        assert!(room.is_some());

        let room = room.unwrap();
        assert!(room.get_is_real());
        assert!(room.get_num() == 987);
    }

    #[test]
    fn test_room_parse3() {
        let room = Room::parse_from("not-a-real-room-404[oarel]");
        assert!(room.is_some());

        let room = room.unwrap();
        assert!(room.get_is_real());
        assert!(room.get_num() == 404);
    }
 
    #[test]
    fn test_room_parse4() {
        let room = Room::parse_from("totally-real-room-200[decoy]");
        assert!(room.is_some());

        let room = room.unwrap();
        assert!(!room.get_is_real());
    }

    #[test]
    fn test_room_parse5() {
        let room = Room::parse_from("aaaaa-bbb-z-y-x-123[abxyd]");
        assert!(room.is_some());

        let room = room.unwrap();
        assert!(!room.get_is_real());
    }

    #[test]
    fn test_dycypher_room_name() {
        let room = Room::parse_from("qzmt-zixmtkozy-ivhz-343[zimth]");
        assert!(room.is_some());

        let room = room.unwrap();
        assert!(room.get_is_real());
        assert!(room.get_decrypted_name() == "very encrypted name");
        assert!(room.get_num() == 343);
    }
    
}