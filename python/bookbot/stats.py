from collections import Counter


def get_num_words(book_text):
    split_text = book_text.split()
    nm = len(split_text)
    return nm


def count_characters(book_text):
    return Counter(book_text.lower())


def sorted_report(book_dict: Counter):
    # Sort by count descending (not alphabetically)
    sorted_chars = sorted(book_dict.items(), key=lambda x: x[1], reverse=True)

    # Print each character
    for char, count in sorted_chars:
        print(f"{char}: {count}")
