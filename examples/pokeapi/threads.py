import requests
import threading
from time import time

error_names = []
count_by_type = {}


def get_pokemon(pokemon, i, errors, types):
    r = requests.get("https://pokeapi.co/api/v2/pokemon/" + pokemon)

    if r.status_code != 200:
        errors[i] = pokemon
    else:
        types[i] = [t["type"]["name"] for t in r.json()["types"]]

    print("finished", i)


if __name__ == "__main__":
    initial = time()

    with open("pokemons.txt") as file:
        names = file.readlines()

        threads = [None] * len(names)
        errors = [None] * len(names)
        types = [None] * len(names)

        for i, pokemon in enumerate(names):
            pokemon = pokemon.strip()

            t = threading.Thread(
                target=get_pokemon,
                args=(
                    pokemon,
                    i,
                    errors,
                    types,
                ),
            )

            threads[i] = t

        for t in threads:
            t.start()

        for t in threads:
            t.join()

        count_by_type = {}
        for type_list in types:
            for t in type_list:
                if t not in count_by_type:
                    count_by_type[t] = 0
                else:
                    count_by_type[t] += 1

        print("errors: ", errors)
        print("count by type: ", count_by_type)
        print("elapsed time: ", time() - initial)
