import json
movies = json.load(open('latest_movies.json', 'r', encoding='utf-8'))
ids = {}
for m in movies:
    iid = m['id']
    if iid in ids: ids[iid].append(m['title'])
    else: ids[iid] = [m['title']]
dup_ids = {k: v for k, v in ids.items() if len(v) > 1}
print(f'Total de IDs duplicados: {len(dup_ids)}')
print()
for iid, titles in sorted(dup_ids.items()):
    print(f'  ID {iid}: {titles}')
