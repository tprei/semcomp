import requests
from time import time

if __name__ == "__main__":
    initial = time()

    error_names = []
    count_by_type = {}

    with open("pokemons.txt") as file:
        names = file.readlines()
        for pokemon in names:
            pokemon = pokemon.strip()

            print("getting", pokemon)
            url = "https://pokeapi.co/api/v2/pokemon/" + pokemon
            r = requests.get(url)

            if r.status_code != 200:
                error_names.append(pokemon)
            else:
                for poke_type in r.json()["types"]:
                    if poke_type["type"]["name"] not in count_by_type:
                        count_by_type[poke_type["type"]["name"]] = 0
                    else:
                        count_by_type[poke_type["type"]["name"]] += 1

    print("errors: ", error_names)
    print("count by type: ", count_by_type)
    print("elapsed time: ", time() - initial)
