import json

# Load the JSON file
with open('gerryModel.json', 'r') as f:
    data = json.load(f)

# Extract the spool map and frequency matrix
spool_map = data['spool_map']
freq_mat = data['freq_mat']

# Create a new map for the word-to-word frequency
word_freq_map = {}

spool_map = { str(v): k for (k,v) in spool_map.items()}

for id, mapping in freq_mat.items():
    spool_map['0'] = 'you'
    word1 = spool_map[id]
    total = sum([n for n in mapping.values()])

    mappings = { spool_map[k]: v/total for (k,v) in mapping.items() }
    
    word_freq_map[word1] = mappings


# Save the new map to a JSON file
with open('word_freq_map.json', 'w') as f:
    json.dump(word_freq_map, f, indent=4)