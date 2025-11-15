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
    return ints


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


def get_character_record(name, server, level, rank):
    return {
        "name": name,
        "server": server,
        "level": level,
        "rank": rank,
        "id": f"{name}#{server}",
    }


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


def get_most_common_enemy(enemies_dict):
    if not enemies_dict:
        return None

    common_val = max(enemies_dict.values())
    for key, value in enemies_dict.items():
        if value == common_val:
            return key
        else:
            return None


def get_quest_status(progress):
    character = progress["entity"]["character"]
    q = character.get("quests").get("bridge_run").get("status")
    return q


def merge(dict1, dict2):
    merged_dict = {}
    for key, value in dict1.items():
        merged_dict.update({key: value})
    for key, value in dict2.items():
        merged_dict.update({key: value})
    return merged_dict


def remove_duplicates(spells):
    seen_spells = set(spells)
    u_spells = []
    for spell in seen_spells:
        u_spells.append(spell)
    return u_spells


def count_vowels(text):
    v_count = 0
    words = text.split()

    unique_vowels = set()
    for word in words:
        for letter in word:
            if letter in "aeiou" or letter in "AEIOU":
                v_count += 1
                unique_vowels.update(letter)
    return v_count, unique_vowels


def number_sum(n):
    tally = 0
    for i in range(1, n + 1):
        tally += i
    return tally


def find_min(nums):
    nums.sort()
    return nums[0]


def remove_nonints(nums):
    ints = []
    for i in nums:
        if type(i) is int:
            ints.append(i)
    return ints


def pyhonic_remove_nonints(nums):
    # note: see how we this is done in the [],
    # that means the result is getting added in []
    return [x for x in nums if type(x) is int]


def factorial(n):
    if n == 0 or n == 1:
        return 1
    f = n * factorial((n - 1))
    print(f)
    return f


def area_sum(rectangle):
    sum = 0
    for rect in rectangle:
        h = rect["height"]
        w = rect["width"]
        s = h * w
        sum += s
    return sum


def divide_list(nums, divisor):
    return [num / divisor for num in nums]


def join_strings(strings):
    return ",".join(strings)


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
    # print(split_players_into_teams(["Mike", "Walter", "Skyler", "Tuco"]))
    # recipe = ["Dragon Scale", "Unicorn Hair", "Phoenix Feather", "Troll Tusk"]
    # ingredients = ["Dragon Scale", "Phoenix Feather", "Troll Tusk"]
    # percentage, missing_ingredients = check_ingredient_match(recipe, ingredients)
    # print('Expected : 75.00 ["Unicorn Hair"]"')
    # print(f"Actual: {percentage}, {missing_ingredients}")
    # print(get_most_common_enemy({}))
    # progress = {
    #     "entity": {
    #         "character": {
    #             "name": "Shallan",
    #             "quests": {
    #                 "bridge_run": {
    #                     "status": "Completed",
    #                 },
    #                 "talk_to_syl": {
    #                     "status": "In Progress",
    #                 },
    #             },
    #         }
    #     }
    # }
    # print(get_quest_status(progress))
    # text = "cat dog boiii"
    # print(count_vowels(text))
    # print(number_sum(5))
    # print(find_min([1, 5, 3, 2, -1, 5, 3]))
    # nums = ["200", 300, 2, False, "something", 7, "something else"]
    # print(remove_nonints(nums))
    # print(pyhonic_remove_nonints(nums))
    # print(factorial(5))
    # rectangles = [{"height": 5, "width": 6}, {"height": 4, "width": 7}]
    # print(area_sum(rectangles))
    # print(divide_list([1, 2, 3, 4, 5], 2))
    print(join_strings(["dog", "bat", "cat"]))


if __name__ == "__main__":
    # This block will only execute if the script is run directly.
    # It will not execute if the script is imported as a module.
    main()
