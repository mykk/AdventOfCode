mod firewall_rules {

    pub(crate) fn parse(blocked_ips: &[&str]) -> Vec<(u32, u32)>{
        blocked_ips.iter().map(|x| {
            let (start, end) = x.split_once('-').unwrap();
            (start.parse::<u32>().unwrap(), end.parse::<u32>().unwrap())
        }).collect()
    }

    pub(crate) fn find_smallest_allowed_ip(blocked_ips: &[(u32, u32)], from: u32) -> Option<u32> {
        let mut maybe_smallest_ip = from;
        while let Some(current) = blocked_ips.iter().find(|ip_range| ip_range.0 <= maybe_smallest_ip && ip_range.1 >= maybe_smallest_ip) {
            if u32::MAX == current.1 { 
                return None;
            }
            maybe_smallest_ip = current.1 + 1;
        }
        Some(maybe_smallest_ip)
    }

    pub(crate) fn find_allowed_ip_count(blocked_ips: &[(u32, u32)], from: u32) -> u32 {
        if let Some(current_ip) = find_smallest_allowed_ip(blocked_ips, from) {
            if let Some(min_blocked) = blocked_ips.iter().filter(|ip_range|ip_range.0 > current_ip).min_by_key(|ip_range|ip_range.0) {
                return min_blocked.0 - current_ip + find_allowed_ip_count(blocked_ips, min_blocked.1);
            }
            return u32::MAX - current_ip;
        }

        0
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::firewall_rules::{parse, find_smallest_allowed_ip, find_allowed_ip_count};

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    let lines: Vec<_> = content.lines().collect();

    let blocked_ips = parse(&lines);
    println!("part1: {}", find_smallest_allowed_ip(&blocked_ips, 0).unwrap());
    println!("part2: {}", find_allowed_ip_count(&blocked_ips, 0));
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_example() {
        use crate::firewall_rules::{find_smallest_allowed_ip, find_allowed_ip_count};
        assert_eq!(3, find_smallest_allowed_ip(&[(5, 8), (0, 2), (4, 7)], 0).unwrap());
        assert_eq!(u32::MAX - 8, find_allowed_ip_count(&[(5, 8), (0, 2), (4, 7)], 0));
    }
}
