mod tls_support {
    use rayon::iter::{IntoParallelRefIterator, ParallelIterator};

    fn collect_sequences(ip: &str) -> (Vec<&str>, Vec<&str>) {
        ip.split(']').fold((Vec::new(), Vec::new()), |(mut hyper_seqences, mut regular_seqences), sub_str| {
            let (regular, hypernet) = sub_str.split_once('[').unwrap_or((sub_str, ""));
            if !regular.is_empty() {
                regular_seqences.push(regular);
            }

            if !hypernet.is_empty() {
                hyper_seqences.push(hypernet);
            }

            (hyper_seqences, regular_seqences)
        })
    }

    fn sequence_contains_abba(sequence: &str) -> bool {
        for (i, c) in sequence.chars().enumerate() {
            let remainder = sequence.chars().skip(i + 1).take(3).collect::<Vec<_>>();
            if remainder.len() != 3 {
                return false;
            }

            if c != remainder[0] && c == remainder[2] && remainder[0] == remainder[1] {
                return true;
            } 
        }
        false
    }

    fn collect_abas(seq: &str) -> Vec<String> {
        seq.chars().enumerate().fold(Vec::new(), |mut abas, (i, c)| {
            let remainder = seq.chars().skip(i + 1).take(2).collect::<Vec<_>>();
            if remainder.len() == 2 && c == remainder[1] && c != remainder[0] {
                abas.push(seq.chars().skip(i).take(3).collect());
            }
            abas
        })
    }

    fn contains_bab(seq: &str, aba: &str) -> bool {
        seq.chars().enumerate().any(|(i, _)| {
            let bab: Vec<char> = seq.chars().skip(i).take(3).collect();
            if bab.len() != 3 || bab[0] != bab[2]{
                return false;
            }
            aba.chars().next().unwrap() == bab[1] && aba.chars().nth(1).unwrap() == bab[0]
        })
    }

    pub(crate) fn ip_supports_abba(ip: &str) -> bool {
        let (hyper, regular) = collect_sequences(ip);
        !hyper.iter().any(|seq|sequence_contains_abba(seq)) && regular.iter().any(|seq|sequence_contains_abba(seq))
    }

    pub(crate) fn ip_supports_sll(ip: &str) -> bool {
        let (hyper_sequences, regular_sequences) = collect_sequences(ip);

        let abas = regular_sequences.iter().fold(Vec::new(), |mut abas, sequence| {
            abas.append(&mut collect_abas(sequence));
            abas
        });
        
        hyper_sequences.iter().any(|seqence|abas.iter().any(|aba|contains_bab(seqence, aba)))
    }

    pub fn count_tls_support(ips: &[&str]) -> usize {
        ips.par_iter().filter(|ip| ip_supports_abba(ip)).count()
    }

    pub fn count_ssl_support(ips: &[&str]) -> usize {
        ips.par_iter().filter(|ip| ip_supports_sll(ip)).count()
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    use crate::tls_support::{count_tls_support, count_ssl_support};

    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();
    let lines = content.lines().collect::<Vec<&str>>();

    println!("part1: {}", count_tls_support(&lines));
    println!("part2: {}", count_ssl_support(&lines));
}

#[cfg(test)]
mod tests {
    use crate::tls_support::{count_tls_support, ip_supports_abba, ip_supports_sll, count_ssl_support};

    #[test]
    fn test_tls_support() {
        assert!(ip_supports_abba("abba[mnop]qrst"));
        assert!(!ip_supports_abba("abcd[bddb]xyyx"));
        assert!(!ip_supports_abba("aaaa[qwer]tyui"));
        assert!(ip_supports_abba("ioxxoj[asdfgh]zxcvbn"));
        assert!(!ip_supports_abba("abba[mnop]qrst[abba]"));
        assert!(!ip_supports_abba("[abba]abba[mnop]qrst[abba]"));
    }

    #[test]
    fn test_tls_support_count() {
        assert!(2 == count_tls_support(&["abba[mnop]qrst", "abcd[bddb]xyyx", "aaaa[qwer]tyui", "ioxxoj[asdfgh]zxcvbn"]));
    }

   #[test]
   fn test_ssl_support() {
       assert!(ip_supports_sll("aba[bab]xyz"));
       assert!(!ip_supports_sll("xyx[xyx]xyx"));
       assert!(ip_supports_sll("aaa[kek]eke"));
       assert!(ip_supports_sll("zazbz[bzb]cdb"));
   }

   #[test]
   fn test_ssl_support_count() {
       assert!(3 == count_ssl_support(&["aba[bab]xyz", "xyx[xyx]xyx", "aaa[kek]eke", "zazbz[bzb]cdb"]));
   }
}
