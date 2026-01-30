import sys
from stats import get_num_words
from stats import count_characters
from stats import sorted_report


def get_book_text(filepath):
    with open(filepath, "r") as f:
        return f.read()


def main():
    if len(sys.argv) != 2:
        print("Usage: python3 main.py <path_to_book>")
        sys.exit(1)

    book_path = sys.argv[1]  # <-- Using sys.argv, not hardcoded path
    book_contents = get_book_text(book_path)
    num_words = get_num_words(book_contents)
    print(f"Found {num_words} total words")
    book_dict = count_characters(book_contents)
    print(book_dict)
    sorted_report(book_dict)


if __name__ == "__main__":
    main()
