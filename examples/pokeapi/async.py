import asyncio
import aiohttp
from time import time

error_names = []
count_by_type = {}
total = 0


async def get(pokemon, session):
    async with session.get("https://pokeapi.co/api/v2/pokemon/" + pokemon) as response:
        status_code = response.status
        if status_code != 200:
            error_names.append(pokemon)
        else:
            resp = await response.json()
            for poke_type in resp["types"]:
                if poke_type["type"]["name"] not in count_by_type:
                    count_by_type[poke_type["type"]["name"]] = 0
                else:
                    count_by_type[poke_type["type"]["name"]] += 1
                    global total
                    total += 1


async def main():
    jobs = []
    with open("pokemons.txt") as file:
        names = file.readlines()
        for pokemon in names:
            pokemon = pokemon.strip()

            jobs.append(pokemon)

    async with aiohttp.ClientSession() as session:
        coroutines = [get(pokemon, session) for pokemon in jobs]
        await asyncio.gather(*coroutines)


if __name__ == "__main__":
    initial = time()
    asyncio.run(main())
    print("errors: ", error_names)
    print("count by type: ", count_by_type)
    print("total: ", total)
    print("elapsed time: ", time() - initial)
