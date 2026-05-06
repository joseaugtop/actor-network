import json
movies = json.load(open('latest_movies.json', 'r', encoding='utf-8'))
titles = {}
dupes = {}
for m in movies:
    t = m['title']
    if t in titles:
        if t not in dupes: dupes[t] = [titles[t]]
        dupes[t].append(m['id'])
    else: titles[t] = m['id']
print(f'Total de titulos duplicados: {len(dupes)}')
print()
for t, ids in sorted(dupes.items()):
    print(f'  "{t}" -> ids: {ids}')
