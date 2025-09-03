## 1. Translation Exercises

### Subtypes

- **block_build** — build translation from blocks
    
- **free_text** — type translation manually
    
- **matching** — match pairs of words
    

### Example: block_build

```json
{
  "id": "tr_001",
  "type": "translation",
  "subtype": "block_build",
  "instruction": "Build the Amharic translation from the blocks.",
  "data": {
    "source_text": "The cat is sleeping under the table.",
    "source_lang": "en",
    "target_lang": "am",
    "blocks": ["ታች", "ጠረጴዛው", "ድመቷ", "ተኝታ", "ነው", "ከ"],
    "correct_sequences": [
      ["ድመቷ", "ከ", "ጠረጴዛው", "ታች", "ተኝታ", "ነው"]
    ]
  }
}
```

### Example: free_text

```json
{
  "id": "tr_002",
  "type": "translation",
  "subtype": "free_text",
  "instruction": "Type the Amharic translation.",
  "data": {
    "source_text": "Good morning",
    "source_lang": "en",
    "target_lang": "am",
    "correct_answers": ["ጤና ይስጥልኝ"],
  }
}
```

### Example: matching

```json
{
  "id": "tr_003",
  "type": "translation",
  "subtype": "matching",
  "instruction": "Match the English words to their Amharic equivalents.",
  "data": {
    "pairs": [
      { "left_id": "l1", "left_text": "Lion", "right_id": "r1", "right_text": "አንበሳ" },
      { "left_id": "l2", "left_text": "Dog",  "right_id": "r2", "right_text": "ውሻ" }
    ],
    "shuffle": true
  }
}
```

---

## 2. Complete-the-Sentence Exercises

### Subtypes

- **given_start** — provide first n words
    
- **given_end** — provide last n words
    
- **select_from_blocks** — omit crucial words, fill from blocks
    

### Example: given_start

```json
{
  "id": "cs_001",
  "type": "complete_sentence",
  "subtype": "given_start",
  "instruction": "Continue the sentence in Amharic.",
  "data": {
    "target_sentence": "እኔ ወደ ገንዘብ ማዕከል እሄዳለሁ።",
    "provided_text": "እኔ ወደ",
    "correct_answers": ["ገንዘብ ማዕከል እሄዳለሁ።"]
  }
}
```

### Example: select_from_blocks

```json
{
  "id": "cs_002",
  "type": "complete_sentence",
  "subtype": "select_from_blocks",
  "instruction": "Fill the missing words using the blocks.",
  "data": {
    "target_sentence": "የቤተሰብ አባት በዓለም ላይ ውድ ነው።",
    "display_with_blanks": "የቤተሰብ ____ በዓለም ላይ ____ ነው።",
    "blocks": ["አባት", "እናቱ", "ውድ"],
    "correct_sequences": [["አባት", "ውድ"]]
  }
}
```

---

## 3. Fill-in-the-Blank Exercises

### Example: one missing word

```json
{
  "id": "fb_001",
  "type": "fill_in_blank",
  "instruction": "Choose the correct word for the blank.",
  "data": {
    "sentence_with_placeholders": "እኔ ____ እበላለሁ።",
    "options": [
      { "id": "o1", "text": "እንጀራ" },
      { "id": "o2", "text": "ስንኩር" },
      { "id": "o3", "text": "ሳምባ" }
    ],
    "correct_option_id": ["o1"]
  }
}
```

### Example: two missing words

```json
{
  "id": "fb_002",
  "type": "fill_in_blank",
  "instruction": "Select the pair that best fits the blanks.",
  "data": {
    "sentence_with_placeholders": "____ ወደ ____ ሄደ።",
    "options": [
      { "id": "o1", "text": "እኔ ... ከቤት" },
      { "id": "o2", "text": "እሱ ... ወደ ትምህርት" },
      { "id": "o3", "text": "እኛ ... ወደ ገበሬ" }
    ],
    "correct_option_id": ["o2"]
  }
}
```

---

## 4. Speaking Exercises

```json
{
  "id": "sp_001",
  "type": "speaking",
  "instruction": "Listen and repeat the sentence.",
  "data": {
    "target_text": "ሰላም እንዴት ነህ?",
    "reference_audio_url": "https://cdn.example.com/audio/selam_endet_neh.mp3",
    "scoring": {
      "min_confidence": 0.7
    },
    "max_record_seconds": 8
  }
}
```

---

## 5. Listening Exercises

### Subtypes

1. **omit_word_two_choices**
    
2. **type_missing**
    
3. **type_sentence**
    
4. **build_from_blocks**
    
5. **matching_audio_text**
    

### Example: omit_word_two_choices

```json
{
  "id": "ls_001",
  "type": "listening",
  "subtype": "omit_word_two_choices",
  "instruction": "Listen and choose the correct missing word.",
  "data": {
    "audio_url": "https://cdn.example.com/audio/sentence_full.mp3",
    "display_text": "እሱ ____ እየጠጣ ነው።",
    "options": [
      { "id": "a", "text": "ነጭ" },
      { "id": "b", "text": "ነጋ" }
    ],
    "correct_option_ids": ["a"]
  }
}
```

### Example: type_sentence

```json
{
  "id": "ls_003",
  "type": "listening",
  "subtype": "type_sentence",
  "instruction": "Listen and type the sentence you hear.",
  "data": {
    "audio_url": "https://cdn.example.com/audio/full_sentence.mp3",
    "correct_answer": ["እኔ በጣም ደስ ብሎኛል"]
  }
}
```

### Example: build_from_blocks

```json
{
  "id": "ls_004",
  "type": "listening",
  "subtype": "build_from_blocks",
  "instruction": "Listen and rebuild the sentence from blocks.",
  "data": {
    "audio_url": "https://cdn.example.com/audio/full_sentence.mp3",
    "blocks": ["እኔ", "ጥሩ", "ነኝ", "በጣም"],
    "correct_sequence": [["እኔ", "በጣም", "ጥሩ", "ነኝ"]]
  }
}
```

### Example: matching_audio_text

```json
{
  "id": "ls_005",
  "type": "listening",
  "subtype": "matching_audio_text",
  "instruction": "Listen to the Amharic audio and match to the English translation.",
  "data": {
    "pairs": [
      { "left_id": "la1", "left_audio": "https://.../word1.mp3", "right_id": "rb1", "right_text": "Water" },
      { "left_id": "la2", "left_audio": "https://.../word2.mp3", "right_id": "rb2", "right_text": "Fire" }
    ]
  }
}
```

## 6. Picture Matching

```json
{
  "id": "pm_001",
  "type": "picture_matching",
  "instruction": "Select the correct image",
  "data": {
    "word": "አንበሳ",
    "audio": "https://example.com/audio/lion.mp3",
    "options": [
      {"id": "a", "image_url": "https://example.com/images/lion.png"},
      {"id": "b", "image_url": "https://example.com/images/dog.png"},
      {"id": "c", "image_url": "https://example.com/images/cat.png"},
      {"id": "d", "image_url": "https://example.com/images/sheep.png"}
    ],
    "correct_option_id": "a"
  }
}
```