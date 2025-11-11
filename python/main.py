from tabnanny import check


def calculate_experience_points(level):
    xp = 0
    for i in range(level):
        gain = 5 * i
        xp += gain
        if i == level:
            break
        print(f"xp{xp} gain{gain}")

    return xp


def meditate(mana, max_mana, num_potions):
    while num_potions != 0 and mana != max_mana:
        num_potions -= 1
        mana += 1
        print(f"mana{mana} potions{num_potions}")

    return mana, num_potions


def word_frequency():
    w = "cat bat cat mat bat cat"
    d = {}
    for word in w.split(" "):
        d[word] = d.get(word, 0) + 1
    print(d)


def find_dup_integers():
    ints = {1, 2, 3, 2, 4, 5, 3, 3}


def get_odd_numbers(num):
    odd_numbers = []
    for i in range(0, num):
        if i > 0:
            if i % 2:
                odd_numbers.append(i)
    return odd_numbers


def get_heroes():
    heros = [
        ("Glorfindel", 2093, True),
        ("Gandalf", 1054, False),
        ("Gimli", 389, False),
        ("Aragorn", 87, False),
    ]
    return heros


def reverse_list(items):
    # in place reversal
    # items.reverse()
    # return items

    # new list reversed
    reversed_list = items[::-1]
    return reversed_list


# Expected:
# * filtered messages: ['well it', 'the whole thing', 'kill that knight, it', 'get him!', 'donkey kong', 'oh come on, get them', 'run away from the baddies']
# * words removed: [1, 2, 1, 0, 0, 0, 1]
def filter_messages(messages):
    cleanMessages = []
    dangs = []
    for msg in messages:

        words = msg.split(" ")
        dangCount = 0
        cleanWords = []
        for word in words:
            if word == "dang":
                dangCount += 1
            else:
                cleanWords.append(word)

        cleanSentence = " ".join(cleanWords)
        cleanMessages.append(cleanSentence)
        dangs.append(dangCount)
    return cleanMessages, dangs


def split_players_into_teams(players):
    return players[::2], players[1::2]


def check_ingredient_match(recipe, ingredients):
    matched = [item for item in recipe if item in ingredients]
    percentage = (len(matched) / len(recipe)) * 100
    missing = [item for item in recipe if item not in ingredients]
    return float(f"{percentage:.2f}"), missing


def main():
    # word_frequency()
    # print(calculate_experience_points(4))
    # print(meditate(0, 10, 9))
    # print(f"Expected:[1, 3, 5, 7, 9], Actual:{get_odd_numbers(10)}")
    # print(reverse_list([1, 2, 3, 4, 5]))
    # print(filter_messages(["dang it bobby!", "look at it go"]))
    # print(
    #     filter_messages(
    #         [
    #             "well dang it",
    #             "dang the whole dang thing",
    #             "kill that knight, dang it",
    #             "get him!",
    #             "donkey kong",
    #             "oh come on, get them",
    #             "run away from the dang baddies",
    #         ]
    #     )
    # )

    # Inputs: ['Mike', 'Walter', 'Skyler', 'Tuco']
    # Expected: (['Mike', 'Skyler'], ['Walter', 'Tuco'])
    # print(split_players_into_teams(["Mike", "Walter", "Skyler", "Tuco"]))

    recipe = ["Dragon Scale", "Unicorn Hair", "Phoenix Feather", "Troll Tusk"]
    ingredients = ["Dragon Scale", "Phoenix Feather", "Troll Tusk"]
    percentage, missing_ingredients = check_ingredient_match(recipe, ingredients)

    print('Expected : 75.00 ["Unicorn Hair"]"')
    print(f"Actual: {percentage}, {missing_ingredients}")


if __name__ == "__main__":
    # This block will only execute if the script is run directly.
    # It will not execute if the script is imported as a module.
    main()
