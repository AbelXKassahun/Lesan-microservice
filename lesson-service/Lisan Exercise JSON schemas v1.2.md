## 1. Direct Translation Exercises

### Subtypes

- **block_build** — build translation from blocks
- **free_text** — type translation manually

### Example: block_build

```json
{
  "id": "tr_001",
  "type": "translation",
  "subtype": "block_build",
  "instruction": "Translate this sentence",
  "data": {
    "prompt_text": "My brother is taller than my sister.",
    "prompt_audio_url": "https://cdn.example.com/audio/tr_001.mp3",
    "blocks": ["ይረዝማል", "ከእህቴ", "ወንድሜ"],
    "correct_answer": "ወንድሜ ከእህቴ ይረዝማል"
  }
}
```

### Example: free_text

```json
{
  "id": "tr_002",
  "type": "translation",
  "subtype": "free_text",
  "instruction": "Translate this sentence",
  "data": {
    "prompt_text": "I love to read books.",
    "prompt_audio_url": "https://cdn.example.com/audio/tr_002.mp3",
    "correct_answer": "መጽሐፍ ማንበብ እወዳለሁ።"
  }
}
```

---

## 2. Complete-the-Sentence Exercises

### Subtypes

- **partial_free_text** — type the missing part of the sentence
- **partial_block_build** — build the missing part of the sentence from blocks

### Example: partial_free_text

```json
{
  "id": "cs_001",
  "type": "complete_sentence",
  "subtype": "partial_free_text",
  "instruction": "Complete the sentence",
  "data": {
    "reference_text": "She is a doctor.",
    "display_text": "እሷ ____ ናት።",
    "correct_answer": "ሀኪም"
  }
}
```

### Example: partial_block_build

```json
{
  "id": "cs_002",
  "type": "complete_sentence",
  "subtype": "partial_block_build",
  "instruction": "Complete the sentence",
  "data": {
    "reference_text": "This beautiful flower smells good.",
    "display_text": "ይህ ____ አበባ ጥሩ ____ አለው።",
    "blocks": ["ቆንጆ", "መዓዛ", "ቀለም"],
    "correct_answer": "ይህ ቆንጆ አበባ ጥሩ መዓዛ አለው"
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
  "display_text": "እኔ ____ እጠጣለሁ።", 
  "options": [ 
   { 
    "id": 0, 
    "text": "ውሃ" 
   }, 
   { 
    "id": 1, 
    "text": "ዳቦ" 
   }, 
   { 
    "id": 2, 
    "text": "ወንበር" 
   } 
  ], 
  "correct_option_id": 1 
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
  "display_text": "____ ከገበያ ____ ገዛች።", 
  "options": [ 
   { 
    "id": 0, 
    "text": "እሱ ... ዳቦ" 
   }, 
   { 
    "id": 1, 
    "text": "እሷ ... አትክልት" 
   }, 
   { 
    "id": 1, 
    "text": "እነሱ ... ወተት" 
   } 
  ], 
  "correct_option_id": 2 
 } 
} 
```

---

## 4. Speaking Exercises

```json
{
  "id": "sp_001",
  "type": "speaking",
  "instruction": "Speak this sentence aloud",
  "data": {
    "prompt_text": "ሰላም እንዴት ነህ?",
    "prompt_audio_url": "https://cdn.example.com/audio/sp_001.mp3",
    "correct_answer": "ሰላም እንዴት ነህ?"
  }
}
```

---

## 5. Listening Exercises

### Subtypes

1. **choose_missing** — listen and choose the missing word
2. **type_missing** — listen and type the missing word
3. **free_text** — listen and type the full sentence
4. **block_build** — listen and build the sentence from blocks

### Example: choose_missing (with options - single answer)

```json
{
  "id": "ls_001",
  "type": "listening",
  "subtype": "omit_word_choose",
  "instruction": "Listen and choose the correct missing word.",
  "data": {
    "prompt_audio_url": "https://cdn.example.com/audio/ls_001.mp3",
    "display_text": "እሷ ____ ትወዳለች።",
    "options": [
      {
        "id": 0,
        "option_audio_url": "https://cdn.example.com/audio/ls_001_1.mp3",
        "text": "ቡና"
      },
      {
        "id": 1,
        "option_audio_url": "https://cdn.example.com/audio/ls_001_2.mp3",
        "text": "ሻይ"
      }
    ],
    "correct_option_id": 1
  }
}
```

### Example: type_missing (text input - multiple answers)

```json
{
  "id": "ls_002",
  "type": "listening",
  "subtype": "omit_word_type",
  "instruction": "Listen and type the missing word.",
  "data": {
    "prompt_audio_url": "https://cdn.example.com/audio/ls_002.mp3",
    "display_text": "ዛሬ በጣም ____ ነው።",
    "correct_answer": "ትኩስ"
  }
}
```

### Example: free_text (text input - multiple answers)

```json
{
  "id": "ls_003",
  "type": "listening",
  "subtype": "free_text",
  "instruction": "Type what you hear",
  "data": {
    "prompt_audio_url": "https://cdn.example.com/audio/full_sentence.mp3",
    "correct_answer": "ከየት ነህ?"
  }
}
```

### Example: block_build (block arrangement - multiple answers)

```json
{
  "id": "ls_004",
  "type": "listening",
  "subtype": "block_build",
  "instruction": "Tap what you hear",
  "data": {
    "prompt_audio_url": "https://cdn.example.com/audio/ls_004.mp3",
    "blocks": ["እየሄድኩ", "ትምህርት", "ነው", "ቤት"],
    "correct_answer": "ትምህርት ቤት እየሄድኩ ነው"
  }
}
```
