mod decompose {
    use regex::Regex;
    use once_cell::sync::Lazy;

    fn get_marker_values(marker: &str) -> (usize, usize) {
        static MARKER_PARSE_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"(\d+)x(\d+)").unwrap());
        let caps = MARKER_PARSE_REGEX.captures(marker).unwrap();
        (caps[1].parse::<usize>().unwrap(), caps[2].parse::<usize>().unwrap())
    }

    fn count_decomposed<F>(line: &str, char_count_callback: F) -> usize
    where F: Fn(&regex::Match, &str, usize) -> usize {
        static MARKER_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"\(\d+x\d+\)").unwrap());

        match MARKER_REGEX.find(line) {
            Some(mat) => {
                let (char_count, repeat) = get_marker_values(mat.as_str());

                repeat * char_count_callback(&mat, &line, char_count) +
                count_decomposed(&line[mat.end() + char_count..], char_count_callback) + 
                line[..mat.start()].len()
            }
            None => line.len()
        }
    }

    pub(crate) fn count_decomposed_version1(line: &str) -> usize {
        count_decomposed(line, |_, _, char_count|char_count)
    }

    pub(crate) fn count_decomposed_version2(line: &str) -> usize {
        count_decomposed(line, |mat, line, char_count|
            count_decomposed_version2(&line[mat.end()..mat.end() + char_count]))
    }
}

fn main() {
    use aoc_2016::utils::aoc_file;
    let content = aoc_file::open_and_read_file(&mut std::env::args()).unwrap();

    println!("part1: {}", decompose::count_decomposed_version1(&content.chars().filter(|c|!c.is_whitespace()).collect::<String>()));
    println!("part2: {}", decompose::count_decomposed_version2(&content.chars().filter(|c|!c.is_whitespace()).collect::<String>()));
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_decompositions() {
        assert_eq!(6, decompose::count_decomposed_version1("ADVENT"));
        assert_eq!(7, decompose::count_decomposed_version1("A(1x5)BC"));
        assert_eq!(9, decompose::count_decomposed_version1("(3x3)XYZ"));
        assert_eq!(11, decompose::count_decomposed_version1("A(2x2)BCD(2x2)EFG"));
        assert_eq!(6, decompose::count_decomposed_version1("(6x1)(1x3)A"));
        assert_eq!(18, decompose::count_decomposed_version1("X(8x2)(3x3)ABCY"));
    }

    #[test]
    fn test_decompositions_version2() {
        assert_eq!(6, decompose::count_decomposed_version2("ADVENT"));
        assert_eq!(7, decompose::count_decomposed_version2("A(1x5)BC"));
        assert_eq!(9, decompose::count_decomposed_version2("(3x3)XYZ"));
        assert_eq!(20, decompose::count_decomposed_version2("X(8x2)(3x3)ABCY"));
        assert_eq!(445, decompose::count_decomposed_version2("(25x3)(3x3)ABC(2x3)XY(5x2)PQRSTX(18x9)(3x2)TWO(5x7)SEVEN"));
        assert_eq!(241920, decompose::count_decomposed_version2("(27x12)(20x12)(13x14)(7x10)(1x12)A"));
    }
}
