CREATE TABLE units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    section_count INT,
    order_index INT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    type TEXT,
    lesson_count INT,
    order_index INT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    section_id UUID NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    title TEXT,
    description TEXT,
    exercise_count INT,
    order_index INT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    subtype TEXT,
    instruction TEXT,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    current_unit UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE
    current_section UUID NOT NULL REFERENCES sections(id) ON DELETE CASCADE
    current_lesson UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE
);

-- Lessons
INSERT INTO lessons (unit_id, section_id, title, description, exercise_count, order_index)
VALUES
((SELECT id FROM units WHERE title = 'Greetings'), (SELECT id FROM sections WHERE type = 'core' AND unit_id = (SELECT id FROM units WHERE title = 'Greetings')), 
 'Lesson 1: Basic Greetings', 'Learn basic greetings, enough to get a conversation going', 3, 1),
((SELECT id FROM units WHERE title = 'Greetings'), (SELECT id FROM sections WHERE type = 'core' AND unit_id = (SELECT id FROM units WHERE title = 'Greetings')), 
 'Lesson 2: How are you?', 'Learn simple conversational phrases, that build on the basic greetings to create natural connections', 3, 2),
((SELECT id FROM units WHERE title = 'Food'), (SELECT id FROM sections WHERE type = 'core' AND unit_id = (SELECT id FROM units WHERE title = 'Food')), 
 'Lesson 1: Fruits', 'Learn fruit vocabulary', 3, 1),
((SELECT id FROM units WHERE title = 'Food'), (SELECT id FROM sections WHERE type = 'core' AND unit_id = (SELECT id FROM units WHERE title = 'Food')), 
 'Lesson 2: Meals', 'Learn meal-related vocabulary', 3, 2);

INSERT INTO exercises (lesson_id, type, subtype, instruction, data)
VALUES
-- Exercises for Lesson 1: Basic Greetings
((SELECT id FROM lessons WHERE title = 'Lesson 1: Basic Greetings'), 'translation', 'block_build', 'Build the Amharic translation from the blocks.',
  '{"source_text": "Hello, how are you?", "source_lang": "en", "target_lang": "am", "blocks": ["እንዴት", "ነህ", "ሰላም", "አይ"], "correct_sequences": [["ሰላም", "እንዴት", "ነህ"]]}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Basic Greetings'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"source_text": "Thank you", "source_lang": "en", "target_lang": "am", "correct_answers": ["አመሰግናለሁ"]}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Basic Greetings'), 'translation', 'matching', 'Match the English greetings to their Amharic equivalents.',
  '{"pairs": [{"left_id": "l1", "left_text": "Hello", "right_id": "r1", "right_text": "ሰላም"}, {"left_id": "l2", "left_text": "Goodbye", "right_id": "r2", "right_text": "ደህና ሁን"}], "shuffle": true}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Basic Greetings'), 'fill_in_blank', NULL, 'Choose the correct word for the blank.',
  '{"sentence_with_placeholders": "____ እንዴት ነህ?", "options": [{"id": "o1", "text": "እንጀራ"}, {"id": "o2", "text": "ሰላም"}, {"id": "o3", "text": "ውሻ"}], "correct_option_id": "o2"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Basic Greetings'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"source_text": "Good morning", "source_lang": "en", "target_lang": "am", "correct_answers": ["እንደምን አደርክ"]}'),

-- Exercises for Lesson 2: How are you?
((SELECT id FROM lessons WHERE title = 'Lesson 2: How are you?'), 'translation', 'block_build', 'Build the Amharic translation from the blocks.',
  '{"source_text": "I am fine, thank you.", "source_lang": "en", "target_lang": "am", "blocks": ["ነኝ", "አመሰግናለሁ", "ደህና", "እኔ"], "correct_sequences": [["ደህና", "ነኝ", "አመሰግናለሁ"]]}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: How are you?'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"source_text": "What is your name?", "source_lang": "en", "target_lang": "am", "correct_answers": ["ስምህ ማን ነው?"]}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: How are you?'), 'complete_sentence', 'select_from_blocks', 'Fill the missing words using the blocks.',
  '{"target_sentence": "ስሜ ማርያም ነው", "display_with_blanks": "ስሜ ____ ነው", "blocks": ["ማርያም", "ውሻ", "ድመት"], "correct_sequences": [["ማርያም"]]}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: How are you?'), 'translation', 'matching', 'Match the English phrases to their Amharic equivalents.',
  '{"pairs": [{"left_id": "l1", "left_text": "How are you?", "right_id": "r1", "right_text": "እንዴት ነህ?"}, {"left_id": "l2", "left_text": "I am fine", "right_id": "r2", "right_text": "ደህና ነኝ"}], "shuffle": true}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: How are you?'), 'speaking', NULL, 'Listen and repeat the sentence.',
  '{"target_text": "ደህና ነኝ", "reference_audio_url": "https://example.com/audio/dehena_negn.mp3", "scoring": {"min_confidence": 0.7}, "max_record_seconds": 5}'),

-- Exercises for Lesson 3: Fruits
((SELECT id FROM lessons WHERE title = 'Lesson 1: Fruits'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"source_text": "Apple", "source_lang": "en", "target_lang": "am", "correct_answers": ["ፖም"]}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Fruits'), 'picture_matching', NULL, 'Select the correct image for the word.',
  '{"word": "ብርቱካን", "audio": "https://example.com/audio/brtukan.mp3", "options": [{"id": "a", "image_url": "https://placehold.co/150x150/ff9900/000000?text=Orange"}, {"id": "b", "image_url": "https://placehold.co/150x150/ffffff/000000?text=Apple"}, {"id": "c", "image_url": "https://placehold.co/150x150/ffd700/000000?text=Banana"}], "correct_option_id": "a"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Fruits'), 'translation', 'block_build', 'Build the Amharic translation from the blocks.',
  '{"source_text": "I eat a banana.", "source_lang": "en", "target_lang": "am", "blocks": ["ሙዝ", "እበላለሁ", "እኔ"], "correct_sequences": [["እኔ", "ሙዝ", "እበላለሁ"]]}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Fruits'), 'fill_in_blank', NULL, 'Choose the correct word for the blank.',
  '{"sentence_with_placeholders": "እኔ ____ እበላለሁ።", "options": [{"id": "o1", "text": "ውሻ"}, {"id": "o2", "text": "ፖም"}, {"id": "o3", "text": "መኪና"}], "correct_option_id": "o2"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Fruits'), 'translation', 'matching', 'Match the English words to their Amharic equivalents.',
  '{"pairs": [{"left_id": "l1", "left_text": "Mango", "right_id": "r1", "right_text": "ማንጎ"}, {"left_id": "l2", "left_text": "Pineapple", "right_id": "r2", "right_text": "አናናስ"}], "shuffle": true}'),

-- Exercises for Lesson 4: Meals
((SELECT id FROM lessons WHERE title = 'Lesson 2: Meals'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"source_text": "Breakfast", "source_lang": "en", "target_lang": "am", "correct_answers": ["ቁርስ"]}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: Meals'), 'translation', 'block_build', 'Build the Amharic translation from the blocks.',
  '{"source_text": "I drink water.", "source_lang": "en", "target_lang": "am", "blocks": ["እጠጣለሁ", "ውሃ", "እኔ"], "correct_sequences": [["እኔ", "ውሃ", "እጠጣለሁ"]]}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: Meals'), 'fill_in_blank', NULL, 'Choose the correct word for the blank.',
  '{"sentence_with_placeholders": "ምሳ ____ ነው።", "options": [{"id": "o1", "text": "ጣፋጭ"}, {"id": "o2", "text": "ውሃ"}, {"id": "o3", "text": "ፖም"}], "correct_option_id": "o1"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: Meals'), 'translation', 'matching', 'Match the English words to their Amharic equivalents.',
  '{"pairs": [{"left_id": "l1", "left_text": "Food", "right_id": "r1", "right_text": "ምግብ"}, {"left_id": "l2", "left_text": "Drink", "right_id": "r2", "right_text": "መጠጥ"}], "shuffle": true}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: Meals'), 'speaking', NULL, 'Listen and repeat the sentence.',
  '{"target_text": "ምሳ በጣም ጣፋጭ ነው", "reference_audio_url": "https://example.com/audio/misa_tatafich.mp3", "scoring": {"min_confidence": 0.7}, "max_record_seconds": 8}');

INSERT INTO user_progress (user_id, current_unit, current_section, current_lesson)
VALUES
('5c8530b5-5056-466c-b82f-82874cf226c4', 
'a5fd6809-9e7e-4f4b-8565-f6c1e7b08b8a', 
'cccf142a-e056-4757-a437-4b93e845231b',
'31eca335-5566-4d4d-bf88-c908cef3e3c7'
)